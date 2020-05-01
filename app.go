package pluginkit

import (
	"context"
	"errors"
	"io"

	"github.com/dogmatiq/dogma"
)

// ErrUnknownApplication indicates that NewApplication() was called with an
// unrecognised application key.
var ErrUnknownApplication = errors.New("unknown application")

// ApplicationServiceName is the name of the application service, as returned by
// Plugin.Services().
const ApplicationServiceName = "application"

// ApplicationService is a plugin service that can instantiate Dogma
// applications.
type ApplicationService interface {
	// ApplicationKeys returns the identity keys of the applications that can be
	// created via this service.
	ApplicationKeys() []string

	// NewApplication returns a new instance of an application.
	//
	// k is the application's identity key. If it is unrecocgnized
	// ErrUnknownApplication is returned.
	//
	// The returned io.Closer is used to free any resources allocated for the
	// application, such as closing database connections.
	NewApplication(ctx context.Context, k string) (dogma.Application, io.Closer, error)
}
