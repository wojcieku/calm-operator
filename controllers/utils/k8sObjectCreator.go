package utils

//TODO getService, getDeployment, getJob
import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	NAMESPACE               = "calm-operator-system"
	LABEL_KEY               = "measurement"
	SERVER_IMAGE            = "jwojciech/udp_probe_server"
	CLIENT_IMAGE            = "jwojciech/udp_probe_client"
	NODE_SELECTOR_HOST_NAME = "kubernetes.io/hostname"
)

func PrepareLatencyServerDeployment(deploymentName string, nodeName string, port int, label string) *appsv1.Deployment {
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
						},
					},
					NodeSelector: map[string]string{NODE_SELECTOR_HOST_NAME: nodeName},
				},
			},
		},
	}

	//bind the deployment to this custom resource instance
	//err := ctrl.SetControllerReference(serverStruct, serverDeployment, r.Scheme)

	return deployment
}

func PrepareServiceForLatencyServer(ip string, port int, svcName string, deploymentName string, label string) *corev1.Service {
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
