package controllers

import (
	"context"
	measurementv1alpha1 "gitlab-stud.elka.pw.edu.pl/jwojciec/calm-operator.git/api/v1alpha1"
)

// deployment status enum
const (
	PROGRESSING     = "Progressing"
	TRUE            = "True"
	FALSE           = "False"
	REASON_COMPLETE = "NewReplicaSetAvailable"
	PENDING         = "Pending"
	UNSCHEDULABLE   = "Unschedulable"
)

type LatencyMeasurementHandler interface {
	HandleLatencyMeasurement(ctx context.Context, measurement *measurementv1alpha1.LatencyMeasurement, r *LatencyMeasurementReconciler) (err error)
}

func updateStatusSuccess(ctx context.Context, measurement *measurementv1alpha1.LatencyMeasurement, r *LatencyMeasurementReconciler) {
	measurement.Status.State = SUCCESS
	err := r.Status().Update(ctx, measurement)
	if err != nil {
		logger.Error(err, "LM Status update failed")
	}
}
