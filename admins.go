package summer

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mirrr/mgo-ai"
	"github.com/mirrr/mgo-wrapper"
	"gopkg.in/mgo.v2"
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
		for _, v := range a.GetArr() {
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
						if err := a.AddRaw(AdminsStruct{
							Login:    login,
							Password: password,
							Name:     "Admin",
							Root:     true,
						}); err != nil {
							c.String(400, "Ошибка БД")
							return
						}
						a.FirstStart()
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

// AddRaw adds new admin from struct
func (a *Admins) AddRaw(admin AdminsStruct) error {
	admin.ID = ai.Next("admins")
	admin.Password = H3hash(admin.Password + a.AuthSalt)
	admin.Updated = uint(time.Now().Unix() / 60)
	if err := a.collection.Insert(admin); err != nil {
		ai.Cancel("admins")
		return err
	}
	return nil
}

// GetArr exports array of admins
func (a *Admins) GetArr() (admins []AdminsStruct) {
	if err := a.collection.Find(obj{"deleted": obj{"$ne": true}}).All(&admins); err != nil {
		fmt.Println("Error (admins.GetArr):", err)
	}
	return
}
