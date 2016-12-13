package main

import (
	"errors"
	"koala"
	"log"
	"net/http"

	"labix.org/v2/mgo/bson"
)

func getClassTeachingSyllabus(id string) (string, error) {
	class, err := mgoFind("class", bson.M{"_id": id})
	if err != nil {
		return "", err
	}
	log.Println("class.teachingsyllabus")
	switch class["teachingsyllabus"].(type) {
	case string:
		teachingsyllabus := class["teachingsyllabus"].(string)
		if teachingsyllabus != "" {
			return teachingsyllabus, nil
		}
	case nil:
	default:
		return "", errors.New("error type of class.teachingsyllabus")
	}
	log.Println("courseid")
	courseid := class["courseid"].(string)
	switch class["courseid"].(type) {
	case string:
	default:
		return "", errors.New("error type of courseid")
	}
	course, err := mgoFind("course", bson.M{"_id": courseid})
	log.Println("course.teachingsyllabus")
	log.Println(courseid)
	switch course["teachingsyllabus"].(type) {
	case string:
		teachingsyllabus := course["teachingsyllabus"].(string)
		return teachingsyllabus, nil
	case nil:
	default:
		return "", errors.New("error type of Course.teachingsyllabus")
	}
	return "", nil
}

func updateClassTeachingSyllabus(id string, teachingsyllabus string) error {
	return mgoUpdate("class",
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"teachingsyllabus": teachingsyllabus}})
}

func classTeachingSyllabusHandlers() {
	koala.Get("/class/:id/teachingsyllabus", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		id := p.ParamUrl["id"]
		teachingsyllabus, err := getClassTeachingSyllabus(id)
		if err != nil {
			koala.NotFound(w)
			return
		}
		koala.Render(w, "class_teachingsyllabus.html", map[string]interface{}{
			"title":            courseWeb,
			"id":               id,
			"teachingsyllabus": teachingsyllabus,
			"powers":           getPowersInClass(r, id),
		})
	})

	koala.Post("/class/:id/teachingsyllabus", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		id := p.ParamUrl["id"]
		powers := getPowersInClass(r, id)
		if !powers["TeachingSyllabusUpdate"] {
			koala.NotFound(w)
			return
		}
		teachingsyllabus := p.ParamPost["teachingsyllabus"][0]
		err := updateClassTeachingSyllabus(id, teachingsyllabus)
		if err != nil {
			log.Println(err)
			koala.Relocation(w, "/class/"+id+"/teachingsyllabus", "更新课程大纲失败", "error")
		} else {
			koala.Relocation(w, "/class/"+id+"/teachingsyllabus", "更新课程大纲成功", "success")
		}
	})
}
