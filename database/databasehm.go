// Copyright Peter Lenson.
// All Rights Reserved

package database

import (
	cmnpkg "github.com/plenson/SnippetService/common"
	utpkg "github.com/plenson/SnippetService/utilities"
	"github.com/xyproto/simplebolt"
	"log"
	"os"
)

type SBolt struct {
	SDB simplebolt.HashMap
}

func (self *SBolt) Set(key, value string) (string, error) {
	id := utpkg.GetUniqueID()
	err := self.SDB.Set(id, key, value)
	return id, err
}

func (self *SBolt) Get(id, key string) (string, error) {
	return self.SDB.Get(id, key)
}

func (self *SBolt) Del(id string) error {
	return self.SDB.Del(id)
}

func (self *SBolt) GetAll() ([]string, error) {
	return self.SDB.GetAll()
}

func NewSnippetDB(params cmnpkg.DbParams) (*SBolt, error) {

	// New bolt database
	db, err := simplebolt.New(params.DataVolPath + string(os.PathSeparator) + params.DbName)
	if err != nil {
		log.Fatalf("Could not create database! %s", err)
		db.Close()
		return nil, err
	}
	sDB, err := simplebolt.NewHashMap(db, params.HmName)
	if err != nil {
		log.Fatalf("Could not create a hashmap! %s", err)
		db.Close()
		return nil, err
	}

	ss := &SBolt{}
	ss.SDB = *sDB

	return ss, nil
}

func populateDataBaseBolt(snipDB *SBolt) error {
	id, err := snipDB.Set("Text", "Mary had a little lamb.")
	if err != nil {
		log.Fatalf("Could not add an item to the list! %s", err)
		return err
	}
	log.Println(id)
	id, err = snipDB.Set("Text", "Little lamb")
	if err != nil {
		log.Fatalf("Could not add an item to the list! %s", err)
		return err
	}
	log.Println(id)

	return nil
}

// Sets up the database.
//
// Currently a simple in-memory (durable) database (based on BoltDB) is used.
func SetupDatabaseBolt(argParms cmnpkg.DbParams) (*SBolt, error) {

	sBoltDB, err := NewSnippetDB(argParms)
	if err != nil {
		log.Fatalf("Could not create a hashmap! %s", err)
		return nil, err
	}

	err = populateDataBaseBolt(sBoltDB)
	if err != nil {
		log.Fatalf("Problem populating database! %s", err)
		return nil, err
	}

	return sBoltDB, nil
}
