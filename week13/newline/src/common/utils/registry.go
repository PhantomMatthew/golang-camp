package utils

import (
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-plugins/registry/consul/v2"
	"newline.com/newline/src/common/config"
)

func GetRegistry() registry.Registry {
	// Get registry options
	registryType := config.Get("registry.type").String()

	registryUrls := []string{}

	for _, u := range config.Get("registry.urls").Array() {
		registryUrls = append(registryUrls, u.String())
	}

	var reg interface{}

	switch registryType {
	case "consul":
		reg = consul.NewRegistry(func(op *registry.Options) {
			op.Addrs = registryUrls
		})
	case "etcd":
		reg = etcd.NewRegistry(func(op *registry.Options) {
			op.Addrs = registryUrls
		})
	}
	return reg.(registry.Registry)
}
