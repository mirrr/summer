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
		added        bool
		collection   *mgo.Collection
		fsCollection *mgo.Collection
		fsCount      int
		*Panel
	}
)

func (a *auth) init(panel *Panel) {
	a.Panel = panel
	a.fsCount = -1 // count of records in the collection "firstStart" (if -1, have not looked in the session)
	a.collection = mongo.DB(panel.DBName).C(a.Panel.UsersCollection)
	a.fsCollection = mongo.DB(panel.DBName).C("firstStart")
}

func (a *auth) Auth(g *gin.RouterGroup, disableAuth bool) {
	if !a.DisableAuth {
		middle := a.Login(g.BasePath(), disableAuth)
		g.Use(middle)
	} else {
		g.Use(func(c *gin.Context) {
			user := getDummyUser()
			user.Rights = Rights{
				Groups: []string{"root"},
			}
			c.Set("user", *user)
			c.Set("login", "")
			c.Next()
		})
	}
	if !a.DisableAuth && !a.added {
		a.RouterGroup.GET("/logout", a.Logout(a.RouterGroup.BasePath()))
		middle := a.Login(g.BasePath(), false)
		authGroup := a.RouterGroup.Group("/summer-auth")
		authGroup.Use(middle)
		authGroup.POST("/login", dummy)
		authGroup.POST("/register", dummy)
		a.added = true
	}
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
		c.Abort()
	}
}

func (a *auth) Login(panelPath string, disableAuth bool) gin.HandlerFunc {
	return func(c *gin.Context) {

		// 	First Start
		if a.fsCount == -1 { // have not looked in the session
			a.fsCount, _ = a.fsCollection.Find(obj{"uc": a.UsersCollection}).Count()
		}
		if a.fsCount <= 0 {
			if !disableAuth && !a.DisableFirstStart {
				defer c.Abort()
				login, e1 := c.GetPostForm("admin-z-login")
				password, e2 := c.GetPostForm("admin-z-password")
				password2, e3 := c.GetPostForm("admin-z-password-2")

				if e1 && e2 && e3 {
					if err := a.Users.Add(UsersStruct{
						Login:     login,
						Password:  password,
						Password2: password2,
						Name:      strings.Title(login),
						Root:      true,
						Rights:    Rights{Groups: []string{"root"}, Actions: []string{"all"}},
						Settings:  obj{},
					}); err != nil {
						c.String(400, err.Error())
						return
					}

					a.fsCollection.Insert(obj{"uc": a.UsersCollection, "commit": true})
					a.fsCount = 1
					go a.FirstStart()
					c.String(200, "Ok")
					return
				}

				c.HTML(200, "firstStart.html", gin.H{"panelPath": panelPath})
				c.Abort()
				return
			} else {
				a.fsCollection.Insert(obj{"uc": a.UsersCollection, "commit": true})
				a.fsCount = 1
				go a.FirstStart()
			}
		}

		// авторизация пользователя админки
		login, e1 := c.GetPostForm("admin-z-login")
		password, e2 := c.GetPostForm("admin-z-password")
		if !disableAuth && e1 && e2 {
			if user, exists := a.Users.GetByLogin(login); exists && user.Password == H3hash(password+a.AuthSalt) {
				if !user.Disabled && !user.Deleted {
					setCookie(c, a.AuthPrefix+"login", login)
					setCookie(c, a.AuthPrefix+"hash", H3hash(c.ClientIP()+user.Password+a.AuthSalt))
					c.String(200, "Ok")
				} else {
					c.String(400, "Account disabled or waits for moderation.")
				}
			} else {
				c.String(400, "Wrong password!")
			}
			c.Abort()
			return
		} else {
			login, e1 := c.Cookie(a.AuthPrefix + "login")
			hash, e2 := c.Cookie(a.AuthPrefix + "hash")
			if e1 == nil && e2 == nil {
				if user, exists := a.Users.GetByLogin(login); exists && hash == H3hash(c.ClientIP()+user.Password+a.AuthSalt) {
					if !user.Disabled && !user.Deleted {
						if user.Root {
							user.Rights.Groups = uniqAppend(user.Rights.Groups, []string{"root"})
						}

						c.Set("user", *user)
						c.Set("login", user.Login)
						c.Next()
						return
					} else {
						a.Logout(a.RouterGroup.BasePath())(c)
						return
					}
				}
			}
			if disableAuth {
				user := getDummyUser()
				c.Set("user", *user)
				c.Set("login", "")
				c.Next()
				return
			}
		}
		c.HTML(200, "login.html", gin.H{"panelPath": panelPath})
		c.Abort()
	}
}
