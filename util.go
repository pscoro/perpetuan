package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"

	//"fmt"
	"log"
	"unicode"
)
import "github.com/badoux/checkmail"

func UsernameValid(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) || !unicode.IsDigit(r) || string(r) != "-" || string(r) != "_" {
			return false
		}
	}
	return true
}

func EmailValid(s string) bool {
	err := checkmail.ValidateFormat(s)
	if err != nil {
		return true
	}
	return false
}

func EmailUnique(s string) bool {
	res := connect.userCollection.FindOne(context.Background(), bson.M{"Email": s})
	if res.Err() != nil {
		return true
	}
	return false
}

func UsernameUnique(s string) bool {
	res := connect.userCollection.FindOne(context.Background(), bson.M{"Username": s})
	if res.Err() != nil {
		return true
	}
	return false
}

func PasswordValid(s string) bool {
	return true
}

func SettingValid(s string) bool {
	settingList,ok := updateSettingList([]setting{})

	if !ok {
		log.Fatal("Can not update setting lists")
		return false
	}

	for i := range settingList {
		if s == settingList[i].key {
			return true
		}
	}
	return false
}
