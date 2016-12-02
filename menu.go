package summer

type (
	//Menu struct
	Menu struct {
		Title  string
		Parent *Menu
	}
)

var (
	menusList = []*Menu{}

	// RootMenu is zerro-level menu
	RootMenu = Menu{}

	// MainMenu is main admin-panel menu
	MainMenu = RootMenu.Add("Main Menu")

	// DropMenu is top user dropdown menu
	DropMenu = RootMenu.Add("User Menu")
)

// Add submenu to current menu
func (m *Menu) Add(title string) *Menu {
	menu := &Menu{Title: title, Parent: m}
	menusList = append(menusList, menu)
	return menu
}
