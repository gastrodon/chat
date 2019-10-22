package io

import (
	"chat/models"
	"crypto/sha1"
	"fmt"
	"hash"
	"testing"

)

func Test_random_string(test *testing.T) {
	var size int = 32
	var rand string
	rand, _ = random_string(size)

	if len(rand) != size {
		test.Errorf("random_string size got: %d, expected: %d", len(rand), size)
	}
}

func Test_NewUser(test *testing.T) {
	var rand string
	rand, _ = random_string(4)

	var user models.User
	user, _ = NewUser("foobar", rand)

	if user != Users[user.ID] {
		test.Errorf("NewUser got: %s, expected: %s", Users[user.ID], user)
	}
}

func Test_CheckPasswd(test *testing.T) {
	var passwd string
	passwd, _ = random_string(4)

	var user models.User
	user, _ = NewUser("foobar", passwd)

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

func Test_CheckPasswdNoUser(test *testing.T) {
	var user_id string = "foobar"
	var err_string string = fmt.Sprintf("No user does exist with id %s", user_id)
	var err error
	_, err = CheckPasswd(user_id, "")

	if err.Error() != err_string {
		test.Errorf("CheckPasswd no such ID expected: %s, got: %s", err_string, err.Error())
	}
}

func Test_NewKey(test *testing.T) {
	var passwd string
	passwd, _ = random_string(4)

	var user models.User
	user, _ = NewUser("foobar", passwd)

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

func Test_NewKeyBadPasswd(test *testing.T) {
	var rand string
	rand, _ = random_string(4)

	var user models.User
	user, _ = NewUser("foobar", rand)

	var err error
	_, err = NewKey(user.ID, "")

	var err_string string = "Password passed does not match password stored"
	if err.Error() != err_string {
		test.Errorf("NewKey with bad password expected: %s, got: %s", err_string, err.Error())
	}
}

func Test_UserFromKey(test *testing.T) {
	var passwd string
	passwd, _ = random_string(4)

	var user models.User
	user, _  = NewUser("foobar", passwd)

	var key string
	key, _ = NewKey(user.ID, passwd)

	var fetched models.User
	var exists bool
	var err error
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

func Test_UserFromKeyNoUser(test *testing.T) {
	var exists bool
	var err error
	_, exists, err = UserFromKey("")

	if exists {
		test.Errorf("UserFromKey no user returns exists as true")
	}

	if err != nil {
		test.Fatal(err)
	}

}

func Test_UserFromKeyDanglingID(test *testing.T) {
	var key string = "barbaz"
	var user_id string = "foobar"
	Session[key] = user_id

	var exists bool
	var err error
	_, exists, err = UserFromKey(key)

	if exists {
		test.Errorf("UserFromKey dangling ID returns exists as true")
	}

	var err_string string = fmt.Sprintf("key %s points to user_id %s, which does not exist", key, user_id)
	if err.Error() != err_string {
		test.Errorf("UserFromKey bad key expected: %s, got %s", err_string, err.Error())
	}
}

func Test_NewRoom(test *testing.T) {
	var passwd string
	passwd, _ = random_string(4)

	var user models.User
	user, _ = NewUser("foobar", passwd)

	var room models.Room
	room, _ = NewRoom("fooroom", true, user.ID)

	if Rooms[room.ID].ID != room.ID {
		test.Errorf("NewRoom got: %s, expected: %s", Rooms[room.ID].ID, room.ID)
	}
}

func Test_UpdateUname(test *testing.T) {
	var uname string = "foobar"
	var new_uname string = "Anonymous"
	var rand string
	rand, _ = random_string(4)

	var user models.User
	user, _ = NewUser(uname, rand)

	if user.Name != uname {
		test.Errorf("NewUser expected: %s, got: %s", uname, user.Name)
	}

	UpdateUname(user.ID, new_uname)
	user = Users[user.ID]

	if user.Name != new_uname {
		test.Errorf("UpdateUname expected: %s, got: %s", new_uname, user.Name)
	}
}

func Test_UpdateUnameNoUser(test *testing.T) {
	var user_id string = "foobar"
	var err error
	_, err = UpdateUname(user_id, "")

	var err_string string = fmt.Sprintf("No user of user_id %s", user_id)
	if err.Error() != err_string {
		test.Errorf("UpdateUname no such ID expected: %s, got: %s", err_string, err.Error())
	}
}

func Test_UserFromID(test *testing.T) {
	var uname, passwd string
	uname, _ = random_string(4)
	passwd, _ =  random_string(4)
	var user models.User
	user, _ = NewUser(uname, passwd)

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

func Test_UserFromIDNoUser(test *testing.T) {
	var exists bool
	var err error
	_, exists, err = UserFromID("")

	if exists {
		test.Error("UserFromID no user_id returns exists as true")
	}

	if err != nil {
		test.Fatal(err)
	}
}
