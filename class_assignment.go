package main

import (
	"koala"
	"log"
	"net/http"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Assignment struct {
	Time     string // 作业发布时间
	Deadline string // 作业ddl
	Topic    string // 作业
	Score    int    // 作业分值
	Type     string // 作业类型 选择，填空，大题，上交文档
	Content  string // 作业内容
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
					"time":    "$assignments.time",
					"title":   "$assignments.title",
					"content": "$assignments.content",
				},
			},
			{
				"$group": bson.M{
					"_id": bson.M{
						"time":    "$time",
						"title":   "$title",
						"content": "$content",
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
		koala.Render(w, "class_assignment.html", map[string]interface{}{
			"title":       courseWeb,
			"id":          id,
			"assignments": assignments,
			"permission":  true,
			"power":       power,
		})
	})
}
