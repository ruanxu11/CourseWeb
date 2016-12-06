package main

type Class struct {
	CourseID            int
	Course              string
	TeacherID           []int
	Teacher             []string
	TeachingAssistantID int
	TeachingAssistant   string
	Students            []StudentInClass // 每个学生的学号
	Assignments         []Assignment     // 作业
	Forum               []Post           // 讨论区
	Materials           []Material       // 课程资料
	Announcements       []Announcement   // 课程公告
}

type StudentInClass struct {
	ID   int
	Name []string
	Team int
}

type Announcement struct {
	Time    string // 公告发布时间
	Topic   string // 公告标题
	Content string // 公告内容
}

type Assignment struct {
	Time     string // 作业发布时间
	Deadline string // 作业ddl
	Topic    string // 作业
	Score    int    // 作业分值
	Type     string // 作业类型 选择，填空，大题，上交文档
	Content  string // 作业内容
}

type Material struct {
	Time              string // 资料发布时间
	BriefIntroduction string // 资料简介
	Type              string // 资料类型
	URL               string // 资料地址
}
