package reserve

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
)

type Channel struct {
	Id   string `bson:"_id"`
	Name string `bson:"name"`
}

type Program struct {
	Id        bson.ObjectId `bson:"_id"`
	EventId   int           `bson:"event_id"`
	Title     string        `bson:"title"`
	Detail    string        `bson:"detail"`
	Start     int           `bson:"start"`
	End       int           `bson:"end"`
	Duration  int           `bson:"duration"`
	UpdatedAt int           `bson:"updated_at"`
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
	collection := getCollection("program")
	saved, err := GetProgram(self.EventId)
	if err != nil {
		log.Printf("Add new program: %d %d %s", self.EventId, self.Start, self.Title)
		self.Id = bson.NewObjectId()
		return collection.Insert(self)
	}
	if self.Equal(&saved) {
		return nil
	}

	log.Printf("Update program: %d %d %s", self.EventId, self.Start, self.Title)
	self.Id = saved.Id
	log.Printf("Old: %+v, New: %+v", saved, *self)
	return collection.Update(saved, self)
}

func (self *Program) Equal(other *Program) bool {
	return self.EventId == other.EventId &&
		self.Title == other.Title &&
		self.Detail == other.Detail &&
		self.Start == other.Start &&
		self.End == other.End &&
		self.Duration == other.Duration
}
