package pluginkit_test

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"

	. "github.com/dogmatiq/pluginkit"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// loadPlugin calls pluginkit.Load() with an operating-system specific file
// extension, as build by the Makefile.
func loadPlugin(ctx context.Context, name string) (Plugin, string, error) {
	dir := fmt.Sprintf(
		"artifacts/build/debug/%s/%s",
		runtime.GOOS,
		runtime.GOARCH,
	)

	if runtime.GOOS == "windows" {
		name += ".dll"
	} else {
		name += ".so"
	}

	name = filepath.Join(dir, name)
	p, err := Load(ctx, name)

	return p, name, err
}

var _ = Describe("func Load()", func() {
	Context("version 1", func() {
		var (
			ctx context.Context
			new func(context.Context) (interface{}, error)
		)

		BeforeEach(func() {
			ctx = context.WithValue(
				context.Background(),
				"func",
				func(context.Context) (interface{}, error) {
					return new(ctx)
				},
			)
		})

		It("returns an error if the new-function returns nil", func() {
			new = func(context.Context) (interface{}, error) {
				return nil, nil
			}

			_, _, err := loadPlugin(ctx, "v1-stub")
			Expect(err).To(MatchError(MatchRegexp(
				`^v1-stub\.(so|dll) is not a valid v1 plugin, NewDogmaPluginV1\(\) returned nil$`,
			)))
		})

		It("returns an error if the new-function returns an error", func() {
			new = func(context.Context) (interface{}, error) {
				return nil, errors.New("<error>")
			}

			_, _, err := loadPlugin(ctx, "v1-stub")
			Expect(err).To(MatchError("<error>"))
		})

		It("returns an error if the NewDogmaPluginV1 symbol does not have the expected type", func() {
			_, _, err := loadPlugin(ctx, "v1-wrong-type")
			Expect(err).To(MatchError(MatchRegexp(
				`^v1-wrong-type\.(so|dll) is not a valid v1 plugin, NewDogmaPluginV1 has type \*string, expected func\(context.Context\) \(interface \{\}, error\)$`,
			)))
		})
	})

	It("returns an error if the file is not a valid Go plugin", func() {
		_, err := Load(context.Background(), "/dev/null")
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the file is not a Dogma plugin", func() {
		_, _, err := loadPlugin(context.Background(), "nonplugin")
		Expect(err).To(MatchError(MatchRegexp(
			`^nonplugin\.(so|dll) does not implement any supported Dogma plugin version$`,
		)))
	})
})
