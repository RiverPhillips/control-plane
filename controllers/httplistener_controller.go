/*
Copyright 2023.
*/

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	apigatewayv1alpha1 "github.com/riverphillips/control-plane/api/v1alpha1"
)

// HttpListenerReconciler reconciles a HttpListener object
type HttpListenerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=apigateway.riverphillips.dev,resources=httplisteners,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apigateway.riverphillips.dev,resources=httplisteners/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apigateway.riverphillips.dev,resources=httplisteners/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the HttpListener object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *HttpListenerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HttpListenerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apigatewayv1alpha1.HttpListener{}).
		Complete(r)
}
