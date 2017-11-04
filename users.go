package summer

import (
	"errors"
	"github.com/kennygrant/sanitize"
	"github.com/night-codes/govalidator"
	"github.com/night-codes/mgo-ai"
	"github.com/night-codes/mgo-wrapper"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
		Password2 string `form:"password2" json:"-" bson:"password2"`

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
		rawList    map[string]*bson.Raw    // key - login
		rawListID  map[uint64]*bson.Raw    // key - id
		list       map[string]*UsersStruct // key - login
		listID     map[uint64]*UsersStruct // key - id
		count      int
		collection *mgo.Collection
		sync.Mutex
		mutUsers sync.Mutex
		*Panel
	}
)

func UsersFarm(DBName, UsersCollection, AuthSalt string, AICollection ...string) *Users {
	AIColl := "ai"
	if len(AICollection) > 0 && len(AICollection[0]) > 0 {
		AIColl = AICollection[0]
	}
	if len(UsersCollection) == 0 {
		UsersCollection = "admins"
	}
	if len(DBName) == 0 {
		DBName = "summerPanel"
	}
	if len(AuthSalt) == 0 {
		AuthSalt = "+Af761"
	}
	farm := &Panel{
		Settings: Settings{
			AuthSalt:        AuthSalt,
			DBName:          DBName,
			UsersCollection: UsersCollection,
			AICollection:    AIColl,
		},
		Users: new(Users),
		AI:    *ai.Create(mongo.DB(DBName).C(AIColl)),
	}
	farm.Users.init(farm)
	return farm.Users
}

func (u *Users) init(panel *Panel) {
	u.Mutex = sync.Mutex{}
	u.Panel = panel
	u.collection = mongo.DB(panel.DBName).C(panel.UsersCollection)
	u.rawList = map[string]*bson.Raw{}
	u.rawListID = map[uint64]*bson.Raw{}
	u.list = map[string]*UsersStruct{}
	u.listID = map[uint64]*UsersStruct{}
	u.count, _ = u.collection.Count()

	go func() {
		for range time.Tick(time.Second * 10) {
			u.count, _ = u.collection.Count()
			u.loadUsers()
			u.clearUsers()
		}
	}()
}

// Add new user from struct
func (u *Users) Add(user UsersStruct) error {
	if err := u.Validate(&user); err != nil {
		return err
	}
	if len(user.Password) == 0 {
		return errors.New("Password to short!")
	}
	user.ID = u.AI.Next(u.Panel.UsersCollection)
	user.Name = sanitize.HTML(user.Name)
	user.Login = sanitize.HTML(user.Login)
	user.Notice = sanitize.HTML(user.Notice)
	user.Password = u.Panel.HashDBFn(user.Login, user.Password, u.Panel.AuthSalt)
	user.Created = time.Now().Unix()
	user.Updated = user.Created
	user.Demo = false
	setUserDefaults(&user)
	msh, err := bson.Marshal(user)
	if err != nil {
		return err
	}
	rawUser := bson.Raw{Kind: 3, Data: msh}
	if err := u.collection.Insert(user); err == nil {
		u.Lock()
		u.list[user.Login] = &user
		u.listID[user.ID] = &user
		u.rawList[user.Login] = &rawUser
		u.rawListID[user.ID] = &rawUser
		u.Unlock()
		return nil
	} else {
		if mgo.IsDup(err) {
			return errors.New("User already exists!")
		}
		return errors.New("DB Error")
	}
}

