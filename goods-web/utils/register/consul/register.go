package consul

import (
	"fmt"
	"mxshop-api/goods-web/global"

	"github.com/hashicorp/consul/api"
)

type Registry struct {
	Host string
	Port int
}

type RegistryClient interface {
	Register(address string, port int, name string, tags []string, id string)
	DeRegister(serviceId string) error
}

func NewRegistryClient(host string, port int) RegistryClient {
	return &Registry{
		Host: host,
		Port: port,
	}
}

func (r *Registry) Register(address string, port int, name string, tags []string, id string) {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	check := &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d/health", address, port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	//生成注册对象
	registration := &api.AgentServiceRegistration{
		Name:    name,
		ID:      id,
		Port:    port,
		Tags:    tags,
		Address: address,
		Check:   check,
	}

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}
}

func (r *Registry) DeRegister(serviceId string) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d",global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)


	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	err = client.Agent().ServiceDeregister(serviceId)
	return err
}
