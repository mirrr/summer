package summer

import (
	"reflect"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/night-codes/mgo-wrapper"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
)

type (
	// Func is alias for map[string]func(c *gin.Context)
	Func map[string]func(c *gin.Context)
	// WebFunc is alias for map[string]func(c *websocket.Conn)
	WebFunc map[string]func(c *gin.Context, ws *websocket.Conn)

	//Module struct
	Module struct {
		Panel
		Collection *mgo.Collection
		Settings   *ModuleSettings
	}

	//ModuleSettings struct
	ModuleSettings struct {
		Name             string
		Menu             *Menu
		MenuTitle        string
		MenuOrder        int
		PageRouteName    string
		AjaxRouteName    string
		SocketsRouteName string
		Title            string
		CollectionName   string
		TemplateName     string
		AllowGroups      []string
		AllowRoles       []string
		AllowUsers       []uint64
		Ajax             Func
		Websockets       WebFunc
		Icon             string
		GroupTo          Simple
		GroupTitle       string
	}

	// Simple module interface
	Simple interface {
		Init(settings *ModuleSettings, panel *Panel)
		Page(c *gin.Context)
		Ajax(c *gin.Context)
		Websockets(c *gin.Context)
		GetSettings() *ModuleSettings
	}
)

var (
	modulesList   = map[string]Simple{}
	modulesListMu = sync.Mutex{}
	wsupgrader    = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

// Ajax  is default module's ajax method
func (m *Module) Ajax(c *gin.Context) {
	method := strings.ToLower(c.Param("method"))
	if len(method) > 0 && method[0] == '/' {
		method = method[1:]
	}
	for ajaxRoute, ajaxFunc := range m.Settings.Ajax {
		if method == ajaxRoute {
			ajaxFunc(c)
			return
		}
	}
	c.String(400, `Method not found in module "`+m.Settings.Name+`"!`)
}

// Websockets  is default module's websockets method
func (m *Module) Websockets(c *gin.Context) {
	for websocketsRoute, websocketsFunc := range m.Settings.Websockets {
		if strings.ToLower(c.Param("method")) == websocketsRoute {
			if conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil); err == nil {
				websocketsFunc(c, conn)
				return
			}
			break
		}
	}
	c.String(400, `Method not found in module "`+m.Settings.Name+`"!`)
}

// Page is default module's page rendering method
func (m *Module) Page(c *gin.Context) {
	c.HTML(200, m.Settings.TemplateName+".html", obj{"title": m.Settings.Title, "action": c.Param("action")})
}

// Init is default module's initial method
func (m *Module) Init(settings *ModuleSettings, panel *Panel) {
	m.Settings = settings
	m.Panel = *panel
	if m.Collection == nil {
		m.Collection = mongo.DB(panel.DBName).C(settings.CollectionName)
	}
}

// GetSettings needs for correct settings getting from module struct
func (m *Module) GetSettings() *ModuleSettings {
	return m.Settings
}

// Create new module
func createModule(panel *Panel, settings *ModuleSettings, s Simple) Simple {

	modulesListMu.Lock()
	if settings.Name == "" || modulesList[settings.Name] != nil {
		panic(`Repeated use of module name "` + settings.Name + `"`)
	}
	modulesListMu.Unlock()
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
			switch name {
			case "ajax", "page", "websockets":
				continue
			}
			if method == "func(*gin.Context)" {
				method := st.Method(i).Interface().(func(*gin.Context))
				settings.Ajax[name] = method
			} else if method == "func(*gin.Context, *websocket.Conn)" {
				method := st.Method(i).Interface().(func(*gin.Context, *websocket.Conn))
				settings.Websockets[name] = method
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
	if len(settings.GroupTitle) == 0 {
		settings.GroupTitle = settings.MenuTitle
	}
	if len(settings.CollectionName) == 0 {
		settings.CollectionName = strings.Replace(settings.Name, "/", "-", -1)
	}
	if len(settings.TemplateName) == 0 {
		settings.TemplateName = strings.Replace(settings.Name, "/", "-", -1)
	}

	moduleGroup := panel.RouterGroup.Group(settings.PageRouteName)
	moduleGroup.Use(func(c *gin.Context) {
		c.Header("Module", settings.PageRouteName)
		c.Header("Login", c.MustGet("login").(string))
		c.Header("Title", settings.Title)
		c.Header("Path", panel.Path)
		c.Header("Action", c.Param("action"))
		header := c.Writer.Header()
		header["Css"] = panel.CSS
		header["Js"] = panel.JS
	})

	moduleGroup.GET("/*action", s.Page)
	panel.RouterGroup.POST("/ajax/"+settings.AjaxRouteName+"/*method", s.Ajax)
	panel.RouterGroup.GET("/websocket/"+settings.SocketsRouteName+"/*method", s.Websockets)
	s.Init(settings, panel)

	modulesListMu.Lock()
	modulesList[settings.Name] = s
	modulesListMu.Unlock()
	return s
}
