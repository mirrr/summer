package summer

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mirrr/mgo-ai"
	"github.com/mirrr/mgo-wrapper"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mirrr/types.v1"
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
	}
)

var (
	admins = Admins{}
)

func (a *Admins) Init() {
	a.collection = mongo.DB(settings.DBName).C("admins")
}

func (a *Admins) Page(c *gin.Context) {
	c.HTML(200, "admins.html", gin.H{
		"title":  "Администраторы",
		"user":   c.MustGet("user").(AdminsStruct),
		"active": obj{"admins": true},
	})
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
	admin.Password = H3hash(admin.Password + settings.AuthSalt)
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
	admin.Password = H3hash(admin.Password + settings.AuthSalt)
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
		set["password"] = H3hash(admin.Password + settings.AuthSalt)
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
