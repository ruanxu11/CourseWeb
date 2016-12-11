package main

import (
	"log"

	"labix.org/v2/mgo/bson"
)

func getTypeInClass(classid string, collection string, id string) string {
	var err error
	if collection == "student" {
		_, err = mgoFind("class", bson.M{"_id": classid, "students.id": id})
		if err != nil {
			log.Println(err)
			return "others"
		}
		return "student"
	} else if collection == "teacher" {
		_, err = mgoFind("class", bson.M{"_id": classid, "teachers.id": id})
		if err != nil {
			log.Println(err)
			return "teacher"
		}
		return "teacher"
	} else if collection == "teachingassistant" {
		_, err = mgoFind("class", bson.M{"_id": classid, "teachingassistantid": id})
		if err != nil {
			log.Println(err)
			return "others"
		}
		return "teachingassistant"
	}
	return "others"
}
