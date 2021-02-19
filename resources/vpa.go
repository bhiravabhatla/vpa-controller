package resources

import (
	"reflect"

	apps "k8s.io/api/apps/v1"
	v1 "k8s.io/api/autoscaling/v1"
	autoscaler "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta2"
)

type Vpa struct {
	*Object
}

func (vpa *Vpa) CreateOrUpdateVPA(deployment *Deployment) error {


	desiredVPA := autoscaler.VerticalPodAutoscaler{}
	getVPAForDeployment(&desiredVPA, deployment)

	if vpa.name() == "" {

	vpa.Log.WithValues("Object:", "Vpa").Info("Creating vpa " + desiredVPA.Name)
		err := vpa.Client.Create(vpa.Ctx, &desiredVPA)
		if err != nil {
			return err
		}
		return nil
	}
	if !checkIfVPAsEqual(vpa.vpaObj(), &desiredVPA) {
		vpa.Log.WithValues("Object:", "Vpa").Info("Updating vpa " + desiredVPA.Name)
		desiredVPA.SetResourceVersion(vpa.vpaObj().GetResourceVersion())
		err := vpa.Client.Update(vpa.Ctx, &desiredVPA)
		if err != nil {
			return err
		}
	}
	return nil

}

func (vpa *Vpa) deleteVPAIfExists() error {

	if vpa.name() != "" {
		vpa.Log.WithValues("Object:", "Vpa").Info("Deleting vpa " + vpa.name())
		err := vpa.Client.Delete(vpa.Ctx, vpa.resource)
		if err != nil {
			return err
		}
	}
	return nil
}

func (vpa *Vpa) name() string {
	return vpa.resource.(*autoscaler.VerticalPodAutoscaler).Name
}

func (vpa *Vpa) vpaObj() *autoscaler.VerticalPodAutoscaler {
	return vpa.resource.(*autoscaler.VerticalPodAutoscaler)
}

func getVPAForDeployment(vpa *autoscaler.VerticalPodAutoscaler, deployment *Deployment) {

	updateMode := autoscaler.UpdateModeOff
	vpa.TypeMeta.Kind = "VerticalPodAutoscaler"
	vpa.TypeMeta.APIVersion = "v1"
	vpa.SetName(deployment.name() + "-vpa")
	vpa.SetNamespace(deployment.namespace())
	vpa.SetLabels(getVPALabels(deployment.deploymentObj()))
	vpa.SetOwnerReferences(getOwnerReference(deployment.deploymentObj()))

	vpa.Spec = autoscaler.VerticalPodAutoscalerSpec{
		TargetRef:      &v1.CrossVersionObjectReference{
			Kind: deployment.deploymentObj().Kind,
			APIVersion: deployment.deploymentObj().APIVersion,
			Name: deployment.name(),
		},
		UpdatePolicy:   &autoscaler.PodUpdatePolicy{UpdateMode: &updateMode},
		ResourcePolicy: nil,
	}


}

func checkIfVPAsEqual(actual *autoscaler.VerticalPodAutoscaler, desired *autoscaler.VerticalPodAutoscaler) bool {
	if reflect.DeepEqual(actual.Spec.TargetRef, desired.Spec.TargetRef) && reflect.DeepEqual(actual.Spec.UpdatePolicy, desired.Spec.UpdatePolicy) {
		return true
	}

	return false
}

func getVPALabels(deployment *apps.Deployment) map[string]string {
	labels := map[string]string{}

	labels["app.kubernetes.io/for-deployment"] = deployment.Name
	return labels
}
