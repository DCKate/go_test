package kkutil

import (
	"fmt"
	"testing"
)

type User struct {
	Who        string `kk:"name"`
	ID         int    `kk:"number"`
	Descrption string `kk:"text"`
}

func TestSerialize(t *testing.T) {
	mm := map[string]interface{}{
		"name":   "John",
		"number": 11,
		"text":   "hello world",
	}
	tmp := serialize(mm, User{}).Interface()
	uu := tmp.(*User)
	fmt.Printf("[%d]%s : %s\n", uu.ID, uu.Who, uu.Descrption)
}
