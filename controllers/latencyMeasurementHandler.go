package controllers

import (
	"context"
	"gitlab-stud.elka.pw.edu.pl/jwojciec/calm-operator.git/api/v1alpha1"
)

// deployment status enum
const (
	PROGRESSING     = "Progressing"
	TRUE            = "True"
	FALSE           = "False"
	REASON_COMPLETE = "NewReplicaSetAvailable"
)

type LatencyMeasurementHandler interface {
	//pewnie przyjmuje clientSet, LatencyMeasurement; zwraca error? coś jeszcze?
	HandleLatencyMeasurement(ctx context.Context, measurement *v1alpha1.LatencyMeasurement, r *LatencyMeasurementReconciler) (err error)
}
