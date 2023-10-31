package controllers

import (
	"context"
	measurementv1alpha1 "gitlab-stud.elka.pw.edu.pl/jwojciec/calm-operator.git/api/v1alpha1"
	"gitlab-stud.elka.pw.edu.pl/jwojciec/calm-operator.git/controllers/utils"
	appsv1 "k8s.io/api/apps/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
)

type ServerSideHandler struct {
}

func (handler *ServerSideHandler) HandleLatencyMeasurement(measurement *measurementv1alpha1.LatencyMeasurement, r *LatencyMeasurementReconciler, ctx context.Context, req ctrl.Request) error {
	//TODO implement me
	//deployment name = CR name + nodeName

	//analyze() - list of desired deployments and services

	desiredServers := measurement.Spec.Servers

	currentDeploys := &appsv1.DeploymentList{}
	listOpts := []client.ListOption{
		client.InNamespace(measurement.Namespace),
		client.MatchingLabels{utils.LABEL_KEY: measurement.GetName()},
	}
	err := r.List(ctx, currentDeploys, listOpts...)
	if err != nil {
		logger.Error(err, "Error during listing deployments")
		return err
	}

	//verify() - compare desired state with actual state
	var missingServers []measurementv1alpha1.Server
	deploysInProgressCounter := 0
	for _, server := range desiredServers {
		completed := false
		inProgress := false

		//TODO detect deployment failure
		for _, deployment := range currentDeploys.Items {
			for i, condition := range deployment.Status.Conditions {
				logger.Info("Deployment condition " + strconv.Itoa(i) + ": " + condition.String())
			}
			if deployment.Name == (getDeploymentName(measurement, server)) {

				logger.Info(strconv.Itoa(int(deployment.Status.ReadyReplicas)))
				if deployment.Status.ReadyReplicas == 1 {
					completed = true
				} else {
					inProgress = true
				}
			}
		}
		if !completed && !inProgress {
			missingServers = append(missingServers, server)
		} else if inProgress {
			deploysInProgressCounter++
		}
	}
	/* list services
	svcList := &corev1.ServiceList{}
	listOpts := []client.ListOption{
		client.InNamespace(measurement.Namespace),
		client.MatchingLabels{utils.LABEL_KEY : measurement.GetName()},
	}
	r.List(ctx, svcList, listOpts...)
	*/

	//adjust() - perform actions if needed
	//TODO set controller reference
	for _, server := range missingServers {
		deploymentName := getDeploymentName(measurement, server)
		depl := utils.PrepareLatencyServerDeployment(deploymentName, server.Node, server.Port, measurement.Name)

		//for k8s garbage collection
		_ = ctrl.SetControllerReference(measurement, depl, r.Scheme)

		err := r.Create(ctx, depl)
		if err != nil {
			logger.Error(err, "Error during server deployment creation")
			return err
		}
	}
	//svc := utils.CreateService(measurement.Spec.Servers[0].IpAddress, )

	//updateStatus() - set suitable status of CR
	logger.Info("deploys in progress: " + strconv.Itoa(deploysInProgressCounter))
	logger.Info("missing deploys: " + strconv.Itoa(len(missingServers)))
	if len(missingServers) == 0 && deploysInProgressCounter == 0 {
		logger.Info("All servers deployed successfully")
		measurement.Status.State = SUCCESS
		err := r.Status().Update(ctx, measurement)
		if err != nil {
			logger.Error(err, "LM Status update failed")
		}
	}
	return nil
}

func getDeploymentName(measurement *measurementv1alpha1.LatencyMeasurement, server measurementv1alpha1.Server) string {
	return measurement.Name + "-" + server.Node
}

//znalezienie obiektów w tym namespace i z takimi labelkami
//listOpts := []client.ListOption{
//	client.InNamespace(memcached.Namespace),
//	client.MatchingLabels(),
//}
