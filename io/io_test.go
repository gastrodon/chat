package io

import (
	"chat/models"
	"crypto/sha1"
	"hash"
	"testing"
)

func Test_random_string(test *testing.T) {
	var size int = 32
	var rand string = random_string(size)

	if len(rand) != size {
		test.Errorf("random_string size got: %d, expected: %d", len(rand), size)
	}
}

func Test_NewUser(test *testing.T) {
	var passwd string = random_string(4)
	var user models.User = NewUser("foobar", passwd)

	if user != Users[user.ID] {
		test.Errorf("NewUser got: %s, expected: %s", Users[user.ID], user)
	}
}

func Test_CheckPasswd(test *testing.T) {
	var passwd string = random_string(4)
	var user models.User = NewUser("foobar", passwd)
	var check bool
	var err error
	check, err = CheckPasswd(user.ID, passwd)

	if err != nil {
		test.Error(err)
	}

	if !check {
		var hashed hash.Hash = sha1.New()
		hashed.Write([]byte(salt + passwd))
		test.Errorf("CheckPasswd got: %s", hashed.Sum(nil))
	}
}

func Test_NewKey(test *testing.T) {
	var passwd string = random_string(4)
	var user models.User = NewUser("foobar", passwd)
	var key string
	var err error
	key, err = NewKey(user.ID, passwd)

	if err != nil {
		test.Error(err)
	}

	if len(key) != 32 {
		test.Errorf("NewKey len got: %d, expected %d", len(key), 32)
	}
}

func Test_UserFromKey(test *testing.T) {
	var passwd string = random_string(4)
	var user models.User = NewUser("foobar", passwd)
	var key string
	var fetched models.User
	var err error
	key, _ = NewKey(user.ID, passwd)

	fetched, err = UserFromKey(key)

	if err != nil {
		test.Error(err)
	}

	if fetched.ID != user.ID {
		test.Errorf("UserFromKey got: %s, expected: %s", fetched.ID, user.ID)
	}
}

func Test_NewRoom(test *testing.T) {
	var passwd string = random_string(4)
	var user models.User = NewUser("foobar", passwd)
	var room models.Room = NewRoom("fooroom", true, user.ID)

	if Rooms[room.ID].ID != room.ID {
		test.Errorf("NewRoom got: %s, expected: %s", Rooms[room.ID].ID, room.ID)
	}
}
