package summer

import (
	"reflect"
	"strings"
	"text/template"
	"time"
	"ttpl"

	"github.com/gin-gonic/gin"
	"github.com/mirrr/mgo-ai"
	"github.com/mirrr/mgo-wrapper"
	"gopkg.in/mirrr/types.v1"
)

type (

	// Settings intended for data transmission into the Create method of package
	Settings struct {
		Port           uint
		Title          string
		AuthSalt       string
		AuthPrefix     string
		Path           string // URL path of panel - "/" by default
		Views          string // file path of ./templates directory
		ViewsDoT       string // file path of doT.js templates directory
		Files          string // file path of ./files directory
		TMPs           string // file path of /tmp directory
		DBName         string // MongoDB database name
		DefaultPage    string
		Language       string
		UserCollection string
		Vars           map[string]interface{}
		TFuncMap       template.FuncMap
		FirstStart     func()
		RouterGroup    *gin.RouterGroup
		Engine         *gin.Engine
	}

	//Panel struct
	Panel struct {
		Settings
	}
)

// Create new panel
func Create(s Settings) *Panel {
	panel := Panel{
		Settings: Settings{
			Port:           8080,
			AuthSalt:       "+Af761",
			AuthPrefix:     "adm-summer-",
			Title:          "Summer Panel",
			Path:           "",
			Views:          "templates/main",
			ViewsDoT:       "templates/doT.js",
			Files:          "files",
			TMPs:           "/tmp",
			Language:       "EN",
			DBName:         "summerPanel",
			DefaultPage:    "/settings",
			UserCollection: "admins",
			Vars:           map[string]interface{}{},
			FirstStart:     func() {},
			Engine:         gin.New(),
		},
	}
	// apply default settings
	extend(&panel.Settings, &s)
	if panel.Vars == nil {
		panel.Vars = make(map[string]interface{})
	}
	panel.Vars["panelPath"] = panel.Path
	panel.Vars["title"] = panel.Title
	panel.Vars["modules"] = &modulesList
	panel.Vars["menus"] = &menusList

	// init autoincrement module
	ai.Connect(mongo.DB(panel.DBName).C("ai"))

	funcMap := template.FuncMap{"jsoner": jsoner, "var": func(key string) interface{} {
		return panel.Vars[key]
	}}
	ttpl.Use(panel.Engine, []string{PackagePath() + "/templates/main/", panel.Views + "/"}, panel.ViewsDoT, funcMap)

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

// AddModule provide adding new panel module
func (panel *Panel) AddModule(settings *ModuleSettings, s Simple) Simple {
	if settings.Ajax == nil {
		settings.Ajax = Func{}
		st := reflect.ValueOf(s)
		for i := 0; i < st.NumMethod(); i++ {
			if st.Method(i).Type().String() == "func(*gin.Context)" {
				name := strings.ToLower(st.Type().Method(i).Name)
				if name != "ajax" && name != "page" {
					method := st.Method(i).Interface().(func(*gin.Context))
					settings.Ajax[name] = method
				}
			}
		}
	}

	// default settings for some fields
	if len(settings.PageRouteName) == 0 {
		settings.PageRouteName = settings.Name
	}
	if len(settings.AjaxRouteName) == 0 {
		settings.AjaxRouteName = settings.PageRouteName
	}
	if len(settings.Title) == 0 {
		settings.Title = strings.Replace(settings.Name, "/", " ", -1)
	}
	if len(settings.MenuTitle) == 0 {
		settings.MenuTitle = settings.Title
	}
	if len(settings.CollectionName) == 0 {
		settings.CollectionName = strings.Replace(settings.Name, "/", "-", -1)
	}
	if len(settings.TemplateName) == 0 {
		settings.TemplateName = strings.Replace(settings.Name, "/", "-", -1)
	}

	moduleGroup := panel.RouterGroup.Group(settings.PageRouteName)
	moduleGroup.Use(func(c *gin.Context) {
		c.Set("moduleName", settings.PageRouteName)
	})
	moduleGroup.GET("/", s.Page)
	panel.RouterGroup.POST("/ajax/"+settings.AjaxRouteName+"/:method", s.Ajax)
	s.Init(settings, panel)

	modulesList[settings.Name] = s
	return s
}

func Wait() {
	for {
		time.Sleep(time.Second)
	}
}
