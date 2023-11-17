package resources

import (
	"fmt"
	v1 "github.com/xuxant/voyager-operator/api/v1"
	"github.com/xuxant/voyager-operator/pkg/constants"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewResourceObjectMeta(ship *v1.Ship) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      GetResourceName(ship),
		Namespace: ship.ObjectMeta.Namespace,
		Labels:    BuildResourceLabels(ship),
	}
}

func BuildResourceLabels(ship *v1.Ship) map[string]string {
	return map[string]string{
		constants.LabelAppKey:    constants.LabelAppValue,
		constants.LabelShipCRKey: ship.Name,
	}
}

func GetResourceName(ship *v1.Ship) string {
	return fmt.Sprintf("%s-%s", constants.LabelAppValue, ship.ObjectMeta.Name)
}
