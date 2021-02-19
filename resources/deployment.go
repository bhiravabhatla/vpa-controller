package resources

import (
	apps "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Deployment struct {
	*Object
}

func (d *Deployment) SkipVPAAnnotationExists() bool {

	if d.name() != "" {
		val, ok := d.deploymentObj().Annotations["thoughtworks.org/skip-vpa"]
		if ok {
			if val == "true" {
				return true
			}
		}
	}
	return false
}

func (d *Deployment) CheckAndAddFinalizer() error {

	if d.name() != "" {
		if !containsString(d.resource.(*apps.Deployment).GetFinalizers(), finalizer) {
			d.Log.WithValues("Object:", "Deployment: "+d.NamespacedName.String()).Info("Adding Finalizers")
			d.resource.(*apps.Deployment).SetFinalizers(append(d.resource.(*apps.Deployment).GetFinalizers(), finalizer))
			err := d.Client.Update(d.Ctx, d.resource)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func (d *Deployment) HandleDeleteDeployment(v *Vpa) (bool, error) {

	if d.isDeleted() {
		err := v.deleteVPAIfExists()
		if err != nil {
			return false, err
		}
		d.Log.WithValues("Object:", "Deployment").Info("Removing Finalizers for deployment - " + d.name())
		d.deploymentObj().SetFinalizers(removeString(d.deploymentObj().GetFinalizers(), finalizer))
		err = d.Client.Update(d.Ctx, d.resource)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func (d *Deployment) isDeleted() bool {
	return !d.resource.(*apps.Deployment).DeletionTimestamp.IsZero()
}

func (d *Deployment) name() string {
	return d.resource.(*apps.Deployment).Name
}

func (d *Deployment) namespace() string {
	return d.resource.(*apps.Deployment).Namespace
}

func (d *Deployment) deploymentObj() *apps.Deployment {
	return d.resource.(*apps.Deployment)
}

func containsString(slice []string, str string) bool {

	for _, s := range slice {
		if str == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}
func getOwnerReference(deployment *apps.Deployment) []metav1.OwnerReference {
	return []metav1.OwnerReference{*metav1.NewControllerRef(deployment, deployment.GroupVersionKind())}
}
