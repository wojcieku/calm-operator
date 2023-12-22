package controllers

import (
	"context"
	"errors"
	measurementv1alpha1 "gitlab-stud.elka.pw.edu.pl/jwojciec/calm-operator.git/api/v1alpha1"
	"gitlab-stud.elka.pw.edu.pl/jwojciec/calm-operator.git/controllers/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
)

type ServerSideHandler struct{}

// HandleLatencyMeasurement - create Deployments and Services objects based on Measurement description.
// Objects' names are created according to pattern: Measurement.Name-NodeName-server.
// If creation fails at any stage, error is returned and for further handle.
// Objects are grouped with labels ["measurement" : LatencyMeasurement.Name].
func (handler *ServerSideHandler) HandleLatencyMeasurement(ctx context.Context, measurement *measurementv1alpha1.LatencyMeasurement, r *LatencyMeasurementReconciler) error {
	desiredServers := measurement.Spec.Servers

	logger.Info("Entering Deployment part")
	deploymentComplete, err := handleServerDeployments(ctx, measurement, r, desiredServers)
	if err != nil {
		return err
	}
	if !deploymentComplete {
		// no error, wait for next reconciliation iteration
		return nil
	}

	logger.Info("Entering SVC part")
	servicesCreationComplete, err := handleServerServices(ctx, measurement, r, desiredServers)
	if err != nil {
		return err
	}
	if !servicesCreationComplete {
		return nil
	}

	// updateStatus() - set suitable status of CR
	if servicesCreationComplete && measurement.Status.State != SUCCESS {
		logger.Info("All servers and services deployed successfully")
		updateStatusSuccess(ctx, measurement, r)
	}
	return nil
}

func handleServerServices(ctx context.Context, measurement *measurementv1alpha1.LatencyMeasurement, r *LatencyMeasurementReconciler, desiredServers []measurementv1alpha1.Server) (bool, error) {
	missingServices, err := verifyServices(ctx, measurement, r, desiredServers)
	serviceCreationComplete := false
	if err != nil {
		return false, err
	}
	err = createMissingServices(ctx, measurement, r, missingServices)
	if err != nil {
		return false, err
	}
	if len(missingServices) == 0 {
		serviceCreationComplete = true
	}
	return serviceCreationComplete, nil
}

func createMissingServices(ctx context.Context, measurement *measurementv1alpha1.LatencyMeasurement, r *LatencyMeasurementReconciler, missingServices []measurementv1alpha1.Server) error {
	for _, server := range missingServices {
		logger.Info("creating service for server")
		serviceName := getServerObjectsName(measurement, server)
		svc := utils.PrepareServiceForLatencyServer(serviceName, measurement.Name, serviceName, server.IPAddress, server.Port)

		// for k8s garbage collection
		_ = ctrl.SetControllerReference(measurement, svc, r.Scheme)

		err := r.Create(ctx, svc)
		if err != nil {
			logger.Error(err, "Error during server deployment creation")
			return err
		}
	}
	return nil
}

func verifyServices(ctx context.Context, measurement *measurementv1alpha1.LatencyMeasurement, r *LatencyMeasurementReconciler, desiredServers []measurementv1alpha1.Server) ([]measurementv1alpha1.Server, error) {
	currentServices, err := getCurrentServices(ctx, measurement, r)
	if err != nil {
		return nil, err
	}
	var missingServices []measurementv1alpha1.Server
	for _, server := range desiredServers {
		exists := false
		for _, service := range currentServices.Items {
			if service.Name == getServerObjectsName(measurement, server) {
				exists = true
			}
		}
		if !exists {
			missingServices = append(missingServices, server)
		}
	}
	return missingServices, nil
}

func handleServerDeployments(ctx context.Context, measurement *measurementv1alpha1.LatencyMeasurement, r *LatencyMeasurementReconciler, desiredServers []measurementv1alpha1.Server) (bool, error) {
	currentDeploys, err := getCurrentDeployments(ctx, measurement, r)
	deployCompleted := false
	if err != nil {
		return false, err
	}

	// verify - compare desired state with actual state
	missingDeployments, deploysInProgress, err := verifyDeployments(measurement, desiredServers, currentDeploys)
	if err != nil {
		return false, err
	}

	// adjust() - perform actions if needed
	err = createMissingDeployments(ctx, measurement, r, missingDeployments)
	if err != nil {
		return false, err
	}

	logger.Info("deploys in progress: " + strconv.Itoa(deploysInProgress))
	logger.Info("missing deploys: " + strconv.Itoa(len(missingDeployments)))

	if len(missingDeployments) == 0 {
		if deploysInProgress == 0 {
			deployCompleted = true
		} else {
			// check if any Pod schedule failed
			err = checkServerPodsScheduleStatus(ctx, measurement, r, desiredServers)
			if err != nil {
				return false, err
			}
		}
	}
	return deployCompleted, nil
}

