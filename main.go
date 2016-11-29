package summer

import (
	"fmt"
	"reflect"
	"text/template"
	"time"
	"ttpl"

	"github.com/gin-gonic/gin"
	"github.com/mirrr/mgo-ai"
	"github.com/mirrr/mgo-wrapper"
	"gopkg.in/mgo.v2"
	"gopkg.in/mirrr/types.v1"
)

type (
	Func map[string]func(c *gin.Context)

	// Settings intended for data transmission into the Init method of package
	Settings struct {
		Port        uint
		Title       string
		AuthSalt    string
		AuthPrefix  string
		Path        string // URL path of panel - "/" by default
		Views       string // file path of ./templates directory
		ViewsDoT    string // file path of doT.js templates directory
		Files       string // file path of ./files directory
		TMPs        string // file path of /tmp directory
		DBName      string // MongoDB database name
		DefaultPage string
		Language    string
		Vars        map[string]interface{}
		TFuncMap    template.FuncMap
		FirstStart  func()
		RouterGroup *gin.RouterGroup
		Engine      *gin.Engine
	}

	//Panel ...
	Panel struct {
		Settings
	}

	//Module ...
	Module struct {
		Panel
		Collection *mgo.Collection
		Settings   *ModuleSettings
	}

	//ModuleSettings ...
	ModuleSettings struct {
		Name           string
		PageRouteName  string
		AjaxRouteName  string
		Title          string
		CollectionName string
		AllowGroups    []string
		AllowRoles     []string
		Ajax           Func
	}

	// Simple module interface
	Simple interface {
		Init(settings *ModuleSettings, panel *Panel)
		Page(c *gin.Context)
		Ajax(c *gin.Context)
	}
)

func (panel *Panel) AddModule(settings *ModuleSettings, s Simple) Simple {
	st := reflect.ValueOf(s)
	for i := 0; i < st.NumMethod(); i++ {
		mtd := st.Method(i).Type()
		mtd2 := st.Type().Method(i)
		fmt.Println(mtd.String(), mtd2.Name)
	}

	// default settings for some fields
	if len(settings.PageRouteName) == 0 {
		settings.PageRouteName = settings.Name
	}
	if len(settings.AjaxRouteName) == 0 {
		settings.AjaxRouteName = settings.PageRouteName
	}
	if len(settings.Title) == 0 {
		settings.Title = settings.Name
	}
	if len(settings.CollectionName) == 0 {
		settings.CollectionName = settings.Name
	}
	if settings.Ajax == nil {
		settings.Ajax = Func{}
	}

	module := panel.RouterGroup.Group(settings.PageRouteName)
	module.Use(func(c *gin.Context) {
		c.Set("moduleName", settings.PageRouteName)
	})
	module.GET("/", s.Page)
	panel.RouterGroup.POST("/ajax/"+settings.AjaxRouteName+"/:method", s.Ajax)
	s.Init(settings, panel)

	return s
}

func (m *Module) Init(settings *ModuleSettings, panel *Panel) {
	m.Settings = settings
	m.Panel = *panel
	if m.Collection == nil {
		m.Collection = mongo.DB(panel.DBName).C(settings.Name)
	}
}

// Ajax - chooise method for module "admins"
func (m *Module) Ajax(c *gin.Context) {
	methodFound := false
	for ajaxRoute, ajaxFunc := range m.Settings.Ajax {
		if c.Param("method") == ajaxRoute {
			ajaxFunc(c)
			methodFound = true
			break
		}
	}

	if !methodFound {
		c.String(400, "Method not found in module \"AdminsModule\"!")
	}
}
func (m *Module) Page(c *gin.Context) {
	c.HTML(200, m.Settings.Name+".html", gin.H{
		"title": m.Settings.Title,
		"user":  c.MustGet("user"),
	})
}

// Create new panel
func Create(s Settings) *Panel {
	panel := Panel{
		Settings: Settings{
			Port:        8080,
			AuthSalt:    "+Af761",
			AuthPrefix:  "adm-summer-",
			Title:       "Summer Panel",
			Path:        "/admin",
			Views:       "./templates/main",
			ViewsDoT:    "./templates/doT.js",
			Files:       "./files",
			TMPs:        "/tmp",
			Language:    "EN",
			DBName:      "summerPanel",
			DefaultPage: "/settings",
			Vars:        map[string]interface{}{},
			FirstStart:  func() {},
			Engine:      gin.New(),
		},
	}
	// apply default settings
	Extend(&panel.Settings, &s)
	if panel.Vars == nil {
		panel.Vars = make(map[string]interface{})
	}
	panel.Vars["panelPath"] = panel.Path
	panel.Vars["title"] = panel.Title

	// init autoincrement module
	ai.Connect(mongo.DB(panel.DBName).C("ai"))

	funcMap := template.FuncMap{"dot": dot, "jsoner": jsoner, "var": func(key string) interface{} {
		return panel.Vars[key]
	}}
	ttpl.Use(panel.Engine, []string{PackagePath() + "/templates/main/*", panel.Views + "/*"}, panel.ViewsDoT, funcMap)

	// включение статических файлов
	panel.Engine.Static(panel.Path+"/files", panel.Files)
	panel.Engine.Static(panel.Path+"/pkgFiles", PackagePath()+"/files")

	// запуск веб-сервера
	go func() {
		panic(panel.Engine.Run(":" + types.String(panel.Port)))
	}()

	admins.Init(&panel)
	panel.RouterGroup = panel.Engine.Group(panel.Path)
	admins.Auth(panel.RouterGroup)
	panel.RouterGroup.GET("/", func(c *gin.Context) {
		c.Header("Expires", time.Now().String())
		c.Header("Cache-Control", "no-cache")
		c.Redirect(301, panel.DefaultPage)
	})
	return &panel
}

func Wait() {
	for {
		time.Sleep(time.Second)
	}
}
