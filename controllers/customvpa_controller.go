package controllers

import (
	"context"

	"github.com/go-logr/logr"
	apps "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/bhiravabhatla/vpa-controller/resources"
)

// CustomVpaReconciler reconciles a CustomVpa object
type CustomVpaReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=vpaextensions.thoughtworks.org,resources=customvpas,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=vpaextensions.thoughtworks.org,resources=customvpas/status,verbs=get;update;patch

func (r *CustomVpaReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {

	ctx := context.Background()
	log := r.Log.WithValues("Event triggered for:", req.NamespacedName)

	service := resources.NewService(ctx, r.Client, req.NamespacedName, log)

	deployref := resources.DeploymentRef{ObjectRef: service}
	cmref := resources.ConfigMapRef{ObjectRef: service}

	annotationExists, err := deployref.SkipVPAAnnotationExists()
	if err != nil {
		return ctrl.Result{}, err
	}
	//If skip annotation is present do nothing.
	if annotationExists {
		r.Log.Info("Skipping create/update cm for " + req.NamespacedName.String())
		return ctrl.Result{}, nil
	}
	err = deployref.CheckAndAddFinalizer()
	if err != nil {
		return ctrl.Result{}, err
	}
	isDeleted, err := deployref.HandleDeleteDeployment(&cmref)
	if err != nil {
		return ctrl.Result{}, err
	}
	if isDeleted {
		return ctrl.Result{}, nil
	}
	err = cmref.CreateOrUpdateConfigMap()
	if err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *CustomVpaReconciler) SetupWithManager(mgr ctrl.Manager) error {

	vpaController, err := controller.New("vpa-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}
	err = vpaController.Watch(&source.Kind{Type: &apps.Deployment{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}
	err = vpaController.Watch(&source.Kind{Type: &v1.ConfigMap{}}, &handler.EnqueueRequestForOwner{OwnerType: &apps.Deployment{}, IsController: true})
	if err != nil {
		return err
	}

	return nil
}
