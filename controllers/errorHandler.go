package controllers

import (
	"context"
	"errors"
	measurementv1alpha1 "gitlab-stud.elka.pw.edu.pl/jwojciec/calm-operator.git/api/v1alpha1"
)

func handleIncorrectSide(ctx context.Context, measurement *measurementv1alpha1.LatencyMeasurement, side string, r *LatencyMeasurementReconciler) {
	if measurement.Status.State != FAILURE {
		logger.Error(errors.New("Unknown side specified in LatencyMeasurement Spec with name: "+measurement.Name), "Wrong Latency Measurement spec")
		measurement.Status = measurementv1alpha1.LatencyMeasurementStatus{State: FAILURE, Details: "Expected Server or Client in spec.side, got: " + side}

		err := r.Status().Update(ctx, measurement)
		if err != nil {
			logger.Error(err, "LM Status update failed")
		}
	}
}
