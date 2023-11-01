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

/*
	list services

svcList := &corev1.ServiceList{}

	listOpts := []client.ListOption{
		client.InNamespace(measurement.Namespace),
		client.MatchingLabels{utils.LABEL_KEY : measurement.GetName()},
	}

r.List(ctx, svcList, listOpts...)
*/
type ServerSideHandler struct{}

func (handler *ServerSideHandler) HandleLatencyMeasurement(ctx context.Context, measurement *measurementv1alpha1.LatencyMeasurement, r *LatencyMeasurementReconciler) error {
	// every Server is mapped to Deployment and Service with name: <LatencyMeasurementName>-<NodeName>

	// analyze() - list of desired deployments and services

	desiredServers := measurement.Spec.Servers

	deploymentComplete, err := handleServerDeployments(ctx, measurement, r, desiredServers)
	if err != nil {
		return err
	}
	if !deploymentComplete {
		return nil
	}
	// TODO implement services creation
	currentServices, err := getCurrentServices(ctx, measurement, r)
	if err != nil {
		return err
	}
	var missingServices []measurementv1alpha1.Server
	for _, server := range desiredServers {
		exists := false
		for _, service := range currentServices.Items {
			if service.Name == getServerObjectsName(measurement, server) {
				// logi zeby zobaczyc co sie dzieje z servicem
				for i, condition := range service.Status.Conditions {
					logger.Info("Service condition #" + strconv.Itoa(i) + condition.String())
				}
				exists = true
			}
		}
		if !exists {
			missingServices = append(missingServices, server)
		}
	}

	for _, server := range missingServices {
		logger.Info("creating service for server")
		serviceName := getServerObjectsName(measurement, server)
		svc := utils.PrepareServiceForLatencyServer(server.IpAddress, server.Port, serviceName, serviceName, measurement.Name)

		// for k8s garbage collection
		_ = ctrl.SetControllerReference(measurement, svc, r.Scheme)

		err := r.Create(ctx, svc)
		if err != nil {
			logger.Error(err, "Error during server deployment creation")
			return err
		}
	}

	// svc := utils.CreateService(measurement.Spec.Servers[0].IpAddress, )

	// updateStatus() - set suitable status of CR
	logger.Info("All servers and services deployed successfully")
	measurement.Status.State = SUCCESS
	err = r.Status().Update(ctx, measurement)
	if err != nil {
		logger.Error(err, "LM Status update failed")
	}

	return nil
}

func handleServerDeployments(ctx context.Context, measurement *measurementv1alpha1.LatencyMeasurement, r *LatencyMeasurementReconciler, desiredServers []measurementv1alpha1.Server) (bool, error) {
	currentDeploys, err := getCurrentDeployments(ctx, measurement, r)
	deployCompleted := false
	if err != nil {
		return false, err
	}

	// verify() - compare desired state with actual state
	missingDeployments, deploysInProgress, err := verifyDeployments(measurement, desiredServers, currentDeploys)
	if err != nil {
		return false, err
	}

	// adjust() - perform actions if needed
	// TODO extract creation
	for _, server := range missingDeployments {
		logger.Info("creating server deployment")
		deploymentName := getServerObjectsName(measurement, server)
		depl := utils.PrepareLatencyServerDeployment(deploymentName, server.Node, server.Port, measurement.Name)

		// for k8s garbage collection
		_ = ctrl.SetControllerReference(measurement, depl, r.Scheme)

		err := r.Create(ctx, depl)
		if err != nil {
			logger.Error(err, "Error during server deployment creation")
			return false, err
		}
	}
	logger.Info("deploys in progress: " + strconv.Itoa(deploysInProgress))
	logger.Info("missing deploys: " + strconv.Itoa(len(missingDeployments)))
	if len(missingDeployments) == 0 && deploysInProgress == 0 {
		deployCompleted = true
	}
	return deployCompleted, nil
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
	logger.Info("Current deployments: " + currentDeploys.String())
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
		logger.Error(err, "Error during listing deployments")
		return nil, err
	}
	logger.Info("Current deployments: " + currentServices.String())
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
				for i, condition := range deployment.Status.Conditions {
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
					logger.Info("Deployment condition " + strconv.Itoa(i) + ": " + condition.String())
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
	return measurement.Name + "-" + server.Node
}

//znalezienie obiektów w tym namespace i z takimi labelkami
//listOpts := []client.ListOption{
//	client.InNamespace(memcached.Namespace),
//	client.MatchingLabels(),
//}
