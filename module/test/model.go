package test

import (
	"encoding/json"
	"go_poc/core"
	"io"
)

// Test table model
type Test struct {
	Id   int64  `form:"-" json:"-" sqlParameterName:"id"`
	Name string `form:"name" json:"name" sqlParameterName:"name"`
}

// TestList is a shortcut to a list of Test
type TestList []*Test

// IsValid check if Test object is balid
func (t *Test) IsValid() *core.TYPoc {
	if t.Name == "" {
		return core.NewModelError("Test.IsValid", "name", "boo")
	}
	return nil
}

// ToJson serializes the bot patch to json.
func (t *Test) ToJson() []byte {
	data, err := json.Marshal(t)
	if err != nil {
		return nil
	}

	return data
}

// BotPatchFromJson deserializes a bot patch from json.
func TestFromJson(data io.Reader) *Test {
	decoder := json.NewDecoder(data)
	var testData Test
	err := decoder.Decode(&testData)
	if err != nil {
		return nil
	}

	return &testData
}

// ToJson serializes the bot patch to json.
func (tl *TestList) ToJson() []byte {
	data, err := json.Marshal(tl)
	if err != nil {
		return nil
	}

	return data
}

// BotPatchFromJson deserializes a bot patch from json.
func TestListFromJson(data io.Reader) *TestList {
	decoder := json.NewDecoder(data)
	var testList TestList
	err := decoder.Decode(&testList)
	if err != nil {
		return nil
	}

	return &testList
}
