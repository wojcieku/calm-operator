package utils

//TODO getService, getDeployment, getJob
import (
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
					"app": deploymentName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": deploymentName,
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
				"apps": deploymentName,
			},
			ExternalIPs: []string{ip},
		},
	}
	return svc
}

func PrepareJobForLatencyClient(jobName string, label string, ip string, port int, interval int, duration int) *batchv1.Job {
	envs := prepareEnvs(ip, port, interval, duration)

	var backOffLimit int32 = 1
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: NAMESPACE,
			Labels:    map[string]string{LABEL_KEY: label},
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "probe-client",
							Image: CLIENT_IMAGE,
							Env:   envs,
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
			BackoffLimit: &backOffLimit,
		},
	}
	return job
}

func prepareEnvs(ip string, port int, interval int, duration int) []corev1.EnvVar {
	envs := []corev1.EnvVar{}
	envs = append(envs, corev1.EnvVar{Name: ADDRESS, Value: ip})
	envs = append(envs, corev1.EnvVar{Name: PORT, Value: strconv.Itoa(port)})
	envs = append(envs, corev1.EnvVar{Name: INTERVAL, Value: strconv.Itoa(interval)})
	envs = append(envs, corev1.EnvVar{Name: DURATION, Value: strconv.Itoa(duration)})
	return envs
}
