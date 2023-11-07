package controllers

import (
	"context"
	"errors"
	measurementv1alpha1 "gitlab-stud.elka.pw.edu.pl/jwojciec/calm-operator.git/api/v1alpha1"
	"gitlab-stud.elka.pw.edu.pl/jwojciec/calm-operator.git/controllers/utils"
	batchv1 "k8s.io/api/batch/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ClientSideHandler struct {
}

func (handler *ClientSideHandler) HandleLatencyMeasurement(ctx context.Context, measurement *measurementv1alpha1.LatencyMeasurement, r *LatencyMeasurementReconciler) error {
	desiredClients := measurement.Spec.Clients

	// getJobs
	currentJobs := &batchv1.JobList{}
	listOpts := []client.ListOption{
		client.InNamespace(measurement.Namespace),
		client.MatchingLabels{utils.LABEL_KEY: measurement.GetName()},
	}
	err := r.List(ctx, currentJobs, listOpts...)
	if err != nil {
		logger.Error(err, "Error during listing Jobs")
		return err
	}

	// verifyJobs() - check Jobs statuses, find missing Jobs
	var missingClients []measurementv1alpha1.Client
	inProgress := false
	for _, c := range desiredClients {
		exists := false
		for _, job := range currentJobs.Items {
			if job.Name == getClientObjectsName(measurement, c) {
				exists = true
				logger.Info("Job " + job.Name + " status: " + job.Status.String())
				if job.Status.Failed == 1 {
					err = errors.New("Job failed for client: " + job.Name)
					return err
				}
				if job.Status.Succeeded != 1 {
					inProgress = true
				}

			}
			if !exists {
				missingClients = append(missingClients, c)
			}
		}
	}

	// create Jobs if missing
	for _, missingClient := range missingClients {
		job := utils.PrepareJobForLatencyClient(getClientObjectsName(measurement, missingClient), measurement.Name,
			missingClient.IpAddress, missingClient.Port, missingClient.Interval, missingClient.Duration)

		_ = ctrl.SetControllerReference(measurement, job, r.Scheme)
		err := r.Create(ctx, job)
		if err != nil {
			logger.Error(err, "Error during client job creation")
			return err
		}
	}

	// return if all succeeded
	if len(missingClients) == 0 && !inProgress && measurement.Status.State != SUCCESS {
		logger.Info("All Jobs completed successfully")
		measurement.Status.State = SUCCESS
		err = r.Status().Update(ctx, measurement)
		if err != nil {
			logger.Error(err, "LM Status update failed")
		}
		return nil
	}
	return nil
}

func getClientObjectsName(measurement *measurementv1alpha1.LatencyMeasurement, client measurementv1alpha1.Client) string {
	return measurement.Name + "-" + client.Node
}
