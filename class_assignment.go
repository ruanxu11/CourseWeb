package main

import (
	"io"
	"koala"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Assignment struct {
	ID         string
	Time       string      // 作业发布时间
	Deadline   string      // 作业ddl
	Topic      string      // 作业
	MaxScore   string      // 作业分值
	Type       string      // 作业类型 选择，填空，问答题，上交文档
	Content    interface{} // 作业内容
	AttachPath string
}

func addClassAssignment(id string, assignment *Assignment) error {
	assignment.Time = time.Now().Format("2006-01-02 15:04:05")
	assignment.ID = koala.HashString(time.Now().Format(time.UnixDate) + assignment.Topic + assignment.Deadline)
	return mgoUpdate("class",
		bson.M{"_id": id},
		bson.M{"$push": bson.M{"assignments": &assignment}})
}

func getClassAssignments(id string) ([]map[string]interface{}, error) {
	assignments := []map[string]interface{}{}
	q := func(c *mgo.Collection) (*mgo.ChangeInfo, error) {
		pipe := c.Pipe([]bson.M{
			{
				"$unwind": "$assignments",
			},
			{
				"$project": bson.M{
					"time":               "$assignments.time",
					"deadline":           "$assignments.deadline",
					"topic":              "$assignments.topic",
					"maxscore":           "$assignments.maxscore",
					"studentassignments": "$assignments.studentassignments",
					"type":               "$assignments.type",
					"content":            "$assignments.content",
					"attachpath":         "$assignments.attachpath",
				},
			},
			{
				"$group": bson.M{
					"_id": bson.M{
						"time":               "$time",
						"deadline":           "$deadline",
						"topic":              "$topic",
						"maxscore":           "$maxscore",
						"studentassignments": "$studentassignments",
						"type":               "$type",
						"content":            "$content",
						"attachpath":         "$attachpath",
					},
				},
			},
			{
				"$sort": bson.M{
					"_id.time": -1,
				},
			},
		})
		iter := pipe.Iter()
		tag := bson.M{}
		for iter.Next(&tag) {
			assignments = append(assignments, tag["_id"].(bson.M))
		}
		if err := iter.Close(); err != nil {
			return nil, err
		}
		return nil, nil
	}
	_, err := withCollection("class", q)
	return assignments, err
}

func classAssignmentHandlers() {
	koala.Get("/class/:id/assignment", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		id := p.ParamUrl["id"]
		assignments, err := getClassAssignments(id)
		if err != nil {
			koala.NotFound(w)
			return
		}
		addpower := false
		viewpower := false
		if koala.ExistSession(r, "sessionID") {
			session := koala.GetSession(r, w, "sessionID")
			typesInClass := getTypeInClass(id, session.Values["collection"].(string), session.Values["id"].(string))
			log.Println(typesInClass)
			if typesInClass == "teacher" {
				addpower = true
			}
			viewpower = true
		}
		koala.Render(w, "class_assignment.html", map[string]interface{}{
			"title":       courseWeb,
			"id":          id,
			"assignments": assignments,
			"addpower":    addpower,
			"viewpower":   viewpower,
		})
	})

	koala.Post("/class/:id/assignment/add", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		id := p.ParamUrl["id"]
		power := false
		if koala.ExistSession(r, "sessionID") {
			session := koala.GetSession(r, w, "sessionID")
			typesInClass := getTypeInClass(id, session.Values["collection"].(string), session.Values["id"].(string))
			if typesInClass == "teacher" {
				power = true
			}
		}
		if !power {
			koala.NotFound(w)
			return
		}
		content := ""
		ddl := ""
		Type := ""
		MaxScore := ""
		r.ParseMultipartForm(32 << 20)
		if r.MultipartForm != nil {
			content = r.MultipartForm.Value["content"][0]
			ddl = r.MultipartForm.Value["ddl"][0]
			Type = r.MultipartForm.Value["Type"][0]
			MaxScore = r.MultipartForm.Value["maxScore"][0]
		}
		file, handle, err := r.FormFile("file")
		AttachPath := ""
		if err == nil {

			filename := handle.Filename
			suffix := path.Ext(filename)
			log.Println(suffix)
			courseid, err := getCourseID(id)
			if err != nil {
				log.Println(err)
				koala.Relocation(w, "/class/"+id+"/assignment", "不存在这门课程", "error")
				return
			}
			AttachPath = "/assignment/" + courseid + "/" + filename
			filepath := "./static/upload/assignment/" + courseid + "/" + filename
			os.MkdirAll(path.Dir(filepath), 0777)
			f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				log.Println(err)
				koala.Relocation(w, "/class/"+id+"/assignment", "新增作业失败", "error")
				return
			}
			_, err = io.Copy(f, file)
			if err != nil {
				log.Println(err)
				koala.Relocation(w, "/class/"+id+"/assignment", "新增作业失败", "error")
				return
			}
			defer f.Close()
			defer file.Close()
		}
		err = addClassAssignment(id, &Assignment{
			Deadline:   ddl,
			Topic:      content,
			MaxScore:   MaxScore,
			Type:       Type,
			AttachPath: AttachPath,
		})
		if err != nil {
			log.Println(err)
			koala.Relocation(w, "/class/"+id+"/assignment", "新增作业失败", "error")
			return
		}
		koala.Relocation(w, "/class/"+id+"/assignment", "新增作业成功", "success")
	})
}
