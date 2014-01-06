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
	Ch   string
	Sid  int
}

type Program struct {
	Id         bson.ObjectId `bson:"_id,omitempty"`
	Hash       []byte
	EventId    string `bson:"event_id"`
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

	if n, err := collection.Find(self).Count(); err != nil {
		return err
	} else if n != 0 {
		return nil
	}

	info, err := collection.UpsertId(self.Id, self)
	if err != nil {
		return err
	}
	if info.UpsertedId != nil {
		log.Printf("Add new channel: %+v", *self)
	}
	if info.Updated > 0 {
		log.Printf("Update channel: %+v", *self)
	}
	return nil
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
		log.Printf("Add new program: %s %d %s", self.EventId, self.Start, self.Title)
	}
	if info.Updated > 0 {
		log.Printf("Update program: %s %d %s", self.EventId, self.Start, self.Title)
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
