package user

import (
	"encoding/json"
	"go_poc/core"
	"io"
)

// User table model
type User struct {
	Id    int64  `form:"-" json:"-" sqlParameterName:"id"`
	Name  string `form:"name" json:"name" sqlParameterName:"name"`
	Email string `form:"email" json:"email" sqlParameterName:"email"`
}

// UserList is a shortcut to a list of User
type UserList []*User

// IsValid check if User object is valid
func (u *User) IsValid() *core.TYPoc {
	if u.Name == "" {
		return core.NewModelError("User.IsValid", "name", "name required")
	}
	if u.Email == "" {
		return core.NewModelError("User.IsValid", "email", "email required")
	}
	return nil
}

// ToJson serializes the bot patch to json.
func (u *User) ToJson() []byte {
	data, err := json.Marshal(u)
	if err != nil {
		return nil
	}

	return data
}

// BotPatchFromJson deserializes a bot patch from json.
func UserFromJson(data io.Reader) *User {
	decoder := json.NewDecoder(data)
	var userData User
	err := decoder.Decode(&userData)
	if err != nil {
		return nil
	}

	return &userData
}

// ToJson serializes the bot patch to json.
func (ul *UserList) ToJson() []byte {
	data, err := json.Marshal(ul)
	if err != nil {
		return nil
	}

	return data
}

// BotPatchFromJson deserializes a bot patch from json.
func UserListFromJson(data io.Reader) *UserList {
	decoder := json.NewDecoder(data)
	var userList UserList
	err := decoder.Decode(&userList)
	if err != nil {
		return nil
	}

	return &userList
}
