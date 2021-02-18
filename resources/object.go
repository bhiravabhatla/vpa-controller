package resources

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ObjectRef struct {
	Ctx            context.Context
	Client         client.Client
	NamespacedName types.NamespacedName
	Log            logr.Logger
}

func NewService(ctx context.Context, client client.Client, namespacedName types.NamespacedName, log logr.Logger) *ObjectRef {
	return &ObjectRef{Ctx: ctx, Client: client, NamespacedName: namespacedName, Log: log}
}

func (s *ObjectRef) getObjectIfExists(objKey client.ObjectKey, obj runtime.Object) (bool, error) {

	err := s.Client.Get(s.Ctx, objKey, obj)
	if errors.IsNotFound(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil

}
