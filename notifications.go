package summer

import (
	"sync"
	"time"

	"github.com/night-codes/mgo-wrapper"
	"gopkg.in/mgo.v2"
)

type (
	// NotifyStruct data struct
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

func (n *notify) init(panel *Panel) {
	n.Mutex = sync.Mutex{}
	n.Panel = panel
	n.collection = mongo.DB(panel.DBName).C(panel.NotifyCollection)
	n.list = map[uint64]*NotifyStruct{}
	go func() {
		n.tick()
		for range time.Tick(time.Second * 10) {
			n.tick()
		}
	}()
}

// Add new notify from struct
func (n *notify) Add(ntf NotifyStruct) error {
	ntf.ID = n.AI.Next(n.Panel.NotifyCollection)
	ntf.Created = uint(time.Now().Unix() / 60)
	ntf.Updated = ntf.Created

	if err := n.collection.Insert(ntf); err != nil {
		return err
	}

	n.Lock()
	defer n.Unlock()
	if len(n.list) == 0 {
		n.collection.EnsureIndex(mgo.Index{Key: []string{"login"}, Unique: true})
	}
	n.list[ntf.ID] = &ntf
	return nil
}

// get array of notify
func (n *notify) tick() {
	result := []NotifyStruct{}
	n.collection.Find(obj{"deleted": false}).All(&result)

	n.Lock()
	defer n.Unlock()
	n.list = map[uint64]*NotifyStruct{}
	for key, ntf := range result {
		n.list[ntf.ID] = &result[key]
	}
}

// Length of array of notify
func (n *notify) Length() int {
	n.Lock()
	defer n.Unlock()
	return len(n.list)
}
