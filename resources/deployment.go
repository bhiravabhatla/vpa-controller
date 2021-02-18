package resources

import (
	"context"

	apps "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
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
			deployment.SetFinalizers(append(deployment.GetFinalizers(), finalizer))
			err = d.Client.Update(d.Ctx, &deployment)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func (d *DeploymentRef) HandleDeleteDeployment(c *ConfigMapRef) (bool,error) {

	deployment := apps.Deployment{}
	_, err := d.getObjectIfExists(d.NamespacedName, &deployment)
	if err != nil {
		return false, err
	}

	if !deployment.ObjectMeta.DeletionTimestamp.IsZero() {
		err := c.deleteConfigMapIfExists()
		if err != nil {
			return false,err
		}

		//Get deployment again - to get a fresh revision. https://github.com/operator-framework/operator-sdk/issues/3968
		_, err = d.getObjectIfExists(d.NamespacedName, &deployment)
		if err != nil {
			return false, err
		}
		deployment.SetFinalizers(removeString(deployment.GetFinalizers(), finalizer))
		err = d.Client.Update(context.Background(), &deployment)
		if err != nil {
			return false,err
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
