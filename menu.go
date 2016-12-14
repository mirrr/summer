package summer

import (
	"sort"
	"sync"
)

type (
	//Menu struct
	Menu struct {
		Title  string
		Order  int
		Parent *Menu
		Link   string
	}

	menuItem struct {
		Order   int
		Title   string
		Parent  *Menu
		Current *Menu
		Link    string
		SubMenu bool
	}

	menuItems []*menuItem
)

var (
	menusList  = []*Menu{}
	menuListMu = sync.Mutex{}
)

func (slice menuItems) Len() int {
	return len(slice)
}

func (slice menuItems) Less(i, j int) bool {
	if slice[i].Order != slice[j].Order {
		return slice[i].Order < slice[j].Order
	}
	if slice[i].SubMenu != slice[j].SubMenu {
		return slice[i].SubMenu
	}
	return slice[i].Title < slice[j].Title
}

func (slice menuItems) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// Add submenu to current menu
func (m *Menu) Add(title string, order ...int) *Menu {
	menu := &Menu{Title: title, Parent: m}
	if len(order) > 0 {
		menu.Order = order[0]
	}
	menuListMu.Lock()
	menusList = append(menusList, menu)
	menuListMu.Unlock()
	return menu
}

func getMenuItems(m *Menu) interface{} {
	menuItemsList := menuItems{}
	menuListMu.Lock()
	for _, menu := range menusList {
		if menu.Parent == m {
			menuItemsList = append(menuItemsList, &menuItem{
				Order:   menu.Order,
				Title:   menu.Title,
				Parent:  menu.Parent,
				Current: menu,
				Link:    menu.Link,
				SubMenu: len(menu.Link) == 0,
			})
		}
	}
	menuListMu.Unlock()
	modulesListMu.Lock()
	for _, module := range modulesList {
		if module.GetSettings().Menu != nil && module.GetSettings().Menu == m {
			menuItemsList = append(menuItemsList, &menuItem{
				Order:   module.GetSettings().MenuOrder,
				Title:   module.GetSettings().MenuTitle,
				Parent:  module.GetSettings().Menu,
				Link:    "/" + module.GetSettings().PageRouteName + "/",
				SubMenu: false,
			})
		}
	}
	modulesListMu.Unlock()
	sort.Sort(menuItemsList)
	return menuItemsList
}
