package main

import (
	"errors"
	"koala"
	"log"
	"net/http"

	"labix.org/v2/mgo/bson"
)

func getClassIntroduction(id string) (string, error) {
	class, err := mgoFind("class", bson.M{"_id": id})
	if err != nil {
		return "", err
	}
	switch class["introduction"].(type) {
	case string:
		introduction := class["introduction"].(string)
		if introduction != "" {
			return introduction, nil
		}
	case nil:
	default:
		return "", errors.New("error type of class.introduction")
	}
	courseid := class["courseid"].(string)
	switch class["courseid"].(type) {
	case string:
	default:
		return "", errors.New("error type of courseid")
	}
	course, err := mgoFind("course", bson.M{"_id": courseid})
	switch course["introduction"].(type) {
	case string:
		introduction := course["introduction"].(string)
		return introduction, nil
	case nil:
	default:
		return "", errors.New("error type of Course.introduction")
	}
	return "", nil
}

func updateClassIntroduction(id string, introduction string) error {
	return mgoUpdate("class",
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"introduction": introduction}})
}

func classIntroductionHandlers() {
	koala.Get("/class/:id/introduction", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		id := p.ParamUrl["id"]
		introduction, err := getClassIntroduction(id)
		if err != nil {
			koala.NotFound(w)
			return
		}
		power := false
		log.Println("typesInClass:")
		if koala.ExistSession(r, "sessionID") {
			session := koala.GetSession(r, w, "sessionID")
			typesInClass := getTypeInClass(id, session.Values["collection"].(string), session.Values["id"].(string))
			log.Println(typesInClass)
			if typesInClass == "teacher" {
				power = true
			}
		}
		koala.Render(w, "class_introduction.html", map[string]interface{}{
			"title":        courseWeb,
			"id":           id,
			"introduction": introduction,
			"permission":   true,
			"power":        power,
		})
	})

	koala.Post("/class/:id/introduction", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		id := p.ParamUrl["id"]
		introduction := p.ParamPost["introduction"][0]
		power := false
		log.Println("typesInClass:")
		if koala.ExistSession(r, "sessionID") {
			session := koala.GetSession(r, w, "sessionID")
			typesInClass := getTypeInClass(id, session.Values["collection"].(string), session.Values["id"].(string))
			log.Println(typesInClass)
			if typesInClass == "teacher" {
				power = true
			}
		}
		if !power {
			return
		}
		err := updateClassIntroduction(id, introduction)
		if err != nil {
			log.Println(err)
			koala.Relocation(w, "/class/"+id+"/introduction", "更新课程介绍失败", "error")
		} else {
			koala.Relocation(w, "/class/"+id+"/introduction", "更新课程介绍成功", "success")
		}
	})
}
