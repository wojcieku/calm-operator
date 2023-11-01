package controllers

import (
	"context"
	"errors"
	measurementv1alpha1 "gitlab-stud.elka.pw.edu.pl/jwojciec/calm-operator.git/api/v1alpha1"
	"gitlab-stud.elka.pw.edu.pl/jwojciec/calm-operator.git/controllers/utils"
	appsv1 "k8s.io/api/apps/v1"
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

func (handler *ServerSideHandler) HandleLatencyMeasurement(measurement *measurementv1alpha1.LatencyMeasurement, r *LatencyMeasurementReconciler, ctx context.Context, req ctrl.Request) error {
	//deployment name = CR name + nodeName

	// analyze() - list of desired deployments and services

	desiredServers := measurement.Spec.Servers

	currentDeploys, err := getCurrentDeployments(ctx, measurement, r)
	if err != nil {
		return err
	}

	// verify() - compare desired state with actual state
	missingServers, deploysInProgress, deploymentErr := verifyDeployments(measurement, desiredServers, currentDeploys)
	if deploymentErr != nil {
		return deploymentErr
	}

	// adjust() - perform actions if needed
	for _, server := range missingServers {
		deploymentName := getDeploymentName(measurement, server)
		depl := utils.PrepareLatencyServerDeployment(deploymentName, server.Node, server.Port, measurement.Name)

		// for k8s garbage collection
		_ = ctrl.SetControllerReference(measurement, depl, r.Scheme)

		err := r.Create(ctx, depl)
		if err != nil {
			logger.Error(err, "Error during server deployment creation")
			return err
		}
	}
	// svc := utils.CreateService(measurement.Spec.Servers[0].IpAddress, )

	// updateStatus() - set suitable status of CR
	logger.Info("deploys in progress: " + strconv.Itoa(deploysInProgress))
	logger.Info("missing deploys: " + strconv.Itoa(len(missingServers)))

	if len(missingServers) == 0 && deploysInProgress == 0 {
		logger.Info("All servers deployed successfully")
		measurement.Status.State = SUCCESS
		err := r.Status().Update(ctx, measurement)
		if err != nil {
			logger.Error(err, "LM Status update failed")
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

func verifyDeployments(measurement *measurementv1alpha1.LatencyMeasurement, desiredServers []measurementv1alpha1.Server, currentDeploys *appsv1.DeploymentList) ([]measurementv1alpha1.Server, int, error) {
	var missingServers []measurementv1alpha1.Server
	deploysInProgress := 0

	var err error
	for _, server := range desiredServers {
		completed := false
		inProgress := false
		failed := false

		for _, deployment := range currentDeploys.Items {
			if deployment.Name == (getDeploymentName(measurement, server)) {
				for i, condition := range deployment.Status.Conditions {
					if condition.Type == PROGRESSING {
						switch condition.Reason {
						case TRUE:
							if condition.Status == REASON_COMPLETE {
								completed = true
							} else {
								inProgress = true
							}
						case FALSE:
							failed = true
						}
					}
					logger.Info("Deployment condition " + strconv.Itoa(i) + ": " + condition.String())
				}
			}
		}
		switch {
		case !completed && !inProgress && !failed:
			missingServers = append(missingServers, server)
		case inProgress:
			deploysInProgress++
		case failed:
			err = errors.New("servers deployment failed")
		}
	}
	return missingServers, deploysInProgress, err
}

func getDeploymentName(measurement *measurementv1alpha1.LatencyMeasurement, server measurementv1alpha1.Server) string {
	return measurement.Name + "-" + server.Node
}

//znalezienie obiektów w tym namespace i z takimi labelkami
//listOpts := []client.ListOption{
//	client.InNamespace(memcached.Namespace),
//	client.MatchingLabels(),
//}
