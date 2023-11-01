package controllers

import (
	"context"
	"gitlab-stud.elka.pw.edu.pl/jwojciec/calm-operator.git/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
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
	HandleLatencyMeasurement(measurement *v1alpha1.LatencyMeasurement, r *LatencyMeasurementReconciler, ctx context.Context, req ctrl.Request) (err error)
}
