package summer

import (
	"github.com/night-codes/mgo-ai"
	"github.com/night-codes/mgo-wrapper"
	"gopkg.in/mgo.v2"
	"sync"
	"time"
)

type (
	NotifyStruct struct {
		ID      uint64 `json:"id"  bson:"_id"`
		UserID  uint64 `json:"userId"  bson:"userId"`
		Title   string `json:"title" bson:"title" binding:"required,min=3"`
		Text    string `json:"text" bson:"text"`
		Created uint   `json:"-" bson:"created"`
		Updated uint   `json:"-" bson:"updated"`
		Deleted bool   `json:"-" bson:"deleted"`
		Demo    bool
	}
	notify struct {
		*Panel
		list       map[uint64]*NotifyStruct // key - login
		collection *mgo.Collection
		sync.Mutex
	}
)

func (u *notify) init(panel *Panel) {
	u.Mutex = sync.Mutex{}
	u.Panel = panel
	u.collection = mongo.DB(panel.DBName).C(panel.NotifyCollection)
	u.list = map[uint64]*NotifyStruct{}
	go func() {
		u.tick()
		for range time.Tick(time.Second * 10) {
			u.tick()
		}
	}()
}

// Add new notify from struct
func (u *notify) Add(ntf NotifyStruct) error {
	ntf.ID = ai.Next(u.Panel.NotifyCollection)
	ntf.Created = uint(time.Now().Unix() / 60)
	ntf.Updated = ntf.Created

	if err := u.collection.Insert(ntf); err == nil {
		u.Lock()
		defer u.Unlock()
		if len(u.list) == 0 {
			u.collection.EnsureIndex(mgo.Index{Key: []string{"login"}, Unique: true})
		}
		u.list[ntf.ID] = &ntf
		return nil
	} else {
		return err
	}
}

// get array of notify
func (u *notify) tick() {
	result := []NotifyStruct{}
	u.collection.Find(obj{"deleted": false}).All(&result)

	u.Lock()
	defer u.Unlock()
	u.list = map[uint64]*NotifyStruct{}
	for key, ntf := range result {
		u.list[ntf.ID] = &result[key]
	}
}

// Length of array of notify
func (u *notify) Length() int {
	u.Lock()
	defer u.Unlock()
	return len(u.list)
}
