package reserve

import (
	"crypto/md5"
	"fmt"
	"labix.org/v2/mgo/bson"
	"log"
)

type Rule struct {
	Id      bson.ObjectId `bson:"_id"`
	Hash    []byte
	Keyword string
}

func GetRule(id *bson.ObjectId) (*Rule, error) {
	collection := getCollection("rule")
	rule := new(Rule)
	return rule, collection.FindId(id).One(rule)
}

func ApplyAllRules(timestamp int) error {
	var rule Rule
	iter := getCollection("rule").Find(nil).Iter()
	for iter.Next(&rule) {
		if err := rule.Apply(timestamp); err != nil {
			log.Print(err)
		}
	}
	if err := iter.Close(); err != nil {
		return err
	}
	return nil
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
		self.Id = bson.NewObjectId()
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

func (self *Rule) Apply(timestamp int) error {
	if !self.Id.Valid() {
		return fmt.Errorf("Rule must be saved")
	}

	collection := getCollection("program")
	query := self.MakeQuery(timestamp)
	info, err := collection.UpdateAll(query, bson.M{"$push": bson.M{"reserved_by": self.Id}})
	if err != nil {
		return err
	} else if info.Updated > 0 {
		log.Printf("Reserved %d programs by %+v", info.Updated, *self)
	}
	return nil
}

func (self *Rule) MakeQuery(timestamp int) bson.M {
	query := bson.M{
		"reserved_by": bson.M{"$ne": self.Id},
		"title":       bson.RegEx{self.Keyword, "i"},
	}
	if timestamp != 0 {
		query["updated_at"] = timestamp
	}
	return query
}
