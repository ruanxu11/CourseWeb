package main

type Course struct {
	ID               int         `bson:"_id"` // 课程代码
	Name             string      // 课程名称
	College          string      // 开课学院
	Type             string      // 课程类型
	HoursPerWeek     int         // 周学时
	Term             string      // 学期
	Credit           float32     // 学分
	ClassRooms       []ClassRoom // 上课教室
	Introduction     string      // 课程简介
	TeachingSyllabus string      // 教学大纲
}

type ClassRoom struct {
	Time    string
	Positon string
}
