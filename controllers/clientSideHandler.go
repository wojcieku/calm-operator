package controllers

import (
	measurementv1alpha1 "gitlab-stud.elka.pw.edu.pl/jwojciec/calm-operator.git/api/v1alpha1"
)

type ClientSideHandler struct {
}

func (handler *ClientSideHandler) HandleLatencyMeasurement(measurement *measurementv1alpha1.LatencyMeasurement, r *LatencyMeasurementReconciler) error {
	//TODO implement me
	//panic("implement me")

	return nil
}
