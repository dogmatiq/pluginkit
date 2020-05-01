package pluginkit

import (
	"context"
	"fmt"
	"path"
	"plugin"
)

// Plugin is an interface for using services provided by a Dogma plugin.
type Plugin interface {
	// File returns the path to the plugin file, as passed to Load().
	File() string

	// Services returns a list of the services implemented by the plugin.
	Services() []string

	// ApplicationService returns the plugin's ApplicationService.
	//
	// It panics if the plugin does not implement this service.
	ApplicationService() ApplicationService

	// Close closes the plugin.
	Close() error
}

// Load returns a Dogma plugin loaded from a file.
func Load(ctx context.Context, file string) (Plugin, error) {
	base := path.Base(file)

	p, err := plugin.Open(file)
	if err != nil {
		return nil, err
	}

	v, err := p.Lookup(v1symbol)
	if err != nil {
		return nil, fmt.Errorf(
			"%s does not implement any supported Dogma plugin version",
			base,
		)
	}

	new, ok := v.(v1new)
	if !ok {
		return nil, fmt.Errorf(
			"%s is not a valid v1 plugin, %s has type %T, expected %T",
			base,
			v1symbol,
			v,
			v1new(nil),
		)
	}

	impl, err := new(ctx)
	if err != nil {
		return nil, err
	}

	if impl == nil {
		return nil, fmt.Errorf(
			"%s is not a valid v1 plugin, %s() returned nil",
			base,
			v1symbol,
		)
	}

	return &v1{
		file: file,
		impl: impl,
	}, nil
}
