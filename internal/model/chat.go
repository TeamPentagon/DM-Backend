package model

import (
	"fmt"

	"github.com/TeamPentagon/DM-Backend/internal/database"
	"google.golang.org/protobuf/proto"
)

func (msg *Message) SaveMessage() error {
	db, err := database.CreateLevelDBDatabase("Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	byte_data, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}
	db.Put([]byte(fmt.Sprintf("chat_%s", msg.GetMsgId())), byte_data, nil)

	return nil

}
func (msg *ChatHistory) SaveChatHistory() error {
	db, err := database.CreateLevelDBDatabase("Database/", "Common", "Shard_0.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	byte_data, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}
	db.Put([]byte(fmt.Sprintf("chat_%s", msg.GetConversationId())), byte_data, nil)

	return nil
}


func (msg *Message) Delete() error {
	//leveldb put
	db, err := database.CreateLevelDBDatabase("Database/", "NOSQL", "/")
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Delete([]byte(fmt.Sprintf("chat_%s", msg.GetMsgId())),nil)
	if err != nil {
		return err
	}
	return nil

}
func (msg *Message) Update() error {
	//leveldb put
	db, err := database.CreateLevelDBDatabase("Database/", "NOSQL", "/")
	if err != nil {
		return err
	}
	defer db.Close()
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	err = db.Delete([]byte(fmt.Sprintf("chat_%s", msg.GetMsgId())), nil)
	if err != nil {
		return err
	}

	err = db.Put([]byte(fmt.Sprintf("chat_%s", msg.GetMsgId())), data, nil)
	if err != nil {
		return err
	}
	return nil

}

// Comments is the database model for comments.
func (msg *Message) Get() error {
	//leveldb get
	db, err := database.CreateLevelDBDatabase("Database/", "NOSQL", "/")
	if err != nil {
		return err
	}
	defer db.Close()
	data, err := db.Get([]byte(fmt.Sprintf("chat_%s", msg.GetMsgId())), nil)
	if err != nil {
		return err
	}
	err = proto.Unmarshal(data, msg)
	if err != nil {
		return err
	}
	return nil
}

