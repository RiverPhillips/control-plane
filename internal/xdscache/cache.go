package xdscache

import (
	"context"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/riverphillips/control-plane/internal/resources"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type XDSCache struct {
	Listeners map[string]resources.Listener
	Routes    map[string]resources.Route
	Clusters  map[string]resources.Cluster
}

func (c *XDSCache) AddListener(ctx context.Context, name string, routeNames []string, address string, port uint32) {
	_ = log.FromContext(ctx)

	c.Listeners[name] = resources.Listener{
		Name:    name,
		Routes:  routeNames,
		Address: address,
		Port:    port,
	}
}

func (c *XDSCache) AddRoute(ctx context.Context, name string, prefix string, service string) {
	_ = log.FromContext(ctx)

	c.Routes[name] = resources.Route{
		Name:    name,
		Prefix:  prefix,
		Service: service,
	}

}

func (c *XDSCache) AddCluster(ctx context.Context, name string) {
	_ = log.FromContext(ctx)

	c.Clusters[name] = resources.Cluster{
		Name: name,
	}
}

func (c *XDSCache) RouteContents() []types.Resource {
	var routes []resources.Route

	for _, r := range c.Routes {
		routes = append(routes, r)
	}

	return []types.Resource{resources.MakeRoute(routes)}
}

func (c *XDSCache) RemoveListener(_ context.Context, name string) {
	delete(c.Listeners, name)
}

func (c *XDSCache) ListenerContents() ([]types.Resource, error) {
	var r []types.Resource

	for _, l := range c.Listeners {
		listener, err := resources.MakeHTTPListener(l.Name, l.Routes[0], l.Address, l.Port)
		if err != nil {
			return nil, err
		}
		r = append(r, listener)
	}

	return r, nil
}

func (c *XDSCache) ClusterContents() []types.Resource {
	var r []types.Resource

	for _, cluster := range c.Clusters {
		r = append(r, resources.MakeCluster(cluster.Name))
	}

	return r
}
