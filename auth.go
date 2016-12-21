package summer

import (
	"fmt"
	"net/http"
	"time"

	"github.com/night-codes/mgo-ai"
	"github.com/night-codes/mgo-wrapper"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
)

type (
	UsersStruct struct {
		ID       uint64 `form:"id"  json:"id"  bson:"_id"`
		Root     bool   `form:"-"  json:"root"  bson:"root"`
		Name     string `form:"name" json:"name" bson:"name" binding:"required,min=3"`
		Notice   string `form:"notice" json:"notice" bson:"notice"`
		Login    string `form:"login" json:"login" bson:"login" binding:"required`
		Password string `form:"password" json:"-" bson:"password"`
		Updated  uint   `form:"-" json:"-" bson:"updated"`
		Deleted  bool   `form:"-" json:"-" bson:"deleted"`
		Settings interface{}
	}

	authAdmins struct {
		collection *mgo.Collection
		Panel
	}
)

var (
	admins = authAdmins{}
)

func (a *authAdmins) Init(panel *Panel) {
	a.Panel = *panel
	a.collection = mongo.DB(panel.DBName).C(a.Panel.UsersCollection)
}

func (a *authAdmins) Auth(g *gin.RouterGroup) {
	if !a.DisableAuth {
		g.Use(a.Login(g.BasePath()))
		g.POST("/z-auth", dummy) // хак для авторизации
		g.POST("/z-register", dummy)
	} else {
		g.Use(func(c *gin.Context) {
			c.Set("user", UsersStruct{})
			c.Set("login", "")
			c.Next()
		})
	}
	g.GET("/logout", a.Logout(g.BasePath()))
}

func (a *authAdmins) Logout(panelPath string) gin.HandlerFunc {
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

func (a *authAdmins) Login(panelPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminsArr := map[string]UsersStruct{}
		for _, v := range a.GetArr() {
			adminsArr[v.Login] = v
		}

		// 	регистрация первого пользователя админки
		if len(adminsArr) == 0 && !a.DisableFirstStart {
			defer c.Abort()
			login, e1 := c.GetPostForm("admin-z-login")
			password, e2 := c.GetPostForm("admin-z-password")
			password2, e3 := c.GetPostForm("admin-z-password-2")

			if e1 && e2 && e3 {
				if password == password2 {
					if len(login) > 3 && len(password) > 6 {
						a.collection.EnsureIndex(mgo.Index{Key: []string{"login"}, Unique: true})

						if err := a.Add(UsersStruct{
							Login:    login,
							Password: password,
							Name:     "Admin",
							Root:     true,
						}); err != nil {
							c.String(400, "DB Error")
							return
						}

						a.FirstStart()
						c.String(200, "Ok")
					} else {
						c.String(400, "Login or password is too short!")
					}
				} else {
					c.String(400, "Passwords do not match!")
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
					c.Set("login", user.Login)
					c.Next()
					return
				}
			}
		}
		c.HTML(200, "login.html", gin.H{"panelPath": panelPath})
		c.Abort()
	}
}

// Add new admin from struct
func (a *authAdmins) Add(admin UsersStruct) error {
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
func (a *authAdmins) GetArr() (admins []UsersStruct) {
	if err := a.collection.Find(obj{"deleted": obj{"$ne": true}}).All(&admins); err != nil {
		fmt.Println("Error (admins.GetArr):", err)
	}
	return
}
