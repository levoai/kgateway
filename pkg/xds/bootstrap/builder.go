package bootstrap

import (
	"fmt"

	envoybootstrapv3 "github.com/envoyproxy/go-control-plane/envoy/config/bootstrap/v3"
	envoyclusterv3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoycorev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoylistenerv3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	envoyroutev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoy_extensions_filters_network_http_connection_manager_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	envoywellknown "github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"google.golang.org/protobuf/proto"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils"
	"github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/ir"
)

// ConfigBuilder helps construct a partial bootstrap config for validation.
type ConfigBuilder struct {
	filterConfigs       ir.TypedFilterConfigMap
	routes              []*envoyroutev3.Route
	routeConfigurations map[string]*envoyroutev3.RouteConfiguration
	clusters            []*envoyclusterv3.Cluster
	listeners           []*envoylistenerv3.Listener
}

// New creates a new ConfigBuilder.
func New() *ConfigBuilder {
	return &ConfigBuilder{
		filterConfigs:       make(ir.TypedFilterConfigMap),
		routeConfigurations: make(map[string]*envoyroutev3.RouteConfiguration),
	}
}

// AddFilterConfig adds a filter configuration to the builder. Assumes that the
// filter config is a valid proto message and error handling is done by the caller.
func (b *ConfigBuilder) AddFilterConfig(name string, config proto.Message) {
	b.filterConfigs.AddTypedConfig(name, config)
}

// AddRoute adds a route to the builder.
func (b *ConfigBuilder) AddRoutes(routes ...*envoyroutev3.Route) {
	for _, route := range routes {
		b.routes = append(b.routes, proto.Clone(route).(*envoyroutev3.Route))
	}
}

// AddRoute adds a route to the builder.
func (b *ConfigBuilder) AddRouteConfigurations(routeConfigurations ...*envoyroutev3.RouteConfiguration) {
	for _, rc := range routeConfigurations {
		b.routeConfigurations[rc.Name] = proto.Clone(rc).(*envoyroutev3.RouteConfiguration)
	}
}

// AddCluster adds a cluster to the builder.
func (b *ConfigBuilder) AddClusters(clusters ...*envoyclusterv3.Cluster) {
	for _, cluster := range clusters {
		b.clusters = append(b.clusters, proto.Clone(cluster).(*envoyclusterv3.Cluster))
	}
}

// AddCluster adds a listener to the builder.
func (b *ConfigBuilder) AddListener(listeners ...*envoylistenerv3.Listener) {
	for _, listener := range listeners {
		b.listeners = append(b.listeners, proto.Clone(listener).(*envoylistenerv3.Listener))
	}
}

// Build creates a partial bootstrap config suitable for validation.
func (b *ConfigBuilder) Build() (*envoybootstrapv3.Bootstrap, error) {
	staticResources := &envoybootstrapv3.Bootstrap_StaticResources{}
	if len(b.filterConfigs) > 0 || len(b.routes) > 0 {
		vhost := &envoyroutev3.VirtualHost{
			Name:    "placeholder_vhost",
			Domains: []string{"*"},
		}
		if len(b.filterConfigs) > 0 {
			vhost.TypedPerFilterConfig = b.filterConfigs.ToAnyMap()
		}
		if len(b.routes) > 0 {
			vhost.Routes = b.routes
		}

		hcmAny, err := utils.MessageToAny(&envoy_extensions_filters_network_http_connection_manager_v3.HttpConnectionManager{
			StatPrefix: "placeholder",
			RouteSpecifier: &envoy_extensions_filters_network_http_connection_manager_v3.HttpConnectionManager_RouteConfig{
				RouteConfig: &envoyroutev3.RouteConfiguration{
					VirtualHosts: []*envoyroutev3.VirtualHost{vhost},
				},
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to marshal HttpConnectionManager: %w", err)
		}

		staticResources = &envoybootstrapv3.Bootstrap_StaticResources{
			Listeners: []*envoylistenerv3.Listener{{
				Name: "placeholder_listener",
				Address: &envoycorev3.Address{
					Address: &envoycorev3.Address_SocketAddress{
						SocketAddress: &envoycorev3.SocketAddress{
							Address:       "0.0.0.0",
							PortSpecifier: &envoycorev3.SocketAddress_PortValue{PortValue: 8081},
						},
					},
				},
				FilterChains: []*envoylistenerv3.FilterChain{{
					Name: "placeholder_filter_chain",
					Filters: []*envoylistenerv3.Filter{{
						Name: envoywellknown.HTTPConnectionManager,
						ConfigType: &envoylistenerv3.Filter_TypedConfig{
							TypedConfig: hcmAny,
						},
					}},
				}},
			}},
		}
	}
	if len(b.clusters) > 0 {
		staticResources.Clusters = b.clusters
	}
	if len(b.listeners) > 0 {
		staticResources.Listeners = b.listeners
		for _, l := range b.listeners {
			for _, fc := range l.FilterChains {
				for i := range len(fc.Filters) {
					f := fc.Filters[i]
					if f.Name == "envoy.filters.network.http_connection_manager" {
						if _, ok := b.routeConfigurations[fc.Name]; !ok {
							return nil, fmt.Errorf("no RouteConfiguration found with name %s", fc.Name)
						}
						newHCM, err := utils.MessageToAny(&envoy_extensions_filters_network_http_connection_manager_v3.HttpConnectionManager{
							StatPrefix: "placeholder",
							RouteSpecifier: &envoy_extensions_filters_network_http_connection_manager_v3.HttpConnectionManager_RouteConfig{
								RouteConfig: b.routeConfigurations[fc.Name],
							},
						})
						if err != nil {
							return nil, fmt.Errorf("failed to marshal HttpConnectionManager: %w", err)
						}

						fc.Filters[i] = &envoylistenerv3.Filter{
							Name: envoywellknown.HTTPConnectionManager,
							ConfigType: &envoylistenerv3.Filter_TypedConfig{
								TypedConfig: newHCM,
							},
						}
					}
				}
			}
		}
	}

	return &envoybootstrapv3.Bootstrap{
		Node: &envoycorev3.Node{
			Id:      "validation-node-id",
			Cluster: "validation-cluster",
		},
		StaticResources: staticResources,
	}, nil
}
