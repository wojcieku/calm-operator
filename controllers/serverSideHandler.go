package controllers

import (
	measurementv1alpha1 "gitlab-stud.elka.pw.edu.pl/jwojciec/calm-operator.git/api/v1alpha1"
)

type ServerSideHandler struct {
}

func (handler *ServerSideHandler) HandleLatencyMeasurement(measurement *measurementv1alpha1.LatencyMeasurement, r *LatencyMeasurementReconciler) error {
	//TODO implement me

	panic("implement me")

	return nil
}
