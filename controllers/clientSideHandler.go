package controllers

import (
	"context"
	"errors"
	measurementv1alpha1 "gitlab-stud.elka.pw.edu.pl/jwojciec/calm-operator.git/api/v1alpha1"
	"gitlab-stud.elka.pw.edu.pl/jwojciec/calm-operator.git/controllers/utils"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

type ClientSideHandler struct {
}

func (handler *ClientSideHandler) HandleLatencyMeasurement(ctx context.Context, measurement *measurementv1alpha1.LatencyMeasurement, r *LatencyMeasurementReconciler) error {
	desiredClients := measurement.Spec.Clients

	currentJobs, err := getCurrentJobs(ctx, measurement, r)
	if err != nil {
		return err
	}

	// check Jobs statuses, find missing Jobs
	missingClients, inProgress, err := verifyJobs(desiredClients, currentJobs, measurement)
	if err != nil {
		return err
	}

	// create Jobs if missing
	err = createMissingJobs(ctx, measurement, r, missingClients)
	if err != nil {
		logger.Error(err, "Error during listing Jobs")
		return err
	}

	// set success status if all jobs succeeded
	if len(missingClients) == 0 && !inProgress && measurement.Status.State != SUCCESS {
		logger.Info("All Jobs completed successfully")
		updateStatusSuccess(ctx, measurement, r)
	} else {
		err = checkClientPodsScheduleStatus(ctx, measurement, r, desiredClients)
		if err != nil {
			return err
		}
	}
	return nil
}

func createMissingJobs(ctx context.Context, measurement *measurementv1alpha1.LatencyMeasurement, r *LatencyMeasurementReconciler, missingClients []measurementv1alpha1.Client) error {
	for _, missingClient := range missingClients {
		job := utils.PrepareJobForLatencyClient(getClientObjectsName(measurement, missingClient), missingClient.ClientNodeName, measurement.Name,
			missingClient, measurement.Name)

		_ = ctrl.SetControllerReference(measurement, job, r.Scheme)
		err := r.Create(ctx, job)
		if err != nil {
			logger.Error(err, "Error during client job creation")
			// Events arrive one per object change and k8s server may process all creation requests
			// after listing current deployments, so it may already exist - assuming operator is namespace scoped and
			// names are unique - that's ok
			if !strings.Contains(err.Error(), "already exists") {
				return err
			}
		}
	}
	return nil
}

func verifyJobs(desiredClients []measurementv1alpha1.Client, currentJobs *batchv1.JobList, measurement *measurementv1alpha1.LatencyMeasurement) ([]measurementv1alpha1.Client, bool, error) {
	var missingClients []measurementv1alpha1.Client
	inProgress := false
	for _, c := range desiredClients {
		exists := false
		for _, job := range currentJobs.Items {
			if job.Name == getClientObjectsName(measurement, c) {
				exists = true
				logger.Info("Job " + job.Name + " status: " + job.Status.String())
				if job.Status.Failed == 1 {
					err := errors.New("Job failed for client: " + job.Name)
					logger.Error(err, "Job execution error")
					// TODO pobranie logow poda?
					return nil, inProgress, err
				}
				if job.Status.Succeeded != 1 {
					inProgress = true
				}
			}
		}
		if !exists {
			missingClients = append(missingClients, c)
		}
	}
	return missingClients, inProgress, nil
}

func getCurrentJobs(ctx context.Context, measurement *measurementv1alpha1.LatencyMeasurement, r *LatencyMeasurementReconciler) (*batchv1.JobList, error) {
	currentJobs := &batchv1.JobList{}
	listOpts := []client.ListOption{
		client.InNamespace(measurement.Namespace),
		client.MatchingLabels{utils.LABEL_KEY: measurement.GetName()},
	}
	err := r.List(ctx, currentJobs, listOpts...)
	if err != nil {
		logger.Error(err, "Error during listing Jobs")
		return currentJobs, err
	}
	return currentJobs, err
}

// Returns error with details if any Pod failed to be scheduled.
func checkClientPodsScheduleStatus(ctx context.Context, measurement *measurementv1alpha1.LatencyMeasurement, r *LatencyMeasurementReconciler, desiredClients []measurementv1alpha1.Client) error {
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
		for _, c := range desiredClients {
			if pod.ObjectMeta.GetLabels()[utils.APP] == getClientObjectsName(measurement, c) && pod.Status.Phase == PENDING {
				for _, condition := range pod.Status.Conditions {
					if condition.Status == FALSE && condition.Reason == UNSCHEDULABLE {
						logger.Info("POD UNSCHEDULABLE")
						err = errors.New("Pod unschedulable for job: " + pod.ObjectMeta.GetLabels()[utils.APP])
						return err
					}
				}
			}
		}
	}
	return err
}

func getClientObjectsName(measurement *measurementv1alpha1.LatencyMeasurement, client measurementv1alpha1.Client) string {
	return measurement.Name + "-" + client.ClientNodeName + "-to-" + client.ServerNodeName + "-client"
}
