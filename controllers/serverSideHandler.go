package controllers

import (
	"context"
	measurementv1alpha1 "gitlab-stud.elka.pw.edu.pl/jwojciec/calm-operator.git/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
)

type ServerSideHandler struct {
}

func (handler *ServerSideHandler) HandleLatencyMeasurement(measurement *measurementv1alpha1.LatencyMeasurement, r *LatencyMeasurementReconciler, ctx context.Context, req ctrl.Request) error {
	//TODO implement me

	//analyze() - list of desired deployments and services

	//desiredServers := measurement.Spec.Servers

	//verify() - compare desired state with actual state

	/* list services
	svcList := &corev1.ServiceList{}
	listOpts := []client.ListOption{
		client.InNamespace(measurement.Namespace),
		client.MatchingLabels{utils.LABEL_KEY : measurement.GetName()},
	}
	r.List(ctx, svcList, listOpts...)
	*/

	//adjust() - perform actions if needed
	//svc := utils.CreateService(measurement.Spec.Servers[0].IpAddress, )

	//updateStatus() - set suitable status of CR

	return nil
}

//znalezienie obiektów w tym namespace i z takimi labelkami
//listOpts := []client.ListOption{
//	client.InNamespace(memcached.Namespace),
//	client.MatchingLabels(),
//}
