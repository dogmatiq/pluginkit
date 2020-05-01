package pluginkit_test

import (
	"context"
	"errors"

	. "github.com/dogmatiq/pluginkit"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type pluginStub struct {
	ApplicationService
	CloseFunc func() error
}

func (s *pluginStub) Close() error {
	if s.CloseFunc != nil {
		return s.CloseFunc()
	}

	return nil
}

var _ = Describe("type v1", func() {
	var (
		ctx    context.Context
		stub   *pluginStub
		new    func(context.Context) (interface{}, error)
		plugin Plugin
		file   string
	)

	BeforeEach(func() {
		stub = &pluginStub{}

		new = func(context.Context) (interface{}, error) {
			return stub, nil
		}

		ctx = context.WithValue(
			context.Background(),
			"func",
			func(context.Context) (interface{}, error) {
				return new(ctx)
			},
		)
	})

	JustBeforeEach(func() {
		var err error
		plugin, file, err = loadPlugin(ctx, "v1-stub")
		Expect(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		plugin.Close()
	})

	Describe("func File()", func() {
		It("returns the filename passed to Load()", func() {
			Expect(plugin.File()).To(Equal(file))
		})
	})

	Describe("func ApplicationService()", func() {
		It("returns the plugin implementation as an ApplicationService", func() {
			s, ok := plugin.ApplicationService()
			Expect(s).NotTo(BeNil())
			Expect(ok).To(BeTrue())
		})
	})

	Describe("func Close()", func() {
		When("the plugin implements io.Closer", func() {
			It("calls Close() on the plugin", func() {
				stub.CloseFunc = func() error {
					return errors.New("<error>")
				}

				err := plugin.Close()
				Expect(err).To(MatchError("<error>"))
			})
		})

		When("the plugin does not implement io.Closer", func() {
			BeforeEach(func() {
				new = func(context.Context) (interface{}, error) {
					return "<impl>", nil
				}
			})

			It("returns nil", func() {
				err := plugin.Close()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
