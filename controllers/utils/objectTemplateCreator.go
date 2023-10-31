package utils

//TODO getService, getDeployment, getJob
import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	NAMESPACE = "calm-operator-system"
	LABEL_KEY = "measurement"
)

func CreateService(ip string, port int, svcName string, deploymentName string, label string) *corev1.Service {
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
