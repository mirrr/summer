package summer

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/sha3"
	"io"
	"net/http"
	"time"
)

func auth(g *gin.RouterGroup) {
	g.Use(login(g.BasePath()))
	g.POST("/z-auth", dummy) // хак для авторизации
	g.POST("/z-register", dummy)
	g.GET("/logout", logout(g.BasePath()))
}

func login(panelPath string) gin.HandlerFunc {
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
			if user, exists := adminsArr[login]; exists && user.Password == H3hash(password+settings.AuthSalt) {
				setCookie(c, settings.AuthPrefix+"login", login)
				setCookie(c, settings.AuthPrefix+"hash", H3hash(c.ClientIP()+user.Password+settings.AuthSalt))
				c.String(200, "Ok")
			} else {
				c.String(400, "Wrong password!")
			}
			c.Abort()
			return
		} else {
			login, e1 := c.Cookie(settings.AuthPrefix + "login")
			hash, e2 := c.Cookie(settings.AuthPrefix + "hash")
			if e1 == nil && e2 == nil {
				if user, exists := adminsArr[login]; exists && hash == H3hash(c.ClientIP()+user.Password+settings.AuthSalt) {
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

func dummy(c *gin.Context) {
	c.Next()
}

func logout(panelPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		http.SetCookie(c.Writer, &http.Cookie{
			Name:    settings.AuthPrefix + "hash",
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

func H3hash(s string) string {
	h3 := sha3.New512()
	io.WriteString(h3, s)
	return fmt.Sprintf("%x", h3.Sum(nil))
}

func setCookie(c *gin.Context, name string, value string) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:    name,
		Value:   value,
		Path:    "/",
		MaxAge:  32000000,
		Expires: time.Now().AddDate(1, 0, 0),
	})
}

func newStart() {
}
