package summer

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mirrr/mgo-ai"
	"github.com/mirrr/mgo-wrapper"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mirrr/types.v1"
	"net/http"
	"time"
)

type (
	AdminsStruct struct {
		ID       uint64 `form:"id"  json:"id"  bson:"_id"`
		Root     bool   `form:"-"  json:"root"  bson:"root"`
		Name     string `form:"name" json:"name" bson:"name" binding:"required,min=3"`
		Notice   string `form:"notice" json:"notice" bson:"notice"`
		Login    string `form:"login" json:"login" bson:"login" binding:"required`
		Password string `form:"password" json:"-" bson:"password"`
		Updated  uint   `form:"-" json:"-" bson:"updated"`
		Deleted  bool   `form:"-" json:"-" bson:"deleted"`
	}

	Admins struct {
		collection *mgo.Collection
		Panel
	}
)

var (
	admins = Admins{}
)

func (a *Admins) Init(panel *Panel) {
	a.Panel = *panel
	a.collection = mongo.DB(panel.DBName).C("admins")
}

func (a *Admins) Page(c *gin.Context) {
	c.HTML(200, "admins.html", gin.H{
		"title":  "Администраторы",
		"user":   c.MustGet("user").(AdminsStruct),
		"active": obj{"admins": true},
	})
}
func (a *Admins) Auth(g *gin.RouterGroup) {
	g.Use(a.Login(g.BasePath()))
	g.POST("/z-auth", dummy) // хак для авторизации
	g.POST("/z-register", dummy)
	g.GET("/logout", a.Logout(g.BasePath()))
}

func (a *Admins) Logout(panelPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		http.SetCookie(c.Writer, &http.Cookie{
			Name:    a.AuthPrefix + "hash",
			Value:   "",
			Path:    "/",
			MaxAge:  1,
			Expires: time.Now(),
		})
		c.Header("Expires", time.Now().String())
		c.Header("Cache-Control", "no-cache")
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, "<meta http-equiv='refresh' content='0; url="+panelPath+"' />")
	}
}

func (a *Admins) Login(panelPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminsArr := map[string]AdminsStruct{}
		for _, v := range admins.GetArr() {
			adminsArr[v.Login] = v
		}

		// 	регистрация первого пользователя админки
		if len(adminsArr) == 0 {
			defer c.Abort()
			login, e1 := c.GetPostForm("admin-z-login")
			password, e2 := c.GetPostForm("admin-z-password")
			password2, e3 := c.GetPostForm("admin-z-password-2")
			if e1 && e2 && e3 {
				if password == password2 {
					if len(login) > 3 && len(password) > 6 {
						if err := admins.AddRaw(AdminsStruct{
							Login:    login,
							Password: password,
							Name:     "Admin",
							Root:     true,
						}); err != nil {
							c.String(400, "Ошибка БД")
							return
						}
						newStart()
						c.String(200, "Ok")
					} else {
						c.String(400, "Логин или пароль слишком коротки!")
					}
				} else {
					c.String(400, "Пароли не совпадают!")
				}
				return
			}

			// add admin user
			c.HTML(200, "firstStart.html", gin.H{"panelPath": panelPath})
			c.Abort()
			return
		}

		// авторизация пользователя админки
		login, e1 := c.GetPostForm("admin-z-login")
		password, e2 := c.GetPostForm("admin-z-password")
		if e1 && e2 {
			if user, exists := adminsArr[login]; exists && user.Password == H3hash(password+a.AuthSalt) {
				setCookie(c, a.AuthPrefix+"login", login)
				setCookie(c, a.AuthPrefix+"hash", H3hash(c.ClientIP()+user.Password+a.AuthSalt))
				c.String(200, "Ok")
			} else {
				c.String(400, "Wrong password!")
			}
			c.Abort()
			return
		} else {
			login, e1 := c.Cookie(a.AuthPrefix + "login")
			hash, e2 := c.Cookie(a.AuthPrefix + "hash")
			if e1 == nil && e2 == nil {
				if user, exists := adminsArr[login]; exists && hash == H3hash(c.ClientIP()+user.Password+a.AuthSalt) {
					c.Set("user", user)
					c.Next()
					return
				}
			}
		}
		c.HTML(200, "login.html", gin.H{"panelPath": panelPath})
		c.Abort()
	}
}

