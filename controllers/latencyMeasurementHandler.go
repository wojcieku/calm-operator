package controllers

import (
	"gitlab-stud.elka.pw.edu.pl/jwojciec/calm-operator.git/api/v1alpha1"
)

type LatencyMeasurementHandler interface {
	//pewnie przyjmuje clientSet, LatencyMeasurement; zwraca error? coś jeszcze?
	HandleLatencyMeasurement(measurement *v1alpha1.LatencyMeasurement, r *LatencyMeasurementReconciler) (err error)
}

//znalezienie obiektów w tym namespace i z takimi labelkami
//listOpts := []client.ListOption{
//	client.InNamespace(memcached.Namespace),
//	client.MatchingLabels(),
//}
