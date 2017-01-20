package summer

import (
	"errors"
	"github.com/night-codes/govalidator"
	"github.com/night-codes/mgo-ai"
	"github.com/night-codes/mgo-wrapper"
	"gopkg.in/mgo.v2"
	"strings"
	"sync"
	"time"
)

type (
	UsersStruct struct {
		ID     uint64 `form:"id" json:"id" bson:"_id"`
		Login  string `form:"login" json:"login" bson:"login" valid:"required,min(3)"`
		Name   string `form:"name" json:"name" bson:"name"`
		Notice string `form:"notice" json:"notice" bson:"notice"`

		// Is Root-user? Similar as Rights.Groups = ["root"]
		Root bool `form:"-" json:"root" bson:"root"`

		// Information field, if needs auth by email set Login == Email
		Email string `form:"email" json:"email" bson:"email" valid:"email"`

		// sha512 hash of password (but from form can be received string password value)
		Password string `form:"password" json:"-" bson:"password" valid:"required,min(5)"`

		// from form can be received string password value)
		Password2 string `form:"password2" json:"-" bson:"-"`

		// Default user language (Information field)
		Lang string `form:"lang" json:"lang" bson:"lang" valid:"max(3)"`

		// Times of creating or editing (or loading from mongoDB)
		Created int64 `form:"-" json:"-" bson:"created"`
		Updated int64 `form:"-" json:"-" bson:"updated"`
		Loaded  int64 `form:"-" json:"-" bson:"-"`

		// Fields for users auth limitation
		Disabled bool `form:"-" json:"-" bson:"disabled"`
		Deleted  bool `form:"-" json:"-" bson:"deleted"`

		// User access rights (summer.Rights)
		Rights Rights `form:"-" json:"rights" bson:"rights"`

		// IP control fields (coming soon)
		LastIP   uint32 `form:"-" json:"lastIP" bson:"lastIP"`
		IP       uint32 `form:"-" json:"-" bson:"ip"`
		StringIP string `form:"-" json:"ip" bson:"-"`

		// custom data map
		Settings map[string]interface{} `form:"-" json:"settings" bson:"settings"`

		// user without authentication
		Demo bool `form:"-" json:"demo" bson:"-"`
	}
	Users struct {
		list       map[string]*UsersStruct // key - login
		listID     map[uint64]*UsersStruct // key - id
		collection *mgo.Collection
		sync.Mutex
		afterAddFn func(*UsersStruct)
		*Panel
	}
)

func (u *Users) init(panel *Panel) {
	u.Mutex = sync.Mutex{}
	u.Panel = panel
	u.collection = mongo.DB(panel.DBName).C(panel.UsersCollection)
	u.list = map[string]*UsersStruct{}
	u.listID = map[uint64]*UsersStruct{}
	u.afterAddFn = func(*UsersStruct) {}

	go func() {
		for range time.Tick(time.Second * 10) {
			u.loadUsers()
			u.clearUsers()
		}
	}()
}

// SetAddingFn set callback function that will be called after successful user adding
func (u *Users) SetAddingFn(fn func(*UsersStruct)) {
	u.Lock()
	u.afterAddFn = fn
	u.Unlock()
}

// Add new user from struct
func (u *Users) Add(user UsersStruct) error {
	if _, err := govalidator.ValidateStruct(user); err != nil {
		ers := []string{}
		for k, v := range govalidator.ErrorsByField(err) {
			ers = append(ers, k+": "+v)
		}
		return errors.New(strings.Join(ers, "<hr />"))
	}
	if user.Password != user.Password2 {
		return errors.New("Password mismatch!")
	}
	user.ID = ai.Next(u.Panel.UsersCollection)
	user.Password = H3hash(user.Password + u.Panel.AuthSalt)
	user.Created = time.Now().Unix()
	user.Updated = user.Created
	user.Loaded = user.Created

	if count, _ := u.collection.Count(); count == 0 {
		u.collection.EnsureIndex(mgo.Index{Key: []string{"login"}, Unique: true})
		u.collection.EnsureIndex(mgo.Index{Key: []string{"updated"}})
		u.collection.EnsureIndex(mgo.Index{Key: []string{"created"}})
	}

	if err := u.collection.Insert(user); err == nil {
		u.Lock()
		u.list[user.Login] = &user
		u.listID[user.ID] = &user
		u.afterAddFn(&user)
		u.Unlock()
		return nil
	} else {
		if mgo.IsDup(err) {
			return errors.New("User already exists!")
		}
		return errors.New("DB Error")
	}
}

// get changed users from mongoDB
func (u *Users) loadUsers() {
	u.Lock()
	ids := make([]uint64, len(u.listID))
	for id := range u.listID {
		ids = append(ids, id)
	}
	u.Unlock()
	now := time.Now().Unix()
	result := []UsersStruct{}
	request := obj{
		"_id": obj{"$in": ids},
		"$or": arr{
			obj{"updated": obj{"$gte": now - 30}},
			obj{"created": obj{"$gte": now - 30}},
		},
	}
	u.collection.Find(request).All(&result)

	u.Lock()
	for key, user := range result {
		result[key].Loaded = now
		u.list[user.Login] = &result[key]
		u.listID[user.ID] = &result[key]
	}
	u.Unlock()
}

// clear old records
func (u *Users) clearUsers() {
	u.Lock()
	defer u.Unlock()
	for id, user := range u.listID {
		to := time.Now().Unix() - 3660
		if user.Loaded < to || user.Deleted {
			delete(u.list, user.Login)
			delete(u.listID, id)
		}
	}
}

// GetByLogin returns user struct by login
func (u *Users) GetByLogin(login string) (user *UsersStruct, exists bool) {
	u.Lock()
	if user, exists = u.list[login]; !exists {
		u.Unlock() // Unlock 1
		result := &UsersStruct{}
		if err := u.collection.Find(obj{"login": login, "deleted": false}).One(result); err == nil {
			user = result
			user.Loaded = time.Now().Unix()
			exists = true
			u.Lock()
			u.list[user.Login] = user
			u.listID[user.ID] = user
			u.Unlock()
			return
		}
	} else {
		u.list[login].Loaded = time.Now().Unix()
		u.Unlock() // Unlock 2
		return
	}
	user = getDummyUser()
	if u.Panel.DisableAuth {
		user.Rights = Rights{
			Groups: []string{"root", "demo"},
		}
	}
	return
}

// Get returns user struct by id
func (u *Users) Get(id uint64) (user *UsersStruct, exists bool) {
	u.Lock()
	if user, exists = u.listID[id]; !exists {
		u.Unlock() // Unlock 1
		result := &UsersStruct{}
		if err := u.collection.Find(obj{"_id": id, "deleted": false}).One(result); err == nil {
			user = result
			user.Loaded = time.Now().Unix()
			exists = true
			u.Lock()
			u.list[user.Login] = user
			u.listID[user.ID] = user
			u.Unlock()
			return
		}
	} else {
		u.listID[id].Loaded = time.Now().Unix()
		u.Unlock() // Unlock 2
		return
	}
	user = getDummyUser()
	if u.Panel.DisableAuth {
		user.Rights = Rights{
			Groups: []string{"root", "demo"},
		}
	}
	return
}

// Length of users array
func (u *Users) Length() int {
	u.Lock()
	defer u.Unlock()
	return len(u.list)
}

func getDummyUser() *UsersStruct {
	return &UsersStruct{
		Name:  "",
		Login: "",
		Rights: Rights{
			Groups:  []string{"demo"},
			Actions: []string{},
		},
		Settings: obj{},
		Demo:     true,
	}
}
