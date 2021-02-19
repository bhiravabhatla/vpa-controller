package resources

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Object struct {
	Ctx            context.Context
	Client         client.Client
	NamespacedName types.NamespacedName
	Log            logr.Logger

	resource runtime.Object
}

func NewObject(ctx context.Context, client client.Client, namespacedName types.NamespacedName, log logr.Logger, object runtime.Object) *Object {
	return &Object{Ctx: ctx, Client: client, NamespacedName: namespacedName, Log: log, resource: object}
}

func (o *Object) PopulateObjectIfExists(namespacedName types.NamespacedName) (bool, error) {

	err := o.Client.Get(o.Ctx, namespacedName, o.resource)
	if errors.IsNotFound(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil

}
