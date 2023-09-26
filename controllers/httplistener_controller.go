/*
Copyright 2023.
*/

package controllers

import (
	"context"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/riverphillips/control-plane/internal/resources"
	"github.com/riverphillips/control-plane/internal/xdscache"
	"k8s.io/apimachinery/pkg/api/errors"
	"math"
	"math/rand"
	"strconv"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	apigatewayv1alpha1 "github.com/riverphillips/control-plane/api/v1alpha1"
)

// HttpListenerReconciler reconciles a HttpListener object
type HttpListenerReconciler struct {
	client.Client
	Scheme          *runtime.Scheme
	xdsCache        xdscache.XDSCache
	cache           cache.SnapshotCache
	nodeId          string
	snapshotVersion int64
}

func NewHttpListenerReconciler(
	client client.Client,
	scheme *runtime.Scheme,
	nodeId string,
	cache cache.SnapshotCache,
) *HttpListenerReconciler {
	return &HttpListenerReconciler{
		Client:          client,
		Scheme:          scheme,
		nodeId:          nodeId,
		cache:           cache,
		snapshotVersion: rand.Int63n(1000),
		xdsCache: xdscache.XDSCache{
			Routes:    make(map[string]resources.Route),
			Listeners: make(map[string]resources.Listener),
			Clusters:  make(map[string]resources.Cluster),
		},
	}
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
	logger := log.FromContext(ctx)

	logger.Info("Reconciling HttpListener")

	httpListener := apigatewayv1alpha1.HttpListener{}
	err := r.Get(ctx, req.NamespacedName, &httpListener)
	if err != nil && !errors.IsNotFound(err) {
		logger.Error(err, "Failed to get HttpListener", "name", req.NamespacedName, "errorMessage", err.Error())
		return ctrl.Result{}, err
	}

	if errors.IsNotFound(err) {
		// Request object not found, could have been deleted after reconcile request.
		// This should be cleared up by the finalizer. If not we will delete the object from the xDS snapshot
		logger.Error(err, "HttpListener resource not found.")

		r.xdsCache.RemoveListener(ctx, httpListener.Name)

		return ctrl.Result{}, nil
	}

	var lrRoutes []string

	uniqueServices := make(map[string]struct{})

	for _, route := range httpListener.Spec.Routes {
		lrRoutes = append(lrRoutes, route.RouteName)
		r.xdsCache.AddRoute(ctx, route.RouteName, route.Prefix, route.Service)

		// We need to track all unique services from the k8s api server and use that to create an envoy cluster
		uniqueServices[route.Service] = struct{}{}
	}

	r.xdsCache.AddListener(ctx, httpListener.Name, lrRoutes, httpListener.Spec.IpAddress, httpListener.Spec.Port)

	for serviceName := range uniqueServices {
		r.xdsCache.AddCluster(ctx, serviceName)
	}

	listeners, err := r.xdsCache.ListenerContents()
	if err != nil {
		logger.Error(err, "Failed to get listener contents")
		return ctrl.Result{}, err
	}

	snapshot, err := cache.NewSnapshot(r.newSnapshotVersion(), map[resource.Type][]types.Resource{
		resource.RouteType:    r.xdsCache.RouteContents(),
		resource.ListenerType: listeners,
		resource.ClusterType:  r.xdsCache.ClusterContents(),
	})

	if err != nil {
		logger.Error(err, "Failed to create snapshot")
		return ctrl.Result{}, err
	}

	if err := snapshot.Consistent(); err != nil {
		logger.Error(err, "Snapshot inconsistency")
		return ctrl.Result{}, err
	}

	if err := r.cache.SetSnapshot(ctx, r.nodeId, snapshot); err != nil {
		logger.Error(err, "Failed to set snapshot")
		return ctrl.Result{}, err
	}

	logger.Info("Snapshot set successfully")

	return ctrl.Result{
		RequeueAfter: time.Second * 30,
	}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HttpListenerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apigatewayv1alpha1.HttpListener{}).
		Complete(r)
}

func (r *HttpListenerReconciler) newSnapshotVersion() string {

	// Reset the snapshotVersion if it ever hits max size.
	if r.snapshotVersion == math.MaxInt64 {
		r.snapshotVersion = 0
	}

	// Increment the snapshot version & return as string.
	r.snapshotVersion++
	return strconv.FormatInt(r.snapshotVersion, 10)
}
