package summer

import (
	"text/template"
	"time"
	"ttpl"

	"github.com/gin-gonic/gin"
	ai "github.com/night-codes/mgo-ai"
	mongo "github.com/night-codes/mgo-wrapper"
	"gopkg.in/night-codes/types.v1"
)

type (

	// Settings intended for data transmission into the Create method of package
	Settings struct {
		Port              uint
		Title             string
		DefaultPage       string
		AuthSalt          string
		AuthPrefix        string                 // prefix for cookies names
		AuthSkipIP        bool                   // disable using IP as salt for cookie hash generation
		Path              string                 // URL path of panel - "/" by default
		Views             string                 // file path of ./templates directory
		ViewsDoT          string                 // file path of doT.js templates directory
		Files             string                 // file path of ./files directory
		DBName            string                 // MongoDB database name
		UsersCollection   string                 // collection for panel's users
		NotifyCollection  string                 // collection for panel's notifications
		AICollection      string                 // collection for AUTO_INCREMENT
		Debug             bool                   // show `gin` debugging messages
		Vars              map[string]interface{} // variables, that can be used in templates {{var "variableName"}}
		TFuncMap          template.FuncMap       // `gin` template functions
		MakeEngineFn      func(*gin.Engine)      // Engine make callback
		MakeRouterGroupFn func(*gin.RouterGroup) // RouterGroup make callback
		FirstStart        func()                 // function called after first user creation
		DisableAuth       bool                   // if TRUE - without summer auth
		DisableFirstStart bool                   // if TRUE - without first user creating (FirstStart function don't called)
		JS                []string               // external JS resources
		CSS               []string               // external CSS resources
		RouterGroup       *gin.RouterGroup
		Engine            *gin.Engine
		HashDBFn          func(login, password, authSalt string) string
		HashCookieFn      func(login, password, authSalt, ip, userAgent string) string
		DelayStart        bool
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
		// Modules
		Modules *ModuleList
		// Users
		Users *Users

		auth     *auth
		menuList *menuList
		AI       ai.AI
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
	panel := Panel{
		Settings: Settings{
			Port:              8080,
			AuthSalt:          "+Af761",
			AuthPrefix:        "adm-summer-",
			Title:             "Summer Panel",
			Path:              "",
			Views:             "templates/main",
			ViewsDoT:          "templates/doT.js",
			Files:             "files",
			DBName:            "summerPanel",
			DefaultPage:       "/settings",
			UsersCollection:   "admins",
			AICollection:      "ai",
			Vars:              map[string]interface{}{},
			FirstStart:        func() {},
			MakeEngineFn:      func(*gin.Engine) {},
			MakeRouterGroupFn: func(*gin.RouterGroup) {},
			HashDBFn:          func(login, password, authSalt string) string { return H3hash(password + authSalt) },
			HashCookieFn:      func(login, dbpass, authSalt, ip, userAgent string) string { return H3hash(ip + dbpass + authSalt) },
			Engine:            engine,
		},
		menuList: createMenuList(),
		Modules:  createModuleList(),
		RootMenu: &Menu{Title: "[Root]"},
		MainMenu: &Menu{Title: "[Main]"},
		DropMenu: &Menu{Title: "[Drop]"},
		Groups:   new(GroupsList),
		Users:    new(Users),
		auth:     new(auth),
	}
	// apply default settings
	extend(&panel.Settings, &s)
	panel.MakeEngineFn(panel.Engine)
	panel.init()

	return &panel
}

// AddModule provide adding new module to Panel
func (panel *Panel) AddModule(settings *ModuleSettings, s Simple) Simple {
	return createModule(panel, settings, s)
}

// AddPage provide adding simple page module to Panel
func (panel *Panel) AddPage(name string, title string, menu *Menu, originTemplate ...bool) Simple {
	settings := &ModuleSettings{
		Name:  name,
		Title: title,
		Menu:  menu,
	}
	if len(originTemplate) > 0 {
		settings.OriginTemplate = originTemplate[0]
	}
	return panel.AddModule(settings, &Module{})
}

// AddOpenPage provide adding simple page without authorization
func (panel *Panel) AddOpenPage(name string, title string, menu *Menu, originTemplate ...bool) Simple {
	settings := &ModuleSettings{
		Name:        name,
		Title:       title,
		Menu:        menu,
		DisableAuth: true,
	}
	if len(originTemplate) > 0 {
		settings.OriginTemplate = originTemplate[0]
	}
	return panel.AddModule(settings, &Module{})
}

// initial method for *Panel
func (panel *Panel) init() {
	panel.RootMenu.init(panel, nil)
	panel.MainMenu.init(panel, panel.RootMenu)
	panel.DropMenu.init(panel, panel.RootMenu)

	panel.Users.init(panel)
	panel.auth.init(panel)

	panel.correctPath()
	panel.setVariables()
	panel.initTpl()

	// init autoincrement module
	panel.AI = *ai.Create(mongo.DB(panel.DBName).C(panel.AICollection))

	// static files
	panel.Engine.Use(gzipper)
	panel.Engine.Static(panel.Path+"/files", panel.Files)
	panel.Engine.Static(panel.Path+"/pkgFiles", PackagePath()+"/files")

	// main rout group
	panel.RouterGroup = panel.Engine.Group(panel.Path)
	panel.MakeRouterGroupFn(panel.RouterGroup)
	panel.RouterGroup.GET("/", func(c *gin.Context) {
		c.Header("Expires", time.Now().String())
		c.Header("Cache-Control", "no-cache")
		c.Redirect(301, panel.Path+panel.DefaultPage)
	})

	// starting web-server
	if !panel.DelayStart {
		go func() {
			panic(panel.Engine.Run(":" + types.String(panel.Port)))
		}()
	}
}

func (panel *Panel) initTpl() {
	if panel.TFuncMap == nil {
		panel.TFuncMap = template.FuncMap{}
	}

	panel.TFuncMap["jsoner"] = jsoner
	panel.TFuncMap["menu"] = getMenuItems
	panel.TFuncMap["user"] = func(login string) (user *UsersStruct) {
		user, _ = panel.Users.GetByLogin(login)
		return
	}
	panel.TFuncMap["tabs"] = getTabs
	panel.TFuncMap["site"] = getSite
	panel.TFuncMap["var"] = func(key string) interface{} {
		return panel.Vars[key]
	}
	panel.TFuncMap["o"] = func(key string) interface{} {
		return ""
	}

	ttpl.Use(panel.Engine, []string{PackagePath() + "/templates/main/", panel.Views + "/"}, panel.ViewsDoT, panel.TFuncMap)
}

func (panel *Panel) setVariables() {
	if panel.Vars == nil {
		panel.Vars = make(map[string]interface{})
	}
	panel.Vars["panel"] = &panel
	panel.Vars["path"] = panel.Path
	panel.Vars["title"] = panel.Title
	panel.Vars["mainMenu"] = panel.MainMenu
	panel.Vars["dropMenu"] = panel.DropMenu
}

func (panel *Panel) correctPath() {
	if len(panel.Path) > 0 {
		if panel.Path[len(panel.Path)-1] == '/' {
			panel.Path = panel.Path[:len(panel.Path)-1]
		}
		if len(panel.Path) > 0 && panel.Path[0] != '/' {
			panel.Path = "/" + panel.Path
		}
	}
	if len(panel.DefaultPage) == 0 || panel.DefaultPage[0] != '/' {
		panel.DefaultPage = "/" + panel.DefaultPage
	}
}

// Run delayed application
func (panel *Panel) Run() {
	go func() {
		panic(panel.Engine.Run(":" + types.String(panel.Port)))
	}()
}

// Wait locks go-routine
func Wait() {
	select {}
}
