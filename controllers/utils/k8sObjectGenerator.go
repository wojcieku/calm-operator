package utils

import (
	"gitlab-stud.elka.pw.edu.pl/jwojciec/calm-operator.git/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strconv"
)

const (
	NAMESPACE               = "calm-operator-system"
	LABEL_KEY               = "measurement"
	SERVER_IMAGE            = "jwojciech/udp_probe_server"
	CLIENT_IMAGE            = "jwojciech/udp_probe_client"
	NODE_SELECTOR_HOST_NAME = "kubernetes.io/hostname"
	ADDRESS                 = "ADDRESS"
	INTERVAL                = "INTERVAL"
	DURATION                = "DURATION"
	PORT                    = "PORT"
	METRICS_AGGREGATOR      = "METRICS_AGGREGATOR"
	MEASUREMENT_ID          = "MEASUREMENT_ID"
	SRC_NODE                = "SRC_NODE"
	TARGET_NODE             = "TARGET_NODE"
	SRC_CLUSTER             = "SRC_CLUSTER"
	TARGET_CLUSTER          = "TARGET_CLUSTER"
	APP                     = "app"
)

func PrepareLatencyServerDeployment(deploymentName string, label string, nodeName string, port int) *appsv1.Deployment {
	var replicas int32 = 1

	deployment := &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{
		Name:      deploymentName,
		Namespace: NAMESPACE,
		Labels:    map[string]string{LABEL_KEY: label},
	},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					APP: deploymentName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						APP: deploymentName,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "probe-server",
							Image: SERVER_IMAGE,
							Ports: []corev1.ContainerPort{
								{
									Protocol:      corev1.ProtocolUDP,
									ContainerPort: int32(port),
								},
							},
							Env: []corev1.EnvVar{{Name: PORT, Value: strconv.Itoa(port)}},
						},
					},
					NodeSelector: map[string]string{NODE_SELECTOR_HOST_NAME: nodeName},
				},
			},
		},
	}

	return deployment
}

func PrepareServiceForLatencyServer(svcName string, label string, deploymentName string, ip string, port int) *corev1.Service {

	svc := &corev1.Service{ObjectMeta: metav1.ObjectMeta{
		Name:      svcName,
		Namespace: NAMESPACE,
		Labels: map[string]string{
			LABEL_KEY: label,
		},
	},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Protocol:   "UDP",
					Port:       int32(port),
					TargetPort: intstr.FromInt(port),
				},
			},
			Selector: map[string]string{
				APP: deploymentName,
			},
			ExternalIPs: []string{ip},
		},
	}
	return svc
}

func PrepareJobForLatencyClient(jobName string, nodeName string, label string, client v1alpha1.Client, measurementID string) *batchv1.Job {
	envs := prepareEnvs(client, measurementID)

	var backOffLimit int32 = 0
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: NAMESPACE,
			Labels:    map[string]string{LABEL_KEY: label},
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						APP: jobName,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "probe-client",
							Image: CLIENT_IMAGE,
							Env:   envs,
						},
					},
					NodeSelector:  map[string]string{NODE_SELECTOR_HOST_NAME: nodeName},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
			BackoffLimit: &backOffLimit,
		},
	}
	return job
}

func prepareEnvs(client v1alpha1.Client, measurementID string) []corev1.EnvVar {
	envs := []corev1.EnvVar{}
	envs = append(envs, corev1.EnvVar{Name: ADDRESS, Value: client.IPAddress})
	envs = append(envs, corev1.EnvVar{Name: PORT, Value: strconv.Itoa(client.ServerPort)})
	envs = append(envs, corev1.EnvVar{Name: INTERVAL, Value: strconv.Itoa(client.Interval)})
	envs = append(envs, corev1.EnvVar{Name: DURATION, Value: strconv.Itoa(client.Duration)})
	envs = append(envs, corev1.EnvVar{Name: METRICS_AGGREGATOR, Value: client.MetricsAggregatorAddress})
	envs = append(envs, corev1.EnvVar{Name: MEASUREMENT_ID, Value: measurementID})
	envs = append(envs, corev1.EnvVar{Name: SRC_NODE, Value: client.ClientNodeName})
	envs = append(envs, corev1.EnvVar{Name: TARGET_NODE, Value: client.ServerNodeName})
	envs = append(envs, corev1.EnvVar{Name: SRC_CLUSTER, Value: client.ClientSideClusterName})
	envs = append(envs, corev1.EnvVar{Name: TARGET_CLUSTER, Value: client.ServerSideClusterName})
	return envs
}
