package summer

import (
	"text/template"
	"time"
	"ttpl"

	"github.com/gin-gonic/gin"
	"github.com/mirrr/mgo-ai"
	"github.com/mirrr/mgo-wrapper"
	"gopkg.in/mirrr/types.v1"
)

type (
	// Settings intended for data transmission into the Init method of package
	Settings struct {
		Port        uint
		Title       string
		AuthSalt    string
		AuthPrefix  string
		Path        string // URL path of panel - "/" by default
		Views       string // file path of ./templates directory
		Files       string // file path of ./files directory
		TMPs        string // file path of /tmp directory
		DBName      string // MongoDB database name
		DefaultPage string
		Language    string
		Vars        map[string]interface{}
		TFuncMap    template.FuncMap
		FirstStart  func()
		r           *gin.RouterGroup
	}

	//Panel ...
	Panel struct {
		Settings
	}

	Simple interface {
		Init(p *Panel)
		Page(c *gin.Context)
		Ajax(c *gin.Context)
	}
)

func Init(s Settings) *Panel {
	panel := Panel{
		Settings: Settings{
			Port:        8080,
			AuthSalt:    "+Af761",
			AuthPrefix:  "adm-summer-",
			Title:       "Summer Panel",
			Path:        "/admin",
			Views:       "./templates/main",
			Files:       "./files",
			TMPs:        "/tmp",
			Language:    "EN",
			DBName:      "summerPanel",
			DefaultPage: "/settings",
			Vars:        map[string]interface{}{},
			FirstStart:  func() {},
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
	ai.Connect(mongo.DB(s.DBName).C("ai"))

	r := gin.New()
	funcMap := template.FuncMap{"dot": dot, "jsoner": jsoner, "var": func(key string) interface{} {
		return panel.Vars[key]
	}}
	ttpl.Use(r, []string{PackagePath() + "/templates/main/*", s.Views + "/*"}, funcMap)

	// включение статических файлов
	r.Static(panel.Path+"/files", s.Files)
	r.Static(panel.Path+"/pkgFiles", PackagePath()+"/files")

	// запуск веб-сервера
	go func() {
		panic(r.Run(":" + types.String(s.Port)))
	}()

	admins.Init(&panel)
	panel.r = r.Group(panel.Path)
	admins.Auth(panel.r)
	panel.r.GET("/", func(c *gin.Context) {
		c.Header("Expires", time.Now().String())
		c.Header("Cache-Control", "no-cache")
		c.Redirect(301, s.DefaultPage)
	})
	return &panel
}

func Wait() {
	for {
		time.Sleep(time.Second)
	}
}

func (panel *Panel) AddModule(module string, s Simple) Simple {
	panel.r.GET("/"+module, s.Page)
	panel.r.POST("/ajax/"+module+"/:method", s.Ajax)
	s.Init(panel)
	return s
}
