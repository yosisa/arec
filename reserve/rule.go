package reserve

import (
	"crypto/md5"
	"fmt"
	"labix.org/v2/mgo/bson"
	"log"
)

type Rule struct {
	Id      bson.ObjectId `bson:"_id,omitempty"`
	Hash    []byte
	Keyword string
}

func GetRule(id *bson.ObjectId) (*Rule, error) {
	collection := getCollection("rule")
	rule := new(Rule)
	return rule, collection.FindId(id).One(rule)
}

func (self *Rule) Save() error {
	self.Hash = self.MakeHash()
	collection := getCollection("rule")

	// check duplication
	n, err := collection.Find(bson.M{"hash": self.Hash}).Count()
	if err != nil {
		return err
	} else if n != 0 {
		return fmt.Errorf("Duplicated rule: %+v", *self)
	}

	if !self.Id.Valid() {
		log.Printf("Add new rule: %+v", *self)
		return collection.Insert(self)
	} else {
		log.Printf("Update rule: %+v", *self)
		return collection.UpdateId(self.Id, self)
	}
}

func (self *Rule) MakeHash() []byte {
	hasher := md5.New()
	hasher.Write([]byte(self.Keyword))
	return hasher.Sum(nil)
}
