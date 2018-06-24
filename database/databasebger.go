// Copyright Peter Lenson.
// All Rights Reserved

package database

import (
	cmnpkg "github.com/plenson/SnippetService/common"
	"github.com/zippoxer/bow"
	"log"
)

type SnipBow struct {
	Id   bow.Id
	Text string
}

type SBow struct {
	SDB bow.Bucket
}

func (self *SBow) Set(key, value string) (string, error) {
	id := bow.NewId()
	vin := SnipBow{
		Id:   id,
		Text: value,
	}

	err := self.SDB.Put(&vin)
	if err != nil {
		log.Println("Error setting")
		log.Println(err)
	}
	return id.String(), err
}

func (self *SBow) Get(key, value string) (string, error) {

	id, err := bow.ParseId(key)
	if err != nil {
		log.Println("Error parsing Id")
		log.Println(err)
	}
	vin := SnipBow{
		Id:   id,
		Text: "",
	}

	var got SnipBow
	err = self.SDB.Get(vin.Id, &got)
	if err != nil {
		log.Println("Error with getting")
		log.Println(err)
	}
	return got.Text, err
}

func (self *SBow) Del(key string) error {
	id, err := bow.ParseId(key)
	if err != nil {
		log.Println("Error parsing Id")
		log.Println(err)
	}
	vin := SnipBow{
		Id:   id,
		Text: "",
	}
	err = self.SDB.Delete(vin.Id)
	return err
}

func (self *SBow) GetAll() ([]string, error) {

	iter := self.SDB.Iter()
	defer iter.Close()

	snipids := make([]string, 0)

	var snip SnipBow
	for iter.Next(&snip) {
		snipids = append(snipids, snip.Id.String())
	}
	if iter.Err() != nil {
		log.Fatal("error")
	}

	return snipids, nil
}

func NewSnippetDBBow(params cmnpkg.DbParams) (*SBow, error) {

	// Open database under directory
	db, err := bow.Open(params.DataVolPath)
	//	bow.SetCodec(msgp.Codec{}),
	//	bow.SetBadgerOptions(badger.DefaultOptions))
	if err != nil {
		log.Fatalf("Could not create database! %s", err)
		db.Close()
		return nil, err
	}
	sDB := db.Bucket(params.HmName)

	if err != nil {
		log.Fatalf("Could not create a bucket! %s", err)
		db.Close()
		return nil, err
	}
	ss := &SBow{}
	ss.SDB = *sDB

	return ss, nil
}

func populateDataBaseBow(snipDB *SBow) error {
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
func SetupDatabaseBow(argParms cmnpkg.DbParams) (*SBow, error) {

	sBowDB, err := NewSnippetDBBow(argParms)
	if err != nil {
		log.Fatalf("Could not create a hashmap! %s", err)
		return nil, err
	}

	err = populateDataBaseBow(sBowDB)
	if err != nil {
		log.Fatalf("Problem populating database! %s", err)
		return nil, err
	}

	return sBowDB, nil
}
