package summer

import (
	"sync"
)

type (
	//MenuList struct
	menuList struct {
		sync.Mutex
		list []*Menu
	}
)

func createMenuList() *menuList {
	m := &menuList{}
	m.Mutex = sync.Mutex{}
	m.list = []*Menu{}
	return m
}

// Add one submenu
func (m *menuList) Add(menu *Menu) {
	m.Lock()
	defer m.Unlock()
	m.list = append(m.list, menu)
}

// GetList returns menu list
func (m *menuList) GetList() []*Menu {
	m.Lock()
	defer m.Unlock()
	return append([]*Menu{}, m.list...)
}
