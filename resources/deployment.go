package resources

import (
	apps "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeploymentRef struct {
	*ObjectRef
}

func (d *DeploymentRef) SkipVPAAnnotationExists() (bool, error) {

	deployment := apps.Deployment{}
	err := d.Client.Get(d.Ctx, d.NamespacedName, &deployment)
	if errors.IsNotFound(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	val, ok := deployment.Annotations["thoughtworks.org/skip-cm"]
	if ok {
		if val == "true" {
			return true, err
		}
	}
	return false, err
}

func (d *DeploymentRef) CheckAndAddFinalizer() error {

	deployment := apps.Deployment{}
	exists, err := d.getObjectIfExists(d.NamespacedName, &deployment)

	if err != nil {
		return err
	}
	if exists {
		if !containsString(deployment.GetFinalizers(), finalizer) {
			d.Log.WithValues("Object:", "DeploymentRef").Info("Adding Finalizers.")
			deployment.SetFinalizers(append(deployment.GetFinalizers(), finalizer))
			err = d.Client.Update(d.Ctx, &deployment)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func (d *DeploymentRef) HandleDeleteDeployment(c *ConfigMapRef) (bool, error) {

	deployment := apps.Deployment{}
	_, err := d.getObjectIfExists(d.NamespacedName, &deployment)
	if err != nil {
		return false, err
	}

	if !deployment.ObjectMeta.DeletionTimestamp.IsZero() {
		err := c.deleteConfigMapIfExists()
		if err != nil {
			return false, err
		}
		d.Log.WithValues("Object:", "DeploymentRef").Info("Removing Finalizers for deployment - " + d.NamespacedName.String())
		deployment.SetFinalizers(removeString(deployment.GetFinalizers(), finalizer))
		err = d.Client.Update(d.Ctx, &deployment)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
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
