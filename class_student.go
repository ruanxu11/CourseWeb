package main

import (
	"errors"
	"koala"
	"log"
	"net/http"

	"labix.org/v2/mgo/bson"
)

type StudentInClass struct {
	ID            string
	Name          string
	Team          string
	TeamLeader    bool
	AssignmentDos []AssignmentDo // 作业情况
}

func addStudentsInClass(id string, student *StudentInClass) error {
	err := mgoUpdate("student",
		bson.M{"_id": student.ID},
		bson.M{"$push": bson.M{"classes": id}})
	if err != nil {
		return err
	}
	return mgoUpdate("class",
		bson.M{"_id": id},
		bson.M{"$push": bson.M{"students": &student}})
}

func removeStudentInClassByID(classid string, ID string) error {
	err := mgoUpdate("class",
		bson.M{"_id": classid},
		bson.M{"$pull": bson.M{"students": bson.M{"id": ID}}})
	if err != nil {
		return err
	}
	return mgoUpdate("student",
		bson.M{"_id": ID},
		bson.M{"$pull": bson.M{"classes": classid}})
}

func makeTeam(classid string, studentID string, team string, teamleader string) error {
	TeamLeader := false
	if teamleader == "true" {
		TeamLeader = true
	} else if teamleader == "false" {
		TeamLeader = false
	} else {
		return errors.New("未知的teamleader类型")
	}
	return mgoUpdate("class",
		bson.M{"_id": classid, "students.id": studentID},
		bson.M{"$set": bson.M{
			"students.$.team":       team,
			"students.$.teamleader": TeamLeader,
		}})
}

func classStudentHandlers() {
	koala.Post("/class/:id/student/add", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		if !admincheck(w, r) {
			koala.NotFound(w)
			return
		}
		id := p.ParamUrl["id"]
		ID := p.ParamPost["ID"][0]
		Name := p.ParamPost["Name"][0]
		err := addStudentsInClass(id, &StudentInClass{
			ID:   ID,
			Name: Name,
		})
		if err != nil {
			log.Println(err)
			koala.Relocation(w, "/class/"+id, "添加学生失败", "error")
		} else {
			koala.Relocation(w, "/class/"+id, "添加学生成功", "success")
		}
	})

	koala.Post("/class/:id/team/make", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		id := p.ParamUrl["id"]
		powers := getPowers(r, id)
		if !powers["MakeTeam"] {
			koala.NotFound(w)
			return
		}
		ID := p.ParamPost["ID"]
		Team := p.ParamPost["Team"]
		for i := 0; i < len(ID); i++ {
			TeamLeader := p.ParamPost["TeamLeader"+ID[i]][0]
			err := makeTeam(id, ID[i], Team[i], TeamLeader)
			if err != nil {
				log.Println(err)
				koala.Relocation(w, "/class/"+id, "组队失败", "error")
				return
			}
		}
		koala.Relocation(w, "/class/"+id, "组队成功", "success")
	})

	koala.Get("/class/:id/student/remove/:ID", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		if !admincheck(w, r) {
			koala.NotFound(w)
			return
		}
		classid := p.ParamUrl["id"]
		ID := p.ParamUrl["ID"]
		err := removeStudentInClassByID(classid, ID)
		if err != nil {
			log.Println(err)
			koala.Relocation(w, "/class/"+classid, "删除学生失败", "error")
		} else {
			koala.Relocation(w, "/class/"+classid, "删除学生成功", "success")
		}
	})
}
