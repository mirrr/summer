package summer

import (
	"sort"
)

type (
	tab struct {
		Order  int
		Title  string
		Link   string
		Active bool
	}

	tabs []*tab
)

func (slice tabs) Len() int {
	return len(slice)
}

func (slice tabs) Less(i, j int) bool {
	if slice[i].Order != slice[j].Order {
		return slice[i].Order < slice[j].Order
	}
	return slice[i].Title < slice[j].Title
}

func (slice tabs) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func getTabs(name string) interface{} {
	modulesListMu.Lock()
	tabsList := tabs{}
	if modulesList[name] == nil {
		return obj{"title": name, "icon": "", "list": tabsList}
	}

	current := modulesList[name]
	parent := current
	if current.GetSettings().GroupTo != nil {
		parent = current.GetSettings().GroupTo
	}

	for _, module := range modulesList {
		if module.GetSettings().GroupTo == parent {
			tabsList = append(tabsList, &tab{
				Order:  module.GetSettings().MenuOrder,
				Title:  module.GetSettings().GroupTitle,
				Link:   "/" + module.GetSettings().PageRouteName + "/",
				Active: module == current,
			})
		}
	}
	modulesListMu.Unlock()
	sort.Sort(tabsList)
	if len(tabsList) > 0 {
		tabsList = append(tabs{&tab{
			Order:  parent.GetSettings().MenuOrder,
			Title:  parent.GetSettings().GroupTitle,
			Link:   "/" + parent.GetSettings().PageRouteName + "/",
			Active: parent == current,
		}}, tabsList...)
	}
	return obj{"title": parent.GetSettings().Title, "icon": parent.GetSettings().Icon, "list": tabsList}
}
