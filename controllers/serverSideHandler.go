package controllers

import (
	"context"
	measurementv1alpha1 "gitlab-stud.elka.pw.edu.pl/jwojciec/calm-operator.git/api/v1alpha1"
	"gitlab-stud.elka.pw.edu.pl/jwojciec/calm-operator.git/controllers/utils"
	appsv1 "k8s.io/api/apps/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
	for _, server := range desiredServers {
		completed := false
		for _, deployment := range currentDeploys.Items {
			if deployment.Name == (measurement.Name+server.Node) && deployment.Status.ReadyReplicas == 1 {
				completed = true
			}
		}
		if !completed {
			missingServers = append(missingServers, server)
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
		deploymentName := measurement.Name + server.Node
		depl := utils.PrepareLatencyServerDeployment(deploymentName, server.Node, server.Port, measurement.Name)

		err := r.Create(ctx, depl)
		if err != nil {
			logger.Error(err, "Error during server deployment creation")
			return err
		}
	}
	//svc := utils.CreateService(measurement.Spec.Servers[0].IpAddress, )

	//updateStatus() - set suitable status of CR
	if len(missingServers) == 0 {
		logger.Info("All server deployed successfully")
	}
	return nil
}

//znalezienie obiektów w tym namespace i z takimi labelkami
//listOpts := []client.ListOption{
//	client.InNamespace(memcached.Namespace),
//	client.MatchingLabels(),
//}
