package main

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

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

func addTeachingAssistant(teachingAssistant *TeachingAssistant) error {
	return mgoInsert("teachingAssistant", &teachingAssistant)
}

func removeTeachingAssistants(selector map[string]interface{}) (*mgo.ChangeInfo, error) {
	return mgoRemoveAll("teachingAssistant", selector)
}

func removeTeachingAssistantByID(id string) (*mgo.ChangeInfo, error) {
	return mgoRemove("teachingAssistant", bson.M{"_id": id})
}

func updateTeachingAssistants(selector map[string]interface{}, update map[string]interface{}) (*mgo.ChangeInfo, error) {
	return mgoUpdateAll("teachingAssistant", selector, update)
}

func updateTeachingAssistantByID(id string, update map[string]interface{}) error {
	return mgoUpdate("teachingAssistant", bson.M{"_id": id}, update)
}

func searchTeachingAssistants(selector map[string]interface{}) ([]map[string]interface{}, error) {
	return mgoFindAll("teachingAssistant", selector)
}

func getTeachingAssistantInfo(collection string, id string) (*map[string]interface{}, error) {
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
