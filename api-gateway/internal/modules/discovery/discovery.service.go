package discovery

import (
	"net/http"

	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/log"
	"github.com/hashicorp/consul/api"
)

type DiscoveryService interface {
	DiscoverService(serviceName string, tag string) (string, int, error)
}

type DiscoveryServiceImpl struct {
	consulClient *api.Client
	logger       *log.Logger
}

func NewDiscoveryService(consulClient *api.Client) DiscoveryService {
	return &DiscoveryServiceImpl{
		consulClient: consulClient,
	}
}
func (d *DiscoveryServiceImpl) DiscoverService(serviceName string, tag string) (string, int, error) {
	services, _, err := d.consulClient.Catalog().Service(serviceName, tag, nil)
	if err != nil {
		level.Error(*d.logger).Log("msg", "Error getting service from consul", "error", err)
		return "", 0, err
	}
	if len(services) == 0 {
		level.Error(*d.logger).Log("msg", "Service not found", "service", serviceName)
		return "", 0, http.ErrNotSupported
	}
	return services[0].ServiceAddress, services[0].ServicePort, nil
}
