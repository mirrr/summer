package summer

import (
	"text/template"
	"time"
	"ttpl"

	"github.com/night-codes/mgo-ai"
	"github.com/night-codes/mgo-wrapper"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/night-codes/types.v1"
)

type (

	// Settings intended for data transmission into the Create method of package
	Settings struct {
		Port              uint
		Title             string
		AuthSalt          string
		AuthPrefix        string
		Path              string // URL path of panel - "/" by default
		Views             string // file path of ./templates directory
		ViewsDoT          string // file path of doT.js templates directory
		Files             string // file path of ./files directory
		TMPs              string // file path of /tmp directory
		DBName            string // MongoDB database name
		DefaultPage       string
		Language          string
		UsersCollection   string
		Debug             bool
		Vars              map[string]interface{}
		TFuncMap          template.FuncMap
		FirstStart        func()
		RouterGroup       *gin.RouterGroup
		Engine            *gin.Engine
		DisableAuth       bool     // if TRUE - without summer auth
		DisableFirstStart bool     // if TRUE - without first user creating
		JS                []string // external JS resources
		CSS               []string // external CSS resources
	}

	//Panel struct
	Panel struct {
		Settings
		// RootMenu is zerro-level menu
		RootMenu *Menu
		// MainMenu is main admin-panel menu
		MainMenu *Menu
		// DropMenu is top user dropdown menu
		DropMenu *Menu
		// Groups
		Groups *GroupsList
	}
)

// Create new panel
func Create(s Settings) *Panel {
	var engine *gin.Engine
	if s.Debug {
		gin.SetMode(gin.DebugMode)
		engine = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		engine = gin.New()
	}
	rootMenu := &Menu{Title: "[Root]"}
	panel := Panel{
		Settings: Settings{
			Port:            8080,
			AuthSalt:        "+Af761",
			AuthPrefix:      "adm-summer-",
			Title:           "Summer Panel",
			Path:            "",
			Views:           "templates/main",
			ViewsDoT:        "templates/doT.js",
			Files:           "files",
			TMPs:            "/tmp",
			Language:        "EN",
			DBName:          "summerPanel",
			DefaultPage:     "/settings",
			UsersCollection: "admins",
			Vars:            map[string]interface{}{},
			FirstStart:      func() {},
			Engine:          engine,
		},
		RootMenu: rootMenu,
		MainMenu: rootMenu.Add("[Main]"),
		DropMenu: rootMenu.Add("[Drop]"),
		Groups:   &GroupsList{},
	}
	// apply default settings
	extend(&panel.Settings, &s)
	if panel.Vars == nil {
		panel.Vars = make(map[string]interface{})
	}
	// for correct templates render
	if len(panel.Path) > 0 {
		if panel.Path[len(panel.Path)-1] == '/' {
			panel.Path = panel.Path[:len(panel.Path)-1]
		}
		if panel.Path[0] != '/' {
			panel.Path = "/" + panel.Path
		}
	}
	if len(panel.DefaultPage) == 0 || panel.DefaultPage[0] != '/' {
		panel.DefaultPage = "/" + panel.DefaultPage
	}
	panel.Vars["panel"] = &panel
	panel.Vars["path"] = panel.Path
	panel.Vars["title"] = panel.Title
	panel.Vars["mainMenu"] = panel.MainMenu
	panel.Vars["dropMenu"] = panel.DropMenu

	// init autoincrement module
	ai.Connect(mongo.DB(panel.DBName).C("ai"))

	funcMap := template.FuncMap{
		"jsoner": jsoner,
		"menu":   getMenuItems,
		"user":   getByLogin,
		"tabs":   getTabs,
		"site":   getSite,
		"var": func(key string) interface{} {
			return panel.Vars[key]
		},
	}
	ttpl.Use(panel.Engine, []string{PackagePath() + "/templates/main/", panel.Views + "/"}, panel.ViewsDoT, funcMap)

	// включение статических файлов
	panel.Engine.Static(panel.Path+"/files", panel.Files)
	panel.Engine.Static(panel.Path+"/pkgFiles", PackagePath()+"/files")

	// запуск веб-сервера
	go func() {
		panic(panel.Engine.Run(":" + types.String(panel.Port)))
	}()

	auth.Init(&panel)
	Users.init(&panel)
	panel.RouterGroup = panel.Engine.Group(panel.Path)
	auth.Auth(panel.RouterGroup)
	panel.RouterGroup.GET("/", func(c *gin.Context) {
		c.Header("Expires", time.Now().String())
		c.Header("Cache-Control", "no-cache")
		c.Redirect(301, panel.Path+panel.DefaultPage)
	})
	return &panel
}

// AddModule provide adding new panel module
func (panel *Panel) AddModule(settings *ModuleSettings, s Simple) Simple {
	return createModule(panel, settings, s)
}

func Wait() {
	for {
		time.Sleep(time.Second)
	}
}
