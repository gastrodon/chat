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
	var exists bool
	key, _ = NewKey(user.ID, passwd)

	fetched, exists, err = UserFromKey(key)

	if !exists {
		test.Fatalf("User with key %s does not exist", key)
	}

	if err != nil {
		test.Fatal(err)
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

func Test_UpdateUname(test *testing.T) {
	var uname string = "foobar"
	var new_uname string = "Anonymous"
	var passwd string = random_string(4)
	var user models.User = NewUser(uname, passwd)

	if user.Name != uname {
		test.Errorf("NewUser expected: %s, got: %s", uname, user.Name)
	}

	UpdateUname(user.ID, new_uname)
	user = Users[user.ID]

	if user.Name != new_uname {
		test.Errorf("UpdateUname expected: %s, got: %s", new_uname, user.Name)
	}
}

func Test_UserFromID(test *testing.T) {
	var uname, passwd string = random_string(4), random_string(4)
	var user models.User = NewUser(uname, passwd)

	if user.Name != uname {
		test.Errorf("NewUser expected: %s, got: %s", uname, user.Name)
	}

	var passwd_check bool
	var err error
	passwd_check, err = CheckPasswd(user.ID, passwd)

	if err != nil {
		test.Fatal(err)
	}

	if !passwd_check {
		test.Errorf("NewUser password mismatch expected: %s", passwd)
	}

	var from_id models.User
	var exists bool
	from_id, exists, err = UserFromID(user.ID)

	if err != nil {
		test.Fatal(err)
	}

	if !exists {
		test.Errorf("UserFromID no such user by id %s", user.ID)
	}

	if from_id.ID != user.ID {
		test.Errorf("UserFromID expected: %s, got: %s", user.ID, from_id.ID)
	}

}
