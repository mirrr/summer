package summer

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/night-codes/mgo-wrapper"
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

		if !a.added {
			a.RouterGroup.GET("/logout", a.Logout(a.RouterGroup.BasePath()))
			middle := a.Login(g.BasePath(), false)
			authGroup := a.RouterGroup.Group("/summer-auth")
			authGroup.Use(middle)
			authGroup.POST("/login", dummy)
			authGroup.POST("/register", dummy)
			a.added = true
		}
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
			FS := func() {
				a.Users.collection.EnsureIndex(mgo.Index{Key: []string{"login"}, Unique: true})
				a.Users.collection.EnsureIndex(mgo.Index{Key: []string{"updated"}})
				a.Users.collection.EnsureIndex(mgo.Index{Key: []string{"created"}})
				a.fsCollection.Insert(obj{"uc": a.UsersCollection, "commit": true})
				a.fsCount = 1
				go a.FirstStart()
			}
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
					FS()
					c.String(200, "Ok")
					return
				}

				c.HTML(200, "firstStart.html", gin.H{"panelPath": panelPath})
				c.Abort()
				return
			}
			FS()
		}

		// авторизация пользователя админки
		login, e1 := c.GetPostForm("admin-z-login")
		password, e2 := c.GetPostForm("admin-z-password")
		cliIP := c.ClientIP()
		userAgent := c.Request.Header.Get("User-Agent")
		if a.Panel.AuthSkipIP {
			cliIP = ""
		}
		if !disableAuth && e1 && e2 {
			if user, exists := a.Users.GetByLogin(login); exists && user.Password == a.HashDBFn(login, password, a.AuthSalt) {
				if !user.Disabled && !user.Deleted {
					setCookie(c, a.AuthPrefix+"login", login)
					setCookie(c, a.AuthPrefix+"hash", a.HashCookieFn(login, user.Password, a.AuthSalt, cliIP, userAgent))
					c.String(200, "Ok")
				} else {
					c.String(401, "Account disabled or waits for moderation.")
				}
			} else {
				c.String(401, "Wrong password!")
			}
			c.Abort()
			return
		}
		login, err1 := c.Cookie(a.AuthPrefix + "login")
		hash, err2 := c.Cookie(a.AuthPrefix + "hash")
		if err1 == nil && err2 == nil {
			if user, exists := a.Users.GetByLogin(login); exists && hash == a.HashCookieFn(login, user.Password, a.AuthSalt, cliIP, userAgent) {
				if !user.Disabled && !user.Deleted {
					if user.Root {
						user.Rights.Groups = uniqAppend(user.Rights.Groups, []string{"root"})
					}

					c.Set("user", *user)
					c.Set("userID", user.ID)
					c.Set("login", user.Login)
					c.Next()
					return
				}
				a.Logout(a.RouterGroup.BasePath())(c)
				return

			}
		}
		if disableAuth {
			user := a.Users.GetDummyUser()
			c.Set("user", *user)
			c.Set("userID", 0)
			c.Set("login", "")
			c.Next()
			return
		}

		if strings.Contains(c.Request.Header.Get("Upgrade"), "websocket") ||
			strings.Contains(c.Request.Header.Get("X-Requested-With"), "XMLHttpRequest") {

			c.String(403, "Unauthorized request")
		} else {
			c.HTML(200, "login.html", gin.H{"panelPath": panelPath})
		}
		c.Abort()
	}
}
