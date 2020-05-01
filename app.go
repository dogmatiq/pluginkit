package pluginkit

import (
	"context"
	"io"

	"github.com/dogmatiq/dogma"
)

// ApplicationService is a plugin service that can instantiate Dogma
// applications.
type ApplicationService interface {
	// ApplicationKeys returns the identity keys of the applications that can be
	// created via this service.
	ApplicationKeys() []string

	// NewApplication returns a new instance of an application.
	//
	// k is the application's identity key. It returns an error if k is not one
	// of the keys returned by ApplicationKeys().
	//
	// The returned io.Closer, if non-nil, is used to free any resources
	// allocated for the application, such as closing database connections.
	NewApplication(ctx context.Context, k string) (dogma.Application, io.Closer, error)
}
