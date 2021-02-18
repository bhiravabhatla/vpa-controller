package resources

import (
	"reflect"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ConfigMapRef struct {
	*ObjectRef
}

func (c *ConfigMapRef) CreateOrUpdateConfigMap() error {
	actualCM := v1.ConfigMap{}
	exists, err := c.getObjectIfExists(client.ObjectKey{Namespace: "default", Name: "test"}, &actualCM)
	if err != nil {
		return err
	}
	data := map[string]string{}
	data["test"] = "blah"

	desiredCM := v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "default",
		},
		Data: data,
	}

	if !exists {

		c.Log.WithValues("ObjectRef:", "ConfigMapRef").Info("Creating Configmap test")
		err = c.Client.Create(c.Ctx, &desiredCM)
		if err != nil {
			return err
		}
		return nil
	}
	if !checkIfCmsEqual(&actualCM, &desiredCM) {
		c.Log.WithValues("ObjectRef:", "ConfigMapRef").Info("Updating Configmap test")
		err = c.Client.Update(c.Ctx, &desiredCM)
		if err != nil {
			return err
		}
	}
	return nil

}

func (c *ConfigMapRef) deleteConfigMapIfExists() error {
	cm := v1.ConfigMap{}
	_, err := c.getObjectIfExists(client.ObjectKey{Namespace: "default", Name: "test"}, &cm)
	if err != nil {
		return err
	}
	err = c.Client.Delete(c.Ctx, &cm)
	if err != nil {
		return err
	}
	return nil
}

func checkIfCmsEqual(actual *v1.ConfigMap, desired *v1.ConfigMap) bool {
	if reflect.DeepEqual(actual.Data, desired.Data) {
		return true
	}

	return false
}
