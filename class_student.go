package main

import "labix.org/v2/mgo/bson"

type StudentInClass struct {
	ID                 string
	Name               string
	Team               string
	TeamLeader         bool
	StudentAssignments []StudentAssignment // 作业情况
}

type StudentAssignment struct {
	ID      string      // 作业id
	Content interface{} // 提交的作业内容
	Attach  string      // 提交的附件
	Score   string      // 得分
	Comment string      // 评价
}

func addStudentsInClass(id string, student *StudentInClass) error {
	return mgoUpdate("class",
		bson.M{"_id": id},
		bson.M{"$push": bson.M{"students": &student}})
}
