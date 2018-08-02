package stub

import (
	"context"
	"fmt"
	"github.com/jparrill/tboi-operator/pkg/apis/tboi/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	extensionsbetav1 "k8s.io/api/extensions/v1beta1"
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

		// Create the route if it doesn't exist
		ing_cont := routeTboiItems(svc, tboi)
		err = sdk.Create(ing_cont)
		if err != nil && !apierrors.IsAlreadyExists(err) {
			logrus.Errorf("Failed to create route : %v", err)
			return err
		}

		// Ensure the deployment size is the same as the spec ItemSize object
		err = checkTboiReplicas(dc, tboi)
		if err != nil {
			logrus.Errorf("Error checking/updating replicas : %v", err)
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

func checkTboiReplicas(dc *appsv1.Deployment, tboi *v1alpha1.Item) error {
	// Validate that the dc exists
	err := sdk.Get(dc)
	if err != nil {
		logrus.Errorf("Failed to get deployment : %v", err)
		return err
	}

	// Extract ItemSize object from CR
	ItemSize := tboi.Spec.ItemSize
	logrus.Infof("CR ItemSize: %d, DC Replicas: %d", ItemSize, *dc.Spec.Replicas)

	// If not equal, update the DC with CR spec
	if *dc.Spec.Replicas != ItemSize {
		logrus.Infof("Need to update replicas from %d to %d", *dc.Spec.Replicas, ItemSize)
		dc.Spec.Replicas = &ItemSize
		err = sdk.Update(dc)
		if err != nil {
			logrus.Errorf("Failed to update deployment : %v", err)
			return err
		}
	}

	return err
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
	logrus.Debugf("DC: %s", dc)
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
	logrus.Debugf("SVC: %s", svc)
	logrus.Infof("SVC Spec Finished")
	return svc
}

func routeTboiItems(svc *corev1.Service, h *v1alpha1.Item) *extensionsbetav1.Ingress {
	logrus.Infof("Making Ingress Spec")
	Name := svc.ObjectMeta.Name
	Namespace := svc.ObjectMeta.Namespace
	FullRoute := Name + "-" + Namespace + "." + h.Spec.Route.RouteDomain
	logrus.Infof("Route: %s", h.Spec.Route)

	ing_cont := &extensionsbetav1.Ingress{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "extensions/v1beta1",
			Kind:       "Ingress",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      svc.ObjectMeta.Name,
			Namespace: svc.ObjectMeta.Namespace,
		},
		Spec: extensionsbetav1.IngressSpec{
			Rules: []extensionsbetav1.IngressRule{{
				Host: FullRoute,
				IngressRuleValue: extensionsbetav1.IngressRuleValue{
					HTTP: &extensionsbetav1.HTTPIngressRuleValue{
						Paths: []extensionsbetav1.HTTPIngressPath{{
							Path: h.Spec.Route.RoutePath,
							Backend: extensionsbetav1.IngressBackend{
								ServiceName: svc.ObjectMeta.Name,
								ServicePort: intstr.IntOrString{Type: intstr.Int, IntVal: svc.Spec.Ports[0].Port},
							},
						}},
					},
				},
			}},
		},
	}

	logrus.Infof("IngressController: %v", ing_cont)
	logrus.Infof("Ingress Spec Finished")

	return ing_cont

}
