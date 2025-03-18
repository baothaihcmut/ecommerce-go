package bootstrap

import "sync"

type OnApplicationBootstrap interface {
	RunBootstrap()
}

type ApplicationBootstrapContainer interface {
	Register(OnApplicationBootstrap)
	Run()
}

type ApplicationBootstrapContainerImpl struct {
	component []OnApplicationBootstrap
}

// Register implements ApplicationBootstrapContainer.
func (a *ApplicationBootstrapContainerImpl) Register(svc OnApplicationBootstrap) {
	a.component = append(a.component,svc)
}

// Run implements ApplicationBootstrapContainer.
func (a *ApplicationBootstrapContainerImpl) Run() {
	wg:= sync.WaitGroup{}
	for _, svc:= range a.component{
		wg.Add(1)
		go func() {
			defer wg.Done()
			svc.RunBootstrap()
		}()
	}
	wg.Wait()
}

func NewApplicationBootstrapContainer() ApplicationBootstrapContainer {
	return &ApplicationBootstrapContainerImpl{
		component: make([]OnApplicationBootstrap, 0),
	}
}
