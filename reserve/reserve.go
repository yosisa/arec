package reserve

import (
	"crypto/md5"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
)

type Channel struct {
	Id   string `bson:"_id"`
	Name string
}

type Program struct {
	Id         bson.ObjectId `bson:"_id,omitempty"`
	Hash       []byte
	EventId    int `bson:"event_id"`
	Title      string
	Detail     string
	Start      int
	End        int
	Duration   int
	ReservedBy []bson.ObjectId `bson:"reserved_by"`
	UpdatedAt  int             `bson:"updated_at"`
}

var session *mgo.Session

func Connect(uri string) {
	var err error
	session, err = mgo.Dial(uri)
	if err != nil {
		panic(err)
	}

	session.SetSafe(&mgo.Safe{})
}

func getCollection(collection string) *mgo.Collection {
	return session.DB("").C(collection)
}

func GetChannel(id *string) (Channel, error) {
	collection := getCollection("channel")
	var channel Channel
	return channel, collection.FindId(id).One(&channel)
}

func (self *Channel) Save() error {
	collection := getCollection("channel")
	saved, err := GetChannel(&self.Id)
	if err != nil {
		log.Printf("Add new channel: %s %s", self.Id, self.Name)
		return collection.Insert(self)
	}
	if self.Equal(&saved) {
		return nil
	}

	log.Printf("Update channel: %s %s", self.Id, self.Name)
	log.Printf("Old: %+v, New: %+v", saved, *self)
	return collection.Update(bson.M{"_id": self.Id}, self)
}

func (self *Channel) Equal(other *Channel) bool {
	return self.Id == other.Id &&
		self.Name == other.Name
}

func GetProgram(event_id int) (Program, error) {
	collection := getCollection("program")
	var program Program
	return program, collection.Find(bson.M{"event_id": event_id}).One(&program)
}

func (self *Program) Save() error {
	self.Hash = self.MakeHash()
	collection := getCollection("program")

	// check duplication
	if n, err := collection.Find(bson.M{"hash": self.Hash}).Count(); err != nil {
		return err
	} else if n != 0 {
		return nil
	}

	info, err := collection.Upsert(bson.M{"event_id": self.EventId}, self)
	if err != nil {
		return err
	}
	if info.UpsertedId != nil {
		log.Printf("Add new program: %d %d %s", self.EventId, self.Start, self.Title)
	}
	if info.Updated > 0 {
		log.Printf("Update program: %d %d %s", self.EventId, self.Start, self.Title)
	}
	return nil
}

func (self *Program) MakeHash() []byte {
	hasher := md5.New()
	fmt.Fprintf(hasher, "%v", []interface{}{
		self.EventId,
		self.Title,
		self.Detail,
		self.Start,
		self.End,
		self.Duration,
	})
	return hasher.Sum(nil)
}
