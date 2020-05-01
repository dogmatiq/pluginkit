package main

import "context"

func NewDogmaPluginV1(ctx context.Context) (interface{}, error) {
	fn := ctx.Value("func").(func(context.Context) (interface{}, error))
	return fn(ctx)
}

func main() {
}
