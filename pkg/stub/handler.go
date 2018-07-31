package stub

import (
	"context"
	"fmt"
	"github.com/jparrill/tboi-operator/pkg/apis/tboi/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	//schema "k8s.io/apimachinery/pkg/runtime/schema"
)

func NewHandler() sdk.Handler {
	return &Handler{}
}

type Handler struct{}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	logrus.Infof("Start Handling")
	switch o := event.Object.(type) {
	case *v1alpha1.Item:
		tboi := o
		logrus.Infof("Switch Object: %s", tboi)
		logrus.Infof("Handler: %s", h)

		// Create the deployment if it doesn't exist
		dc := dcTboiItems(tboi)
		err := sdk.Create(dc)
		if err != nil && !apierrors.IsAlreadyExists(err) {
			return fmt.Errorf("failed to create deployment: %v", err)
		}

		// Create the svc if it doesn't exist
		svc := svcTboiItems(tboi)
		err = sdk.Create(svc)
		if err != nil && !apierrors.IsAlreadyExists(err) {
			logrus.Errorf("Failed to create service : %v", err)
			return err
		}

		logrus.Infof("Finish Handling")
	}
	return nil
}

func getPodLabels(name string) map[string]string {
	app := "tboi-items-app"
	logrus.Debug("Returning labels")
	return map[string]string{"app": app, "name": name}
}

func dcTboiItems(h *v1alpha1.Item) *appsv1.Deployment {
	logrus.Infof("Making dc spec")
	labels := getPodLabels(h.Name)
	replicas := h.Spec.ItemSize
	dc := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      h.Name,
			Namespace: h.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: "docker.io/padajuan/tboi-operator-app",
						Name:  labels["app"],
						LivenessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/",
									Port: intstr.IntOrString{Type: intstr.Int, IntVal: 5000},
								},
							},
							InitialDelaySeconds: 3,
							PeriodSeconds:       3,
							TimeoutSeconds:      3,
							FailureThreshold:    3,
							SuccessThreshold:    1,
						},
						ReadinessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/",
									Port: intstr.IntOrString{Type: intstr.Int, IntVal: 5000},
								},
							},
							InitialDelaySeconds: 3,
							PeriodSeconds:       3,
							TimeoutSeconds:      3,
							FailureThreshold:    3,
							SuccessThreshold:    3,
						},

						Ports: []corev1.ContainerPort{{
							ContainerPort: 5000,
							Name:          labels["app"],
						}},
					}},
				},
			},
		},
	}
	logrus.Infof("DC: %s", dc)
	logrus.Infof("DC Spec Finished")
	return dc
}

func svcTboiItems(h *v1alpha1.Item) *corev1.Service {
	logrus.Infof("Making svc Spec")
	labels := getPodLabels(h.Name)
	svc := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      h.Name,
			Namespace: h.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeLoadBalancer,
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Name: "http",
					Port: 5000,
				},
			},
		},
	}
	logrus.Infof("SVC: %s", svc)
	logrus.Infof("SVC Spec Finished")
	return svc
}
