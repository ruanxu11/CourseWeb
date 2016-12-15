package main

import (
	"koala"
	"log"
	"net/http"
	"time"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type AssignmentDo struct {
	ID           string `bson:"_id"`
	Time         string
	AssignmentID string      // 作业id
	StudentID    string      // 学生
	Student      string      // 学生姓名
	State        string      // 作业状态
	Content      interface{} // 提交的作业内容
	AttachPath   string      // 提交的附件
	AttachName   string      // 提交的附件
	Score        string      // 得分
	Checker      string      // 批改人
	Comment      string      // 评价
}

func updateClassAssignmentDo(assignmentDo *AssignmentDo) error {
	assignmentDo.Time = time.Now().Format("2006-01-02 15:04:05")
	assignmentDo.ID = koala.HashString(assignmentDo.AssignmentID + assignmentDo.StudentID)
	assignmentDo.State = "已提交"
	_, err := mgoFind("assignmentDo", bson.M{"_id": assignmentDo.ID})
	if err != nil {
		log.Println(err)
		return mgoInsert("assignmentDo", &assignmentDo)
	}
	return mgoUpdate("assignmentDo", bson.M{"_id": assignmentDo.ID}, &assignmentDo)
}

func getClassAssignmentDo(assignmentid string, studentid string) (map[string]interface{}, error) {
	return mgoFind("assignmentDo", bson.M{"assignmentid": assignmentid, "studentid": studentid})
}

func getClassAssignmentDos(assignmentid string) ([]map[string]interface{}, error) {
	return mgoFindAll("assignmentDo", bson.M{"assignmentid": assignmentid})
}

func checkClassAssignmentDos(checker string, assignmentid string, studentids []string, scores []string, comments []string) error {
	for i := 0; i < len(studentids); i++ {
		err := mgoUpdate("assignmentDo",
			bson.M{"assignmentid": assignmentid, "studentid": studentids[i]},
			bson.M{"$set": bson.M{"state": "已批改", "checker": checker, "score": scores[i], "comment": comments[i]}})
		if err != nil {
			return err
		}
	}
	return nil
}

func isTeamLeader(id string, sid string) bool {
	student := make(map[string]interface{})
	q := func(c *mgo.Collection) (*mgo.ChangeInfo, error) {
		pipe := c.Pipe([]bson.M{
			{
				"$unwind": "$students",
			},
			{
				"$project": bson.M{
					"id":         "$students.id",
					"teamleader": "$students.teamleader",
				},
			},
			{
				"$group": bson.M{
					"_id": bson.M{
						"id":         "$id",
						"teamleader": "$teamleader",
					},
				},
			},
			{
				"$match": bson.M{
					"_id.id": sid,
				},
			},
		})
		iter := pipe.Iter()
		tag := make(map[string]interface{})
		iter.Next(&tag)
		student = tag["_id"].(map[string]interface{})
		if err := iter.Close(); err != nil {
			return nil, err
		}
		return nil, nil
	}
	_, err := withCollection("class", q)
	if err != nil {
		return false
	}
	log.Println(student)
	return student["teamleader"].(bool)
}

func classAssignmentDoHandlers() {
	koala.Get("/class/:id/assignment/do/:assignmentid", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		id := p.ParamUrl["id"]
		powers := getPowers(r, id)
		if !powers["AssignmentDo"] {
			koala.NotFound(w)
			return
		}
		assignmentid := p.ParamUrl["assignmentid"]
		assignment, err := getClassAssignment(id, assignmentid)
		if err != nil {
			log.Println(err)
			koala.NotFound(w)
			return
		}
		session := koala.PeekSession(r, "sessionID")
		sid := session.Values["id"].(string)
		if assignment["type"].(string) == "小组作业" && !isTeamLeader(id, sid) {
			koala.Relocation(w, "/class/"+id+"/assignment", "组长才有权限提交作业", "warnning")
			return
		}
		assignmentDo, _ := getClassAssignmentDo(assignmentid, sid)
		koala.Render(w, "class_assignment_do.html", map[string]interface{}{
			"title":        courseWeb,
			"id":           id,
			"assignmentid": assignmentid,
			"assignment":   assignment,
			"assignmentDo": assignmentDo,
		})
	})

	koala.Post("/class/:id/assignment/do/:assignmentid", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		id := p.ParamUrl["id"]
		assignmentid := p.ParamUrl["assignmentid"]
		powers := getPowers(r, id)
		if !powers["AssignmentDo"] {
			koala.NotFound(w)
			return
		}
		content := ""
		r.ParseMultipartForm(32 << 20)
		if r.MultipartForm != nil {
			content = r.MultipartForm.Value["content"][0]
		}
		AttachPath, filename, _, err := koala.SavePostFile(r, "file", "/assignment/"+id+"/"+assignmentid+"/")
		session := koala.PeekSession(r, "sessionID")
		err = updateClassAssignmentDo(&AssignmentDo{
			AssignmentID: assignmentid,
			StudentID:    session.Values["id"].(string),
			Student:      session.Values["name"].(string),
			Content:      content,
			AttachPath:   AttachPath,
			AttachName:   filename,
		})
		if err != nil {
			log.Println(err)
			koala.Relocation(w, "/class/"+id+"/assignment/do/"+assignmentid, "提交作业失败", "error")
			return
		}
		koala.Relocation(w, "/class/"+id+"/assignment/do/"+assignmentid, "提交作业成功", "success")
	})

	koala.Get("/class/:id/assignment/check/:assignmentid", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		id := p.ParamUrl["id"]
		assignmentid := p.ParamUrl["assignmentid"]
		assignment, err := getClassAssignment(id, assignmentid)
		if err != nil {
			log.Println(err)
			koala.NotFound(w)
			return
		}
		assignmentDos, _ := getClassAssignmentDos(assignmentid)
		powers := getPowers(r, id)
		if !powers["AssignmentCheck"] {
			koala.NotFound(w)
			return
		}
		koala.Render(w, "class_assignment_check.html", map[string]interface{}{
			"title":         courseWeb,
			"id":            id,
			"assignmentid":  assignmentid,
			"assignment":    assignment,
			"assignmentDos": assignmentDos,
		})
	})

	koala.Post("/class/:id/assignment/check/:assignmentid", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		id := p.ParamUrl["id"]
		powers := getPowers(r, id)
		if !powers["AssignmentCheck"] {
			koala.NotFound(w)
			return
		}
		assignmentid := p.ParamUrl["assignmentid"]
		studentids := p.ParamPost["studentid"]
		scores := p.ParamPost["score"]
		comments := p.ParamPost["comment"]
		session := koala.PeekSession(r, "sessionID")
		checker := session.Values["name"].(string)
		err := checkClassAssignmentDos(checker, assignmentid, studentids, scores, comments)
		if err != nil {
			log.Println(err)
			koala.Relocation(w, "/class/"+id+"/assignment/check/"+assignmentid, "批改作业失败", "error")
			return
		}
		koala.Relocation(w, "/class/"+id+"/assignment/check/"+assignmentid, "批改作业成功", "success")
	})
}
