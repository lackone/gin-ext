package framework

import (
	"errors"
	"sync"
)

// Container 是一个服务容器，提供绑定服务和获取服务的功能
type Container interface {
	// Bind 绑定一个服务提供者，如果关键字凭证已经存在，会进行替换操作，返回error
	Bind(provider ServiceProvider) error
	// IsBind 关键字凭证是否已经绑定服务提供者
	IsBind(key string) bool

	// Make 根据关键字凭证获取一个服务，
	Make(key string) (interface{}, error)
	// MustMake 根据关键字凭证获取一个服务，如果这个关键字凭证未绑定服务提供者，那么会panic。
	// 所以在使用这个接口的时候请保证服务容器已经为这个关键字凭证绑定了服务提供者。
	MustMake(key string) interface{}
	// MakeNew 根据关键字凭证获取一个服务，只是这个服务并不是单例模式的
	// 它是根据服务提供者注册的启动函数和传递的params参数实例化出来的
	// 这个函数在需要为不同参数启动不同实例的时候非常有用
	MakeNew(key string, params []interface{}) (interface{}, error)
}

// ExtContainer 是服务容器的具体实现
type ExtContainer struct {
	Container //强制要求ExtContainer实现Container接口
	// providers 存储注册的服务提供者，key为字符串凭证
	providers map[string]ServiceProvider
	// instance 存储具体的实例，key为字符串凭证
	instances map[string]interface{}
	// lock 用于锁住对容器的变更操作
	lock sync.RWMutex
}

func NewExtContainer() *ExtContainer {
	return &ExtContainer{
		providers: make(map[string]ServiceProvider),
		instances: make(map[string]interface{}),
		lock:      sync.RWMutex{},
	}
}

func (c *ExtContainer) Bind(provider ServiceProvider) error {
	c.lock.Lock()
	name := provider.Name()
	c.providers[name] = provider
	c.lock.Unlock()

	//如果不延迟实例化
	if provider.IsDefer() == false {
		if err := provider.Boot(c); err != nil {
			return err
		}
		params := provider.Params(c)
		method := provider.Register(c)
		instance, err := method(params...)
		if err != nil {
			return errors.New(err.Error())
		}
		c.instances[name] = instance
	}

	return nil
}

func (c *ExtContainer) IsBind(key string) bool {
	return c.findProvider(key) != nil
}

func (c *ExtContainer) Make(key string) (interface{}, error) {
	return c.make(key, nil, false)
}

func (c *ExtContainer) MakeNew(key string, params []interface{}) (interface{}, error) {
	return c.make(key, params, true)
}

func (c *ExtContainer) MustMake(key string) interface{} {
	serv, err := c.make(key, nil, false)
	if err != nil {
		panic("container not contain key " + key)
	}
	return serv
}

func (c *ExtContainer) make(key string, params []interface{}, isNew bool) (interface{}, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	p := c.findProvider(key)
	if p == nil {
		return nil, errors.New(key + " provider not register")
	}

	if isNew {
		return c.newInstance(p, params)
	}

	//如果容器中已经实例化了，那么直接取容器中的实例
	if ins, ok := c.instances[key]; ok {
		return ins, nil
	}

	//容器中还未实例化，则进行一次实例化
	ins, err := c.newInstance(p, nil)
	if err != nil {
		return nil, err
	}

	c.instances[key] = ins

	return ins, nil
}

func (c *ExtContainer) newInstance(p ServiceProvider, params []interface{}) (interface{}, error) {
	if err := p.Boot(c); err != nil {
		return nil, err
	}
	if params == nil {
		params = p.Params(c)
	}
	method := p.Register(c)
	ins, err := method(params...)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return ins, nil
}

func (c *ExtContainer) findProvider(key string) ServiceProvider {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if p, ok := c.providers[key]; ok {
		return p
	}
	return nil
}
