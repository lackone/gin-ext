package gin

import "github.com/lackone/gin-ext/framework"

func (engine *Engine) Bind(provider framework.ServiceProvider) error {
	return engine.container.Bind(provider)
}

func (engine *Engine) IsBind(key string) bool {
	return engine.container.IsBind(key)
}

func (engine *Engine) SetContainer(container framework.Container) {
	engine.container = container
}

func (engine *Engine) GetContainer() framework.Container {
	return engine.container
}

func (c *Context) Make(key string) (interface{}, error) {
	return c.container.Make(key)
}

func (c *Context) MustMake(key string) interface{} {
	return c.container.MustMake(key)
}

func (c *Context) MakeNew(key string, params []interface{}) (interface{}, error) {
	return c.container.MakeNew(key, params)
}