// Ajax - chooise method for module "admins"
func (a *Admins) Ajax(c *gin.Context) {
	switch c.Param("method") {
	case "add":
		a.Add(c)
	case "getAll":
		a.GetAll(c)
	case "get":
		a.Get(c)
	case "edit":
		a.Edit(c)
	case "remove":
		a.Remove(c)
	default:
		c.String(400, "Method not found in module \"Admins\"!")
	}
}

// Add new admin
func (a *Admins) Add(c *gin.Context) {
	var admin AdminsStruct
	if err := c.Bind(&admin); err != nil {
		c.String(400, "Поле \"Имя\" слишком коротко")
		return
	}
	admin.ID = ai.Next("admins")
	admin.Password = H3hash(admin.Password + a.AuthSalt)
	admin.Updated = uint(time.Now().Unix() / 60)
	if err := a.collection.Insert(admin); err != nil {
		ai.Cancel("admins")
		fmt.Println("Error (userAdd):", err)
		c.String(400, "Неизвестная ошибка")
		return
	}
	c.JSON(200, obj{"data": admin})
}

// AddRaw adds new admin from struct
func (a *Admins) AddRaw(admin AdminsStruct) error {
	admin.ID = ai.Next("admins")
	admin.Password = H3hash(admin.Password + a.AuthSalt)
	admin.Updated = uint(time.Now().Unix() / 60)
	if err := admins.collection.Insert(admin); err != nil {
		ai.Cancel("admins")
		return err
	}
	return nil
}

// Edit element
func (a *Admins) Edit(c *gin.Context) {
	var admin AdminsStruct
	if err := c.Bind(&admin); err != nil {
		c.String(400, "Поле \"Имя\" слишком коротко")
		return
	}

	set := obj{
		"name":    admin.Name,
		"notice":  admin.Notice,
		"email":   admin.Login,
		"updated": uint(time.Now().Unix() / 60),
	}
	if len(admin.Password) > 0 {
		set["password"] = H3hash(admin.Password + a.AuthSalt)
	}
	if err := a.collection.UpdateId(admin.ID, obj{
		"$set": set,
	}); err != nil {
		fmt.Println("Error (userEdit):", err)
		c.String(400, "Неизвестная ошибка")
		return
	}
	c.JSON(200, obj{"data": admin})
}

// Remove element
func (a *Admins) Remove(c *gin.Context) {
	id := types.Uint64(c.PostForm("id"))
	if err := a.collection.UpdateId(id, obj{"$set": obj{"deleted": true}}); err != nil {
		fmt.Println("Error (userRemove):", err)
		c.String(400, "Неизвестная ошибка")
		return
	}
	c.JSON(200, obj{"data": obj{"id": id}})
}

// Get element
func (a *Admins) Get(c *gin.Context) {
	id := types.Uint64(c.PostForm("id"))
	admin := AdminsStruct{}
	if err := a.collection.FindId(id).One(&admin); err != nil {
		fmt.Println("Error (userGet):", err)
	}
	c.JSON(200, obj{"data": admin})
}

// GetAll return all elements
func (a *Admins) GetAll(c *gin.Context) {
	params := struct {
		Search string `form:"search"  json:"search"`
	}{}
	search := obj{"deleted": obj{"$ne": true}}
	if err := c.Bind(&params); err == nil {
		if len(params.Search) > 0 {
			regex := bson.RegEx{Pattern: params.Search, Options: "i"}
			search["$or"] = arr{
				obj{"name": regex},
				obj{"notice": regex},
				obj{"login": regex},
			}
		}
	}
	current := c.MustGet("user").(AdminsStruct)
	if !current.Root {
		c.JSON(200, obj{"data": arr{current}})
		return
	}

	admins := []AdminsStruct{}
	if err := a.collection.Find(search).Sort("-_id").All(&admins); err != nil {
		fmt.Println("Error (admins.GetAll):", err)
	}
	c.JSON(200, obj{"data": admins})
}

// GetArr exports array of admins
func (a *Admins) GetArr() (admins []AdminsStruct) {
	if err := a.collection.Find(obj{"deleted": obj{"$ne": true}}).All(&admins); err != nil {
		fmt.Println("Error (admins.GetArr):", err)
	}
	return
}

func newStart() {
}
