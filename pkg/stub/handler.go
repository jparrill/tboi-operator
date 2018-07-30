package stub

import (
	"context"

	"github.com/jparrill/tboi-operator/pkg/apis/tboi/v1alpha1"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	//	corev1 "k8s.io/api/core/v1"
	//	"k8s.io/apimachinery/pkg/api/errors"
	//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//	"k8s.io/apimachinery/pkg/runtime/schema"
)

func NewHandler() sdk.Handler {
	return &Handler{}
}

type Handler struct {
	// Fill me
}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *v1alpha1.Item:
		object := o
		logrus.Infof("Object: %s", object)
		logrus.Infof("Object: %v", object)
		/*err := sdk.Create(newbusyBoxPod(o))
		if err != nil && !errors.IsAlreadyExists(err) {
			logrus.Errorf("Failed to create busybox pod : %v", err)
			return err
		}*/
		logrus.Infof("Finish Handling")
	}
	return nil
}