package main

var courseWeb = "浙江大学课程网站系统"
var globalPowers = map[string][]string{
	"Forum":                  []string{"teacher", "teachingAssistant", "student", "otherTeacher", "otherTeachingAssistant", "otherStudent", "others"},
	"ForumPost":              []string{"teacher", "teachingAssistant", "student"},
	"ForumReply":             []string{"teacher", "teachingAssistant", "student"},
	"ForumPostRemove":        []string{"teacher", "teachingAssistant"},
	"MaterialAdd":            []string{"teacher"},
	"MaterialRemove":         []string{"teacher"},
	"MaterialDownload":       []string{"teacher", "teachingAssistant", "student"},
	"AssignmentAdd":          []string{"teacher", "teachingAssistant"},
	"AssignmentRemove":       []string{"teacher", "teachingAssistant"},
	"AssignmentUpdate":       []string{"teacher", "teachingAssistant"},
	"AssignmentView":         []string{"teacher", "teachingAssistant", "student"},
	"AssignmentCheck":        []string{"teacher", "teachingAssistant"},
	"AssignmentDo":           []string{"student"},
	"AnnouncementAdd":        []string{"teacher"},
	"AnnouncementRemove":     []string{"teacher"},
	"AnnouncementUpdate":     []string{"teacher"},
	"AnnouncementView":       []string{"teacher", "student"},
	"TeachingSyllabusUpdate": []string{"teacher"},
	"IntroductionUpdate":     []string{"teacher"},
	"StudentList":            []string{"teacher", "teachingAssistant", "student"},
	"MakeTeam":               []string{"teacher"},
	"PowersControl":          []string{"teacher"},
}

func main() {
	// addStudents()
	// go koala.HttpGet()
	// go koala.HttpPost()
	// go koala.HttpGetPost()
	// for i := 10000; i < 20000; i++ {
	// 	getTeachers(i)
	// }
	// for i := 21100000; i < 21200000; i++ {
	// 	getCourses(i) // 21120502
	// }
	// getCourses(21188020)
	// addTeacher("0026966", "林海", "男", "", "", "计算机科学与技术学院")
	// addTeacher("0023257", "邢卫", "男", "", "", "计算机科学与技术学院")
	// go addClasses()
	// updateCourseIntroduction("21191730", "这是一门叫软件需求的课")
	// updateClassIntroduction("8475187ad4e2c50fed755d2ed04e1ed2", "")
	// courses, err := searchCourses(map[string]interface{}{
	// 	"name": "软件需求工程",
	// })
	// if err != nil {
	// 	log.Println(err)
	// } else {
	// 	log.Println(courses)
	// }
	// announcement := &Announcement{
	// 	Title:   "最后2条公告",
	// 	Content: "最后2条公告内容",
	// }
	// err := addClassAnnouncement("8f48b0a04ed189183e224b20dbe156b0", announcement)
	// if err != nil {
	// 	log.Println(err)
	// }
	// err := updateClassAnnouncementByTime("37929210eaed36f72d39377b05dc8484", "2016-12-09 12:43:44", "没标题", "没内容")
	// if err != nil {
	// 	log.Println(err)
	// }
	// err := removeClassAnnouncementByTime("8f48b0a04ed189183e224b20dbe156b0", "2016-12-09 12:56:21")
	// if err != nil {
	// 	log.Println(err)
	// }
	// initStudents()
	// addStudentsInClass("8f48b0a04ed189183e224b20dbe156b0", &StudentInClass{
	// 	ID:   "3140102431",
	// 	Name: "徐亮",
	// })
	// initStudents()
	// err := updateClassPowerTable("c8503371c232707818bff7b17ab9975d", map[string]uint{
	// 	"ForumView": 0xff,
	// 	"ForumPost": 0xf0,
	// })
	// res, err := getClassPowerTable("14592062082a15a1ed3fe4925096996e")
	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println(res)
	RouteStart()
}
