package main

import (
	"log"
	"strconv"

	"labix.org/v2/mgo/bson"
)

func loginCheck(collection string, sid string, password string) (bool, string) {
	id, _ := strconv.Atoi(sid)
	person, err := mgoFind(collection, bson.M{
		"_id":      id,
		"password": password,
	}, 0, 0)
	log.Println(person)
	if err != nil {
		return false, ""
	}
	return true, person["name"].(string)
}
