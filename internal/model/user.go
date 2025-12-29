package model

import (
	"log"

	"github.com/TeamPentagon/DM-Backend/internal/database"
	"google.golang.org/protobuf/proto"
	//"time"
)

func (s *User) SaveUserData() error {
	db, err := database.CreateLevelDBDatabase("Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	bytes, err := proto.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Put([]byte(s.PersonId), bytes, nil)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (g *User) GetUserData(userName string) error {
	db, err := database.CreateLevelDBDatabase("Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	data, err := db.Get([]byte(userName), nil)
	if err != nil {
		log.Fatal(err)
	}
	err = proto.Unmarshal(data, g)
	if err != nil {
		log.Fatal(err)
	}
	return nil

}

// previous data first before calling this
func (u *User) UpdatedUserData(UserName string) error {
	db, err := database.CreateLevelDBDatabase("Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	data, err := db.Get([]byte(UserName), nil)
	if err != nil {
		log.Fatal(err)
	}
	err = proto.Unmarshal(data, u)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (d *User) DeleteUserData(UserName string) error {
	db, err := database.CreateLevelDBDatabase("Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	data, err := db.Get([]byte(UserName), nil)
	if err != nil {
		log.Fatal(err)
	}
	err = proto.Unmarshal(data, d)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
