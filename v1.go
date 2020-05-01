package pluginkit

import (
	"context"
	"io"
)

// v1symbol is the name of the plugin's "new" function for v1 of the plugin API.
const v1symbol = "NewDogmaPluginV1"

// v1new is the expected signature of the plugin's "new" function for v1 of the
// plugin API.
type v1new = func(ctx context.Context) (interface{}, error)

// v1 is an implementation of Plugin for plugin's that implement v1 of the
// plugin API.
type v1 struct {
	file string
	impl interface{}
}

func (p *v1) File() string {
	return p.file
}

func (p *v1) Services() []string {
	var services []string

	if _, ok := p.impl.(ApplicationService); ok {
		services = append(services, ApplicationServiceName)
	}

	return services
}

func (p *v1) ApplicationService() ApplicationService {
	return p.impl.(ApplicationService)
}

func (p *v1) Close() error {
	if c, ok := p.impl.(io.Closer); ok {
		return c.Close()
	}

	return nil
}
