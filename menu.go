package summer

import (
	"sort"
)

type (
	//Menu struct
	Menu struct {
		panel  *Panel
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

func (m *Menu) init(panel *Panel, parent *Menu) {
	m.panel = panel
	m.Parent = parent
}

// Add submenu to current menu
func (m *Menu) Add(title string, order ...int) *Menu {
	menu := &Menu{Title: title, Parent: m, panel: m.panel}
	if len(order) > 0 {
		menu.Order = order[0]
	}
	m.panel.menuList.Add(menu)
	return menu
}

func getMenuItems(panel *Panel, m *Menu, u *UsersStruct) menuItems {
	userActions := uniqAppend(panel.Groups.Get(u.Rights.Groups...), u.Rights.Actions)

	menuItemsList := menuItems{}

	for _, menu := range m.panel.menuList.GetList() {
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

	for _, module := range panel.Modules.GetList() {
		sett := module.GetSettings()
		msr := sett.Rights
		rightsEmpty := len(msr.Groups) == 0 && len(msr.Actions) == 0
		allow := (len(msr.Groups) > 0 && isOverlap(u.Rights.Groups, msr.Groups)) || (len(msr.Actions) > 0 && isOverlap(userActions, msr.Actions))

		if sett.Menu == m && (rightsEmpty || allow) {
			menuItemsList = append(menuItemsList, &menuItem{
				Order:   sett.MenuOrder,
				Title:   sett.MenuTitle,
				Parent:  sett.Menu,
				Link:    "/" + module.GetSettings().PageRouteName + "/",
				SubMenu: false,
			})
		}
	}

	sort.Sort(menuItemsList)
	return menuItemsList
}
