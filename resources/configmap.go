package resources

import (
	"reflect"

	apps "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ConfigMapRef struct {
	*ObjectRef
	configMap *v1.ConfigMap
}

func (c *ConfigMapRef) CreateOrUpdateConfigMap() error {
	deploy := apps.Deployment{}
	deployExists, err := c.getObjectIfExists(c.NamespacedName, &deploy)
	if err != nil {
		return err
	}

	if !deployExists {
		c.Log.WithValues("Object:", "ConfigMapRef").Info("Deployment doesn't exist. Skipping creation of configmap")
		return nil
	}

	actualCM := v1.ConfigMap{}
	configMapExists, err := c.getObjectIfExists(client.ObjectKey{Namespace: deploy.Namespace, Name: deploy.Name + "-config"}, &actualCM)
	if err != nil {
		return err
	}

	desiredCM := v1.ConfigMap{}
	getConfigMapForDeployment(&desiredCM, &deploy)

	if !configMapExists {

		c.Log.WithValues("Object:", "ConfigMapRef").Info("Creating Configmap " + desiredCM.Name)
		err = c.Client.Create(c.Ctx, &desiredCM)
		if err != nil {
			return err
		}
		return nil
	}
	if !checkIfCmsEqual(&actualCM, &desiredCM) {
		c.Log.WithValues("Object:", "ConfigMapRef").Info("Updating Configmap " + desiredCM.Name)
		err = c.Client.Update(c.Ctx, &desiredCM)
		if err != nil {
			return err
		}
	}
	return nil

}

func (c *ConfigMapRef) deleteConfigMapIfExists() error {
	cm := v1.ConfigMap{}
	exists, err := c.getObjectIfExists(client.ObjectKey{Namespace: c.NamespacedName.Namespace, Name: c.NamespacedName.Name + "-config"}, &cm)
	if err != nil {
		return err
	}
	if exists {
		c.Log.WithValues("Object:", "ConfigMapRef").Info("Deleting configmap " + cm.Name)
		err = c.Client.Delete(c.Ctx, &cm)
		if err != nil {
			return err
		}
	}
	return nil
}

func getConfigMapForDeployment(configMap *v1.ConfigMap, deployment *apps.Deployment) {

	data := map[string]string{}
	data["test"] = "blah"

	configMap.TypeMeta.Kind = "ConfigMap"
	configMap.TypeMeta.APIVersion = "v1"
	configMap.SetName(deployment.Name + "-config")
	configMap.SetNamespace(deployment.Namespace)
	configMap.SetLabels(getConfigMapLabels(deployment))
	configMap.SetOwnerReferences(getOwnerReference(deployment))
	configMap.Data = data

}

func checkIfCmsEqual(actual *v1.ConfigMap, desired *v1.ConfigMap) bool {
	if reflect.DeepEqual(actual.Data, desired.Data) {
		return true
	}

	return false
}

func getConfigMapLabels(deployment *apps.Deployment) map[string]string {
	labels := map[string]string{}

	labels["app.kubernetes.io/for-deployment"] = deployment.Name
	return labels
}
