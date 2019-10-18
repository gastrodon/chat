package io

import (
	"chat/models"
	"crypto/rand"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"hash"
	"io"
)

var Rooms map[string]models.Room
var Users map[string]models.User
var Session map[string]string
var Login map[string][]byte
var salt string

func init() {
	assertPRNG()
	Rooms = make(map[string]models.Room)
	Users = make(map[string]models.User)
	Session = make(map[string]string)
	Login = make(map[string][]byte)
	salt = random_string(16)
}

func assertPRNG() {
	var err error
	var buffer []byte = make([]byte, 1)
	_, err = io.ReadFull(rand.Reader, buffer)

	if err != nil {
		panic(err)
	}
}

func random_string(size int) string {
	const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_"
	var byte_array []byte = make([]byte, size)
	rand.Read(byte_array)

	var index int
	var value byte
	for index, value = range byte_array {
		byte_array[index] = alphabet[value%byte(len(alphabet))]
	}

	return string(byte_array)
}

func UpdateUname(user_id string, uname string) (models.User, error){
	var user models.User
	var exists bool
	user, exists = Users[user_id]

	if !exists {
		var err_string string = fmt.Sprintf("No user of user_id %s", user_id)
		return user, errors.New(err_string)
	}

	user.Name = uname
	Users[user_id] = user
	return Users[user_id], nil
}

/**
 * Create a user
 * uname    string  -> Username of this user. Request handlers should take care of name defaults.
 * passwd   string  -> Password string. A good client implementation should be hashing passwords before sending,
 *                  resulting in a double password hash. This second hash is salted with a random salt determined at runtime.
 *
 * return   User    -> Created `User`
 */
func NewUser(uname string, passwd string) models.User {
	var id string = uuid.NewV4().String()
	var hashed hash.Hash = sha1.New()
	hashed.Write([]byte(salt + passwd))

	var user models.User = models.User{
		Name: uname,
		ID:   id,
	}

	Login[id] = hashed.Sum(nil)
	Users[id] = user
	return user
}

/**
 * Check that some password for some user_id is correct
 * user_id  string  -> UUID of this user. Should exists in `Users` map
 * passwd   string  -> Suspected password of this user.
 *
 * return   bool    -> Does the given password belong to the given user?
 *          error   -> No such user by this ID does exist
 */
func CheckPasswd(user_id string, passwd string) (bool, error) {
	var exists bool
	_, exists = Users[user_id]

	if !exists {
		return false, errors.New(fmt.Sprintf("No user does exist with id %s", user_id))
	}

	var hashed hash.Hash = sha1.New()
	hashed.Write([]byte(salt + passwd))
	return string(hashed.Sum(nil)) == string(Login[user_id]), nil
}

/**
 * Generate a session key for some user
 * user_id  string  -> UUID of this user. Should exist in `Users` map
 * passwd   string  -> Password of this user
 *
 * return   string  -> Session key of this user
 *          error   -> Password supplied is not correct
 *                  -> No such user by this ID does exist
 */
func NewKey(user_id string, passwd string) (string, error) {
	var is_passwd bool
	var err error
	is_passwd, err = CheckPasswd(user_id, passwd)

	if err != nil {
		return "", err
	}

	if !is_passwd {
		return "", errors.New("Password passed does not match password stored")
	}

	var key string = random_string(32)
	Session[key] = user_id
	return key, nil
}

/**
 * Get the user that is authenticated by some key
 * key      string  -> Key to find the owner of
 *
 * return   User    -> `User` that does own this key
 *          error   -> No `User` does own this key
 *                  -> This key is registered, but points to nothing
 */
func UserFromKey(key string) (models.User, bool, error) {
	var exists bool
	var user_id string
	var user models.User
	user_id, exists = Session[key]

	if !exists {
		return user, false, nil
	}

	user, exists = Users[user_id]

	if !exists {
		return user, false, errors.New(fmt.Sprintf("key %s points to user_id %s, which does not exist", key, user_id))
	}

	return user, true, nil
}

func UserFromID(id string) (models.User, bool, error) {
	var exists bool
	var user models.User

	user, exists = Users[id]

	if !exists {
		return user, false, nil
	}

	return user, true, nil


}

/**
 * Create a room
 * room_name        -> Name of this room. Request handlers should take care of name defaults.
 * open             -> Is this room open for any join?
 * owner_id         -> UUID of the owner of this room
 *
 * return           -> UUID of this room
 */
func NewRoom(room_name string, open bool, owner_id string) models.Room {
	var id string = uuid.NewV4().String()

	var room models.Room = models.Room{
		Name:        room_name,
		OwnerId:     owner_id,
		ID:          id,
		Open:        open,
		UserCount:   0,
		UserArray:   make([]string, 0, 65_536),
		AdminArray:  make([]string, 0, 65_536),
		InviteArray: make([]string, 0, 65_536),
	}

	Rooms[id] = room
	return room
}
