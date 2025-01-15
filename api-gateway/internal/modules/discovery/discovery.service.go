package discovery

import (
	"net/http"

	"github.com/hashicorp/consul/api"
)

type DiscoveryService interface {
	DiscoverService(serviceName string, tag string) (string, int, error)
}

type DiscoveryServiceImpl struct {
	consulClient *api.Client
}

func NewDiscoveryService(consulClient *api.Client) DiscoveryService {
	return &DiscoveryServiceImpl{
		consulClient: consulClient,
	}
}
func (d *DiscoveryServiceImpl) DiscoverService(serviceName string, tag string) (string, int, error) {
	services, _, err := d.consulClient.Catalog().Service(serviceName, tag, nil)
	if err != nil {
		return "", 0, err
	}
	if len(services) == 0 {
		return "", 0, http.ErrNotSupported
	}
	return services[0].ServiceAddress, services[0].ServicePort, nil
}
