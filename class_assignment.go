package main

import (
	"koala"
	"log"
	"net/http"
	"time"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Assignment struct {
	ID            string
	Time          string      // 作业发布时间
	Deadline      string      // 作业ddl
	Topic         string      // 作业
	MaxScore      string      // 作业分值
	Type          string      // 作业类型 选择，填空，问答题，上交文档
	Content       interface{} // 作业内容
	AttachPath    string
	AttachName    string
	AssignmentDos []AssignmentDo // 提交作业的名单
}

func addClassAssignment(id string, assignment *Assignment) error {
	assignment.Time = time.Now().Format("2006-01-02 15:04:05")
	assignment.ID = koala.HashString(time.Now().Format(time.UnixDate) + assignment.Topic + assignment.Deadline)
	return mgoUpdate("class",
		bson.M{"_id": id},
		bson.M{"$push": bson.M{"assignments": &assignment}})
}

func removeClassAssignmentlByID(id string, ID string) error {
	return mgoUpdate("class",
		bson.M{"_id": id},
		bson.M{"$pull": bson.M{"assignments": bson.M{"id": ID}}})
}

func getClassAssignment(id string, assignmentid string) (map[string]interface{}, error) {
	assignment := map[string]interface{}{}
	q := func(c *mgo.Collection) (*mgo.ChangeInfo, error) {
		pipe := c.Pipe([]bson.M{
			{
				"$unwind": "$assignments",
			},
			{
				"$project": bson.M{
					"id":                 "$assignments.id",
					"time":               "$assignments.time",
					"deadline":           "$assignments.deadline",
					"topic":              "$assignments.topic",
					"maxscore":           "$assignments.maxscore",
					"studentassignments": "$assignments.studentassignments",
					"type":               "$assignments.type",
					"content":            "$assignments.content",
					"attachpath":         "$assignments.attachpath",
					"attachname":         "$assignments.attachname",
				},
			},
			{
				"$group": bson.M{
					"_id": bson.M{
						"id":                 "$id",
						"time":               "$time",
						"deadline":           "$deadline",
						"topic":              "$topic",
						"maxscore":           "$maxscore",
						"studentassignments": "$studentassignments",
						"type":               "$type",
						"content":            "$content",
						"attachpath":         "$attachpath",
						"attachname":         "$attachname",
					},
				},
			},
			{
				"$match": bson.M{
					"_id.id": assignmentid,
				},
			},
		})
		iter := pipe.Iter()
		tag := bson.M{}
		iter.Next(&tag)
		assignment = tag["_id"].(bson.M)
		if err := iter.Close(); err != nil {
			return nil, err
		}
		return nil, nil
	}
	_, err := withCollection("class", q)
	return assignment, err
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
					"id":                 "$assignments.id",
					"time":               "$assignments.time",
					"deadline":           "$assignments.deadline",
					"topic":              "$assignments.topic",
					"maxscore":           "$assignments.maxscore",
					"studentassignments": "$assignments.studentassignments",
					"type":               "$assignments.type",
					"content":            "$assignments.content",
					"attachpath":         "$assignments.attachpath",
					"attachname":         "$assignments.attachname",
				},
			},
			{
				"$group": bson.M{
					"_id": bson.M{
						"id":                 "$id",
						"time":               "$time",
						"deadline":           "$deadline",
						"topic":              "$topic",
						"maxscore":           "$maxscore",
						"studentassignments": "$studentassignments",
						"type":               "$type",
						"content":            "$content",
						"attachpath":         "$attachpath",
						"attachname":         "$attachname",
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
		koala.Render(w, "class_assignment.html", map[string]interface{}{
			"title":       courseWeb,
			"id":          id,
			"assignments": assignments,
			"powers":      getPowers(r, id),
		})
	})

	koala.Post("/class/:id/assignment/add", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		id := p.ParamUrl["id"]
		powers := getPowers(r, id)
		if !powers["AssignmentAdd"] {
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
		AttachPath, filename, _, err := koala.SavePostFile(r, "file", "/assignment/"+id+"/")
		err = addClassAssignment(id, &Assignment{
			Deadline:   ddl,
			Topic:      content,
			MaxScore:   MaxScore,
			Type:       Type,
			AttachPath: AttachPath,
			AttachName: filename,
		})
		if err != nil {
			log.Println(err)
			koala.Relocation(w, "/class/"+id+"/assignment", "新增作业失败", "error")
			return
		}
		koala.Relocation(w, "/class/"+id+"/assignment", "新增作业成功", "success")
	})

	koala.Get("/class/:id/assignment/remove/:assid", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		id := p.ParamUrl["id"]
		assid := p.ParamUrl["assid"]
		powers := getPowers(r, id)
		if !powers["AssignmentRemove"] {
			koala.NotFound(w)
			return
		}
		err := removeClassAssignmentlByID(id, assid)
		if err != nil {
			log.Println(err)
			koala.Relocation(w, "/class/"+id+"/assignment", "删除作业失败", "error")
		} else {
			koala.Relocation(w, "/class/"+id+"/assignment", "删除作业成功", "success")
		}
	})
}