// Returns error with details if any Pod failed to be scheduled.
func checkServerPodsScheduleStatus(ctx context.Context, measurement *measurementv1alpha1.LatencyMeasurement, r *LatencyMeasurementReconciler, desiredServers []measurementv1alpha1.Server) error {
	podsList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(measurement.Namespace),
	}
	err := r.List(ctx, podsList, listOpts...)
	if err != nil {
		logger.Error(err, "Error during listing pods")
		return err
	}
	for _, pod := range podsList.Items {
		for _, server := range desiredServers {
			if pod.ObjectMeta.GetLabels()[utils.APP] == getServerObjectsName(measurement, server) && pod.Status.Phase == PENDING {
				for _, condition := range pod.Status.Conditions {
					if condition.Status == FALSE && condition.Reason == UNSCHEDULABLE {
						logger.Info("POD UNSCHEDULABLE")
						err = errors.New("Pod unschedulable for deployment: " + pod.ObjectMeta.GetLabels()[utils.APP])
						return err
					}
				}
			}
		}
	}
	return err
}

func createMissingDeployments(ctx context.Context, measurement *measurementv1alpha1.LatencyMeasurement, r *LatencyMeasurementReconciler, missingDeployments []measurementv1alpha1.Server) error {
	for _, server := range missingDeployments {
		logger.Info("creating server deployment")
		deploymentName := getServerObjectsName(measurement, server)
		depl := utils.PrepareLatencyServerDeployment(deploymentName, measurement.Name, server.Node, server.Port)

		// for k8s garbage collection
		_ = ctrl.SetControllerReference(measurement, depl, r.Scheme)

		err := r.Create(ctx, depl)
		if err != nil {
			logger.Error(err, "Error during server deployment creation")
			return err
		}
	}
	return nil
}

func getCurrentDeployments(ctx context.Context, measurement *measurementv1alpha1.LatencyMeasurement, r *LatencyMeasurementReconciler) (*appsv1.DeploymentList, error) {
	currentDeploys := &appsv1.DeploymentList{}
	listOpts := []client.ListOption{
		client.InNamespace(measurement.Namespace),
		client.MatchingLabels{utils.LABEL_KEY: measurement.GetName()},
	}
	err := r.List(ctx, currentDeploys, listOpts...)
	if err != nil {
		logger.Error(err, "Error during listing deployments")
		return nil, err
	}
	return currentDeploys, nil
}

func getCurrentServices(ctx context.Context, measurement *measurementv1alpha1.LatencyMeasurement, r *LatencyMeasurementReconciler) (*corev1.ServiceList, error) {
	currentServices := &corev1.ServiceList{}
	listOpts := []client.ListOption{
		client.InNamespace(measurement.Namespace),
		client.MatchingLabels{utils.LABEL_KEY: measurement.GetName()},
	}
	err := r.List(ctx, currentServices, listOpts...)
	if err != nil {
		logger.Error(err, "Error during listing services")
		return nil, err
	}
	return currentServices, nil
}

func verifyDeployments(measurement *measurementv1alpha1.LatencyMeasurement, desiredServers []measurementv1alpha1.Server, currentDeploys *appsv1.DeploymentList) ([]measurementv1alpha1.Server, int, error) {
	var missingDeployments []measurementv1alpha1.Server
	deploysInProgress := 0

	var err error
	for _, server := range desiredServers {
		exists := false
		for _, deployment := range currentDeploys.Items {
			if deployment.Name == (getServerObjectsName(measurement, server)) {
				exists = true
				// deployment initial conditions are empty
				if deployment.Status.Conditions == nil {
					deploysInProgress++
				}
				for _, condition := range deployment.Status.Conditions {
					if condition.Type == PROGRESSING {
						switch condition.Status {
						case TRUE:
							if condition.Reason != REASON_COMPLETE {
								deploysInProgress++
							}
						case FALSE:
							err = errors.New("servers deployment failed")
						}
					}
				}
			}
		}
		if !exists {
			missingDeployments = append(missingDeployments, server)
		}
	}
	return missingDeployments, deploysInProgress, err
}

func getServerObjectsName(measurement *measurementv1alpha1.LatencyMeasurement, server measurementv1alpha1.Server) string {
	return measurement.Name + "-" + server.Node + "-server"
}
