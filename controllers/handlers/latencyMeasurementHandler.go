package handlers

import (
	"gitlab-stud.elka.pw.edu.pl/jwojciec/calm-operator.git/api/v1alpha1"
	"gitlab-stud.elka.pw.edu.pl/jwojciec/calm-operator.git/controllers"
)

type LatencyMeasurementHandler interface {
	//pewnie przyjmuje clientSet, LatencyMeasurement; zwraca error? coś jeszcze?
	HandleLatencyMeasurement(measurement *v1alpha1.LatencyMeasurement, r *controllers.LatencyMeasurementReconciler) (err error)
}

//znalezienie obiektów w tym namespace i z takimi labelkami
//listOpts := []client.ListOption{
//	client.InNamespace(memcached.Namespace),
//	client.MatchingLabels(),
//}
