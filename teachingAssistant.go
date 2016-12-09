package main

import "labix.org/v2/mgo/bson"

type TeachingAssistant struct {
	ID                string `bson:"_id"`
	Password          string
	Name              string
	Sex               string
	Email             string
	Phone             string
	Introduction      string
	SecurityQuestions []SecurityQuestion
}

func showTeachingAssistantInfo(collection string, id string) (*map[string]interface{}, error) {
	person, err := mgoFind(collection,
		bson.M{"_id": id})
	personInfo := make(map[string]interface{})
	personInfo["id"] = person["_id"]
	personInfo["name"] = person["name"]
	personInfo["sex"] = person["sex"]
	personInfo["email"] = person["email"]
	personInfo["phone"] = person["phone"]
	personInfo["introduction"] = person["introduction"]
	return &personInfo, err
}
