package summer

import (
	"net/http"
	"strings"
	"time"

	"github.com/night-codes/mgo-wrapper"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
)

type (
	auth struct {
		collection *mgo.Collection
		*Panel
	}
)

func (a *auth) init(panel *Panel) {
	a.Panel = panel
	a.collection = mongo.DB(panel.DBName).C(a.Panel.UsersCollection)
}

func (a *auth) Auth(g *gin.RouterGroup) {
	if !a.DisableAuth {
		g.Use(a.Login(g.BasePath()))
		g.POST("/z-auth", dummy) // хак для авторизации
		g.POST("/z-register", dummy)
	} else {
		g.Use(func(c *gin.Context) {
			c.Set("user", getDummyUser("demo"))
			c.Set("login", "demo")
			c.Next()
		})
	}
	g.GET("/logout", a.Logout(g.BasePath()))
}

func (a *auth) Logout(panelPath string) gin.HandlerFunc {
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

func (a *auth) Login(panelPath string) gin.HandlerFunc {
	return func(c *gin.Context) {

		// 	регистрация первого пользователя админки
		if a.users.Length() == 0 && !a.DisableFirstStart {
			defer c.Abort()
			login, e1 := c.GetPostForm("admin-z-login")
			password, e2 := c.GetPostForm("admin-z-password")
			password2, e3 := c.GetPostForm("admin-z-password-2")

			if e1 && e2 && e3 {
				if password == password2 {
					if len(login) > 2 && len(password) > 5 {
						if err := a.users.Add(UsersStruct{
							Login:    login,
							Password: password,
							Name:     strings.Title(login),
							Root:     true,
							Rights:   Rights{Groups: []string{"root"}, Actions: []string{"all"}},
							Settings: obj{},
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

			c.HTML(200, "firstStart.html", gin.H{"panelPath": panelPath})
			c.Abort()
			return
		}

		// авторизация пользователя админки
		login, e1 := c.GetPostForm("admin-z-login")
		password, e2 := c.GetPostForm("admin-z-password")
		if e1 && e2 {
			if user := a.users.GetByLogin(login); user.Password == H3hash(password+a.AuthSalt) {
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
				if user := a.users.GetByLogin(login); hash == H3hash(c.ClientIP()+user.Password+a.AuthSalt) {
					if user.Root {
						user.Rights.Groups = uniqAppend(user.Rights.Groups, []string{"root"})
					}

					c.Set("user", *user)
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
