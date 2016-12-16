package summer

import (
	"reflect"
	"strings"
	"text/template"
	"time"
	"ttpl"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/mirrr/mgo-ai"
	"github.com/mirrr/mgo-wrapper"
	"gopkg.in/mirrr/types.v1"
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
		UserCollection    string
		Debug             bool
		Vars              map[string]interface{}
		TFuncMap          template.FuncMap
		FirstStart        func()
		RouterGroup       *gin.RouterGroup
		Engine            *gin.Engine
		DisableAuth       bool // if TRUE - without summer auth
		DisableFirstStart bool // if TRUE - without first user creating
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
			Engine:         engine,
		},
		RootMenu: rootMenu,
		MainMenu: rootMenu.Add("[Main]"),
		DropMenu: rootMenu.Add("[Drop]"),
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
	panel.Vars["panelPath"] = panel.Path
	panel.Vars["title"] = panel.Title
	panel.Vars["mainMenu"] = panel.MainMenu
	panel.Vars["dropMenu"] = panel.DropMenu

	// init autoincrement module
	ai.Connect(mongo.DB(panel.DBName).C("ai"))

	funcMap := template.FuncMap{"jsoner": jsoner, "menu": getMenuItems, "var": func(key string) interface{} {
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
		c.Redirect(301, panel.Path+panel.DefaultPage)
	})
	return &panel
}

// AddModule provide adding new panel module
func (panel *Panel) AddModule(settings *ModuleSettings, s Simple) Simple {
	if settings.Ajax == nil {
		settings.Ajax = Func{}
	}
	if settings.Websockets == nil {
		settings.Websockets = WebFunc{}
	}
	st := reflect.ValueOf(s)
	for i := 0; i < st.NumMethod(); i++ {
		method := st.Method(i).Type().String()
		if len(method) > 17 && method[:17] == "func(*gin.Context" {
			name := strings.ToLower(st.Type().Method(i).Name)
			if name != "ajax" && name != "page" && name != "websockets" {
				if method == "func(*gin.Context)" {
					method := st.Method(i).Interface().(func(*gin.Context))
					settings.Ajax[name] = method
				} else if method == "func(*gin.Context, *websocket.Conn)" {
					method := st.Method(i).Interface().(func(*gin.Context, *websocket.Conn))
					settings.Websockets[name] = method
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
	if len(settings.SocketsRouteName) == 0 {
		settings.SocketsRouteName = settings.PageRouteName
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
	panel.RouterGroup.GET("/websocket/"+settings.SocketsRouteName+"/:method", s.Websockets)
	s.Init(settings, panel)

	modulesListMu.Lock()
	modulesList[settings.Name] = s
	modulesListMu.Unlock()
	return s
}

func Wait() {
	for {
		time.Sleep(time.Second)
	}
}
