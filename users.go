package summer

import (
	"github.com/night-codes/mgo-ai"
	"github.com/night-codes/mgo-wrapper"
	"gopkg.in/mgo.v2"
	"sync"
	"time"
)

type (
	UsersStruct struct {
		ID       uint64                 `form:"id"  json:"id"  bson:"_id"`
		Root     bool                   `form:"-"  json:"root"  bson:"root"`
		Name     string                 `form:"name" json:"name" bson:"name" binding:"required,min=3"`
		Notice   string                 `form:"notice" json:"notice" bson:"notice"`
		Login    string                 `form:"login" json:"login" bson:"login" binding:"required`
		Password string                 `form:"password" json:"-" bson:"password"`
		Created  uint                   `form:"-" json:"-" bson:"created"`
		Updated  uint                   `form:"-" json:"-" bson:"updated"`
		Deleted  bool                   `form:"-" json:"-" bson:"deleted"`
		Rights   Rights                 `form:"-" json:"rights" bson:"rights"`
		Settings map[string]interface{} `form:"-" json:"settings" bson:"-"`
	}
	users struct {
		list       map[string]*UsersStruct // key - login
		collection *mgo.Collection
		Panel
		sync.Mutex
	}
)

var (
	Users = users{}
)

func (u *users) init(panel *Panel) {
	u.Panel = *panel
	u.collection = mongo.DB(panel.DBName).C(u.Panel.UsersCollection)
	u.list = map[string]*UsersStruct{}
	go func() {
		u.tick()
		for range time.Tick(time.Second * 10) {
			u.tick()
		}
	}()
}

// Add new user from struct
func (u *users) Add(user UsersStruct) error {
	user.ID = ai.Next(u.UsersCollection)
	user.Password = H3hash(user.Password + u.AuthSalt)
	user.Created = uint(time.Now().Unix() / 60)
	user.Updated = user.Created

	if err := u.collection.Insert(user); err == nil {
		u.Lock()
		defer u.Unlock()
		if len(u.list) == 0 {
			u.collection.EnsureIndex(mgo.Index{Key: []string{"login"}, Unique: true})
		}
		u.list[user.Login] = &user
		return nil
	} else {
		return err
	}
}

// get array of users
func (u *users) tick() {
	result := []UsersStruct{}
	u.collection.Find(obj{"deleted": false}).All(&result)

	u.Lock()
	defer u.Unlock()
	u.list = map[string]*UsersStruct{}
	for key, user := range result {
		u.list[user.Login] = &result[key]
	}
}

// GetArr exports array of users
func (u *users) GetByLogin(login string) (user *UsersStruct, exists bool) {
	u.Lock()
	defer u.Unlock()
	user, exists = u.list[login]
	return
}

// Length of array of users
func (u *users) Length() int {
	u.Lock()
	defer u.Unlock()
	return len(u.list)
}

func getByLogin(login string) *UsersStruct {
	user, _ := Users.GetByLogin(login)
	return user
}
