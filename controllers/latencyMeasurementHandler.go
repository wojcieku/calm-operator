package controllers

import (
	"context"
	"errors"
	measurementv1alpha1 "gitlab-stud.elka.pw.edu.pl/jwojciec/calm-operator.git/api/v1alpha1"
	"gitlab-stud.elka.pw.edu.pl/jwojciec/calm-operator.git/controllers/utils"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

// Returns error with details if any Pod failed to be scheduled.
func checkPodScheduleStatus[E measurementv1alpha1.Server | measurementv1alpha1.Client](ctx context.Context, measurement *measurementv1alpha1.LatencyMeasurement, r *LatencyMeasurementReconciler, probeObjects []E) error {
	podsList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(measurement.Namespace),
	}
	err := r.List(ctx, podsList, listOpts...)
	if err != nil {
		logger.Error(err, "Error during listing pods")
		return err
	}
	for _, pod := range podsList.Items {
		for _, object := range probeObjects {
			if pod.ObjectMeta.GetLabels()[utils.APP] == getObjectName(measurement, object) && pod.Status.Phase == PENDING {
				for _, condition := range pod.Status.Conditions {
					if condition.Status == FALSE && condition.Reason == UNSCHEDULABLE {
						logger.Info("POD UNSCHEDULABLE")
						err = errors.New("Pod unschedulable for deployment: " + pod.ObjectMeta.GetLabels()[utils.APP])
						return err
					}
				}
			}
		}
	}
	return err
}

func getObjectName(measurement *measurementv1alpha1.LatencyMeasurement, object measurementv1alpha1.DeployableObject) string {
	return measurement.Name + "-" + object.GetNodeName()
}
