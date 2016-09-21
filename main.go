package summer

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mirrr/go-gin-ttpl"
	"github.com/mirrr/mgo-ai"
	"github.com/mirrr/mgo-wrapper"
	"text/template"
	/*	"gopkg.in/mgo.v2"
		"gopkg.in/mgo.v2/bson"
		"gopkg.in/mirrr/types.v1"
		"time"*/)

var (
	settings Settings
)

type (
	Settings struct {
		Title  string
		Path   string // URL path of panel - "/" by default
		Views  string // file path of ./templates directory
		Files  string // file path of ./files directory
		TMPs   string // file path of /tmp directory
		DBName string // MongoDB database name
		Glob   map[string]interface{}
	}
)

func Init(s Settings) {

	if len(s.Title) == 0 {
		s.Title = "Summer Panel"
	}
	if len(s.Path) == 0 {
		s.Path = "/admin"
	}
	if len(s.Views) == 0 {
		s.Views = "./views"
	}
	if len(s.Files) == 0 {
		s.Files = "./files"
	}
	if len(s.TMPs) == 0 {
		s.TMPs = "/tmp"
	}
	if len(s.TMPs) == 0 {
		s.TMPs = "/tmp"
	}
	if len(s.DBName) == 0 {
		s.DBName = "summerPanel"
	}
	if s.Glob == nil {
		s.Glob = make(map[string]interface{})
	}
	s.Glob["panelPath"] = s.Path
	s.Glob["title"] = s.Title

	settings = s
	fmt.Println(settings)

	ai.Connect(mongo.DB(settings.DBName).C("ai"))

	r := gin.New()
	// плагин к шаблонизатору возвращающий глобальный список обьектов
	glob := func(key string) interface{} {
		return settings.Glob[key]
	}
	ttpl.Use(r, settings.Views+"/*", template.FuncMap{"dot": dot, "jsoner": jsoner, "glob": glob})
	r.Static(settings.Path+"/files", settings.Files)
}
