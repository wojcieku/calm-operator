/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	measurementv1alpha1 "gitlab-stud.elka.pw.edu.pl/jwojciec/calm-operator.git/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"strings"
)

const (
	SERVER      = "server"
	CLIENT      = "client"
	SUCCESS     = "Success"
	FAILURE     = "Failure"
	IN_PROGRESS = "In Progress"
)

var logger = logf.Log.WithName("global")

// LatencyMeasurementReconciler reconciles a LatencyMeasurement object
type LatencyMeasurementReconciler struct {
	client.Client
	Scheme  *runtime.Scheme
	handler LatencyMeasurementHandler
}

//+kubebuilder:rbac:groups=measurement.calm.com,namespace=calm-operator-system,resources=latencymeasurements,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=measurement.calm.com,namespace=calm-operator-system,resources=latencymeasurements/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=measurement.calm.com,namespace=calm-operator-system,resources=latencymeasurements/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,namespace=calm-operator-system,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch,namespace=calm-operator-system,resources=jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,namespace=calm-operator-system,resources=pods,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,namespace=calm-operator-system,resources=services,verbs=get;list;watch;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the LatencyMeasurement object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *LatencyMeasurementReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// TODO z requesta wiadomo ktory to CR -> trzeba wykminic identyfikacje pozostalych zasobow
	measurement := &measurementv1alpha1.LatencyMeasurement{}
	err := r.Get(ctx, req.NamespacedName, measurement)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			logger.Info("No LatencyMeasurements found")
		} else {
			logger.Error(err, "Error during LM get:")
		}
		return ctrl.Result{}, nil
	}
	// stop reconciliation if state is failure already
	if measurement.Status.State == FAILURE {
		return ctrl.Result{}, nil
	}

	switch side := measurement.Spec.Side; side {
	case SERVER:
		r.handler = &ServerSideHandler{}
	case CLIENT:
		r.handler = &ClientSideHandler{}
	default:
		handleIncorrectSide(ctx, measurement, side, r)
		return ctrl.Result{}, nil
	}

	err = r.handler.HandleLatencyMeasurement(measurement, r, ctx, req)

	if err != nil {
		measurement.Status = measurementv1alpha1.LatencyMeasurementStatus{State: FAILURE, Details: err.Error()}
		err := r.Status().Update(ctx, measurement)
		if err != nil {
			logger.Error(err, "Failed to update status after handle failure")
		}
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LatencyMeasurementReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&measurementv1alpha1.LatencyMeasurement{}).
		Owns(&batchv1.Job{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
