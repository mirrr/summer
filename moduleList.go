package summer

import (
	"sync"
)

type (
	//ModuleList struct
	ModuleList struct {
		sync.Mutex
		list map[string]Simple
	}
)

func (m *ModuleList) init() {
	m.Mutex = sync.Mutex{}
	m.list = map[string]Simple{}
}

// Get one module by name
func (m *ModuleList) Get(name string) (module Simple, exists bool) {
	m.Lock()
	defer m.Unlock()
	module, exists = m.list[name]
	return
}

// Add one module
func (m *ModuleList) add(module Simple) {
	m.Lock()
	defer m.Unlock()
	m.list[module.GetSettings().Name] = module
}

// GetList one module by name
func (m *ModuleList) GetList() map[string]Simple {
	m.Lock()
	defer m.Unlock()

	ret := map[string]Simple{}
	for name, _ := range m.list {
		ret[name] = m.list[name]
	}
	return ret
}
