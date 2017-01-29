package summer

import (
	"errors"
	"github.com/kennygrant/sanitize"
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
		Name   string `form:"name" json:"name" bson:"name" valid:"max(200)"`
		Notice string `form:"notice" json:"notice" bson:"notice" valid:"max(1000)"`

		// Is Root-user? Similar as Rights.Groups = ["root"]
		Root bool `form:"-" json:"-" bson:"root"`

		// Information field, if needs auth by email set Login == Email
		Email string `form:"email" json:"email" bson:"email" valid:"email"`

		// sha512 hash of password (but from form can be received string password value)
		Password string `form:"password" json:"-" bson:"password" valid:"min(5)"`

		// from form can be received string password value)
		Password2 string `form:"password2" json:"-" bson:"-"`

		// Default user language (Information field)
		Lang string `form:"lang" json:"lang" bson:"lang" valid:"max(3)"`

		// Times of creating or editing (or loading from mongoDB)
		Created int64 `form:"-" json:"created" bson:"created"`
		Updated int64 `form:"-" json:"updated" bson:"updated"`
		Loaded  int64 `form:"-" json:"-" bson:"-"`

		// Fields for users auth limitation
		Disabled bool `form:"-" json:"disabled" bson:"disabled"`
		Deleted  bool `form:"-" json:"deleted" bson:"deleted"`

		// User access rights (summer.Rights)
		Rights Rights `json:"rights" bson:"rights"`

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
		count      int
		collection *mgo.Collection
		sync.Mutex
		afterAddFn  func(*UsersStruct)
		afterSaveFn func(*UsersStruct)
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
	u.afterSaveFn = func(*UsersStruct) {}
	u.count, _ = u.collection.Count()

	go func() {
		for range time.Tick(time.Second * 10) {
			u.count, _ = u.collection.Count()
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

// SetSavingFn set callback function that will be called after successful user saving
func (u *Users) SetSavingFn(fn func(*UsersStruct)) {
	u.Lock()
	u.afterSaveFn = fn
	u.Unlock()
}

// Add new user from struct
func (u *Users) Add(user UsersStruct) error {
	if err := u.Validate(&user); err != nil {
		return err
	}
	if len(user.Password) == 0 {
		return errors.New("Password to short!")
	}
	user.ID = ai.Next(u.Panel.UsersCollection)
	user.Name = sanitize.HTML(user.Name)
	user.Login = sanitize.HTML(user.Login)
	user.Notice = sanitize.HTML(user.Notice)
	user.Password = H3hash(user.Password + u.Panel.AuthSalt)
	user.Created = time.Now().Unix()
	user.Updated = user.Created
	user.Demo = false
	setUserDefaults(&user)

	if err := u.collection.Insert(user); err == nil {
		u.Lock()
		u.list[user.Login] = &user
		u.listID[user.ID] = &user
		u.Unlock()
		go u.afterAddFn(&user)
		return nil
	} else {
		if mgo.IsDup(err) {
			return errors.New("User already exists!")
		}
		return errors.New("DB Error")
	}
}

// Save exists user
func (u *Users) Save(user *UsersStruct) error {
	if err := u.Validate(user); err != nil {
		return err
	}
	prevUser, exists := u.Get(user.ID)
	if !exists {
		return errors.New("User not found!")
	}
	user.Login = prevUser.Login
	user.Created = prevUser.Created
	user.Name = sanitize.HTML(user.Name)
	user.Notice = sanitize.HTML(user.Notice)
	if len(user.Password) > 0 {
		user.Password = H3hash(user.Password + u.Panel.AuthSalt)
	} else {
		user.Password = prevUser.Password
	}
	user.Updated = time.Now().Unix()
	user.Demo = false
	setUserDefaults(user)

	if err := u.collection.UpdateId(user.ID, user); err == nil {
		u.Lock()
		u.list[user.Login] = user
		u.listID[user.ID] = user
		u.Unlock()
		go u.afterSaveFn(user)
		return nil
	} else {
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
			setUserDefaults(user)
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
	user = u.GetDummyUser()
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
			setUserDefaults(user)
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
	user = u.GetDummyUser()
	if u.Panel.DisableAuth {
		user.Rights = Rights{
			Groups: []string{"root", "demo"},
		}
	}
	return
}

// Length of users array
func (u *Users) Length() int {
	return u.count
}

// Length of users array
func (u *Users) CacheLength() int {
	u.Lock()
	defer u.Unlock()
	return len(u.list)
}

func (u *Users) GetDummyUser() *UsersStruct {
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

func (u *Users) Validate(user *UsersStruct) error {
	if _, err := govalidator.ValidateStruct(user); err != nil {
		ers := []string{}
		for k, v := range govalidator.ErrorsByField(err) {
			ers = append(ers, k+": "+v)
		}
		return errors.New(strings.Join(ers, " \n"))
	}
	if user.Password != user.Password2 {
		return errors.New("Password mismatch!")
	}
	return nil
}

func setUserDefaults(user *UsersStruct) {
	user.Loaded = time.Now().Unix()
	if user.Rights.Actions == nil {
		user.Rights.Actions = []string{}
	}
	if user.Rights.Groups == nil {
		user.Rights.Groups = []string{}
	}
	if user.Settings == nil {
		user.Settings = obj{}
	}
}