// Add new user from struct
func (u *Users) AddFrom(data interface{}) (error, uint64) {
	msh, err := bson.Marshal(data)
	if err != nil {
		return err, 0
	}
	rawUser := bson.Raw{Kind: 3, Data: msh}

	user := &UsersStruct{}
	if err := rawUser.Unmarshal(&user); err != nil {
		return err, 0
	}

	if err := u.Validate(user); err != nil {
		return err, 0
	}
	if len(user.Password) == 0 {
		return errors.New("Password to short!"), 0
	}
	user.ID = u.AI.Next(u.Panel.UsersCollection)
	user.Name = sanitize.HTML(user.Name)
	user.Login = sanitize.HTML(user.Login)
	user.Email = sanitize.HTML(user.Email)
	user.Notice = sanitize.HTML(user.Notice)
	user.Password = u.Panel.HashDBFn(user.Login, user.Password, u.Panel.AuthSalt)
	user.Created = time.Now().Unix()
	user.Updated = user.Created
	user.Demo = false
	setUserDefaults(user)

	insert := obj{}
	rawUser.Unmarshal(&insert)
	insert["_id"] = user.ID

	errIns := u.collection.Insert(insert)
	if err := u.collection.UpdateId(user.ID, obj{"$set": user}); err == nil && errIns == nil {
		u.Get(user.ID)
		return nil, user.ID
	} else {
		if mgo.IsDup(errIns) {
			return errors.New("User already exists!"), 0
		}
		return errors.New("DB Error"), 0
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
	user.Email = sanitize.HTML(user.Email)
	user.Notice = sanitize.HTML(user.Notice)
	if len(user.Password) > 0 {
		user.Password = u.Panel.HashDBFn(user.Login, user.Password, u.Panel.AuthSalt)
	} else {
		user.Password = prevUser.Password
	}
	user.Updated = time.Now().Unix()
	user.Demo = false
	setUserDefaults(user)
	msh, err := bson.Marshal(user)
	if err != nil {
		return err
	}
	rawUser := bson.Raw{Kind: 3, Data: msh}

	if err := u.collection.UpdateId(user.ID, user); err == nil {
		u.Lock()
		u.list[user.Login] = user
		u.listID[user.ID] = user
		u.rawList[user.Login] = &rawUser
		u.rawListID[user.ID] = &rawUser
		u.Unlock()
		return nil
	} else {
		return errors.New("DB Error")
	}
}

// Save exists user from own struct
func (u *Users) SaveFrom(data interface{}) error {
	msh, err := bson.Marshal(data)
	if err != nil {
		return err
	}
	rawUser := bson.Raw{Kind: 3, Data: msh}
	user := &UsersStruct{}
	if err := rawUser.Unmarshal(user); err != nil {
		return err
	}
	if err := u.Validate(user); err != nil {
		return err
	}
	if user.ID == 0 {
		return errors.New("Wrong ID")
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
		user.Password = u.Panel.HashDBFn(user.Login, user.Password, u.Panel.AuthSalt)
	} else {
		user.Password = prevUser.Password
	}
	user.Updated = time.Now().Unix()
	user.Demo = false
	setUserDefaults(user)

	u.collection.UpdateId(user.ID, obj{"$set": data})
	if err := u.collection.UpdateId(user.ID, obj{"$set": user}); err == nil {
		u.Lock()
		u.list[user.Login] = user
		u.listID[user.ID] = user
		u.rawList[user.Login] = &rawUser
		u.rawListID[user.ID] = &rawUser
		u.Unlock()
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
	resultRaw := []bson.Raw{}
	request := obj{
		"_id": obj{"$in": ids},
		"$or": arr{
			obj{"updated": obj{"$gte": now - 60}},
			obj{"created": obj{"$gte": now - 60}},
		},
	}
	u.collection.Find(request).All(&resultRaw)

	u.Lock()
	for key, _ := range resultRaw {
		user := UsersStruct{}
		resultRaw[key].Unmarshal(&user)
		user.Loaded = now
		u.list[user.Login] = &user
		u.listID[user.ID] = &user
		u.rawList[user.Login] = &resultRaw[key]
		u.rawListID[user.ID] = &resultRaw[key]
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
			delete(u.rawList, user.Login)
			delete(u.listID, id)
			delete(u.rawListID, id)
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

// FetchByLogin returns user data by login
func (u *Users) GetByLoginTo(login string, user interface{}) (exists bool) {
	rawUser := &bson.Raw{}
	u.Lock() // Lock 1
	if rawUser, exists = u.rawList[login]; !exists {
		u.Unlock() // Unlock 1-1 (IF)
		rawUser = &bson.Raw{}
		if err := u.collection.Find(obj{"login": login, "deleted": false}).One(rawUser); err == nil {
			result := &UsersStruct{}
			rawUser.Unmarshal(result)

			setUserDefaults(result)
			exists = true
			u.Lock() // Lock 2
			u.list[result.Login] = result
			u.listID[result.ID] = result
			u.rawList[result.Login] = rawUser
			u.rawListID[result.ID] = rawUser
			u.Unlock() // Unlock 2
		} else {
			return
		}
	} else {
		exists = true
		u.list[login].Loaded = time.Now().Unix()
		u.Unlock() // Unlock 1-2 (ELSE)
	}

	rawUser.Unmarshal(user)
	return
}

// Get returns user struct by id
func (u *Users) Get(id uint64) (user *UsersStruct, exists bool) {
	u.Lock() // Lock 1
	if user, exists = u.listID[id]; !exists {
		u.Unlock() // Unlock 1-1
		result := &UsersStruct{}
		if err := u.collection.Find(obj{"_id": id, "deleted": false}).One(result); err == nil {
			user = result
			setUserDefaults(user)
			exists = true
			u.Lock() // Lock 2
			u.list[user.Login] = user
			u.listID[user.ID] = user
			u.Unlock() // Unlock 2
			return
		}
	} else {
		u.listID[id].Loaded = time.Now().Unix()
		u.Unlock() // Unlock 1-2
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

// Fetch returns user data by login
func (u *Users) GetTo(id uint64, user interface{}) (exists bool) {
	rawUser := &bson.Raw{}
	u.Lock() // Lock 1
	if rawUser, exists = u.rawListID[id]; !exists {
		u.Unlock() // Unlock 1-1 (IF)
		rawUser = &bson.Raw{}
		if err := u.collection.Find(obj{"_id": id, "deleted": false}).One(rawUser); err == nil {
			result := &UsersStruct{}
			rawUser.Unmarshal(result)

			setUserDefaults(result)
			exists = true
			u.Lock() // Lock 2
			u.list[result.Login] = result
			u.listID[result.ID] = result
			u.rawList[result.Login] = rawUser
			u.rawListID[result.ID] = rawUser
			u.Unlock() // Unlock 2
		} else {
			return
		}
	} else {
		u.listID[id].Loaded = time.Now().Unix()
		u.Unlock() // Unlock 1-2 (ELSE)
	}

	rawUser.Unmarshal(user)
	return
}

func (u *Users) GetFromContextTo(c *gin.Context, user interface{}) (exists bool) {
	u.mutUsers.Lock()
	defer u.mutUsers.Unlock()
	exists = u.GetTo((c.MustGet("user").(UsersStruct)).ID, user)
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
	user.Password2 = ""
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
