package summer

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mirrr/mgo-wrapper"
	"gopkg.in/mgo.v2"
)

type (
	// Func is alias for map[string]func(c *gin.Context)
	Func map[string]func(c *gin.Context)

	//Module struct
	Module struct {
		Panel
		Collection *mgo.Collection
		Settings   *ModuleSettings
	}

	//ModuleSettings struct
	ModuleSettings struct {
		Name           string
		Menu           *Menu
		MenuTitle      string
		PageRouteName  string
		AjaxRouteName  string
		Title          string
		CollectionName string
		TemplateName   string
		AllowGroups    []string
		AllowRoles     []string
		AllowUsers     []uint64
		Ajax           Func
	}

	// Simple module interface
	Simple interface {
		Init(settings *ModuleSettings, panel *Panel)
		Page(c *gin.Context)
		Ajax(c *gin.Context)
	}
)

var (
	ModulesList = map[string]Simple{}
)

// Ajax  is default module's ajax method
func (m *Module) Ajax(c *gin.Context) {
	found := false
	for ajaxRoute, ajaxFunc := range m.Settings.Ajax {
		if strings.ToLower(c.Param("method")) == ajaxRoute {
			ajaxFunc(c)
			found = true
			break
		}
	}

	if !found {
		c.String(400, `Method not found in module "`+m.Settings.Name+`"!`)
	}
}

// Page is default module's page rendering method
func (m *Module) Page(c *gin.Context) {
	c.HTML(200, m.Settings.TemplateName+".html", gin.H{
		"title": m.Settings.Title,
		"user":  c.MustGet("user"),
	})
}

// Init is default module's initial method
func (m *Module) Init(settings *ModuleSettings, panel *Panel) {
	m.Settings = settings
	m.Panel = *panel
	if m.Collection == nil {
		m.Collection = mongo.DB(panel.DBName).C(settings.CollectionName)
	}
}
