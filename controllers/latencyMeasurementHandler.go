package controllers

import (
	"context"
	measurementv1alpha1 "gitlab-stud.elka.pw.edu.pl/jwojciec/calm-operator.git/api/v1alpha1"
)

// Kubernetes object constants
const (
	PROGRESSING     = "Progressing"
	TRUE            = "True"
	FALSE           = "False"
	REASON_COMPLETE = "NewReplicaSetAvailable"
	PENDING         = "Pending"
	UNSCHEDULABLE   = "Unschedulable"
	K8S_ARCH        = "kubernetes.io/arch"
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
