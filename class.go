package main

import (
	"fmt"
	"io/ioutil"
	"koala"
	"net/http"
	"net/url"
	"strings"
	"time"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Class struct {
	ID                  string `bson:"_id"`
	CourseID            string
	Course              string
	Year                string
	Semester            string
	Introduction        string      // 课程介绍 如果为空则用course的简介
	TeachingSyllabus    string      // 教学大纲 如果为空则用course的大纲
	ClassRooms          []ClassRoom // 上课教室
	Teachers            []TeacherInClass
	TeachingAssistantID string
	TeachingAssistant   string
	Students            []StudentInClass // 每个学生的学号
	Assignments         []Assignment     // 作业
	Forum               []Post           // 讨论区
	Materials           []Material       // 课程资料
	Announcements       []Announcement   // 课程公告
	Powers              PowersInClass
}

type ClassRoom struct {
	Time     string
	Position string
}

type TeacherInClass struct {
	ID   string
	Name string
}

func addClasses() {
	v := url.Values{}
	v.Set("CourseID", "21191730")
	v.Set("Course", "软件需求工程")
	v.Set("Year", "2016-2017")
	v.Set("Semester", "秋冬")
	v.Set("ClassRoomNum", "2")
	v.Set("ClassRoomTime", "周六第11,12节{双周}")
	v.Set("ClassRoomPositon", "玉泉曹光彪西-503")
	v.Add("ClassRoomTime", "周一第6,7,8节")
	v.Add("ClassRoomPositon", "玉泉教4-413(多)")
	v.Set("TeacherNum", "2")
	v.Set("TeacherID", "0023257")
	v.Set("TeacherName", "邢卫")
	v.Add("TeacherID", "0026966")
	v.Add("TeacherName", "林海")
	body := ioutil.NopCloser(strings.NewReader(v.Encode())) //把form数据编下码
	client := &http.Client{}
	r, _ := http.NewRequest("POST", "http://localhost:2333/addClass", body)

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded") //这个一定要加，不加form的值post不过去，被坑了两小时
	fmt.Printf("%+v\n\n\n", r)                                        //看下发送的结构

	resp, err := client.Do(r) //发送
	defer resp.Body.Close()   //一定要关闭resp.Body
	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(data), err)
}

func addClass(class *Class) error {
	class.ID = koala.HashString(time.Now().Format(time.UnixDate))
	return mgoInsert("class", &class)
}

func removeClasss(selector map[string]interface{}) (*mgo.ChangeInfo, error) {
	return mgoRemoveAll("class", selector)
}

func removeClassByID(id string) (*mgo.ChangeInfo, error) {
	return mgoRemove("class", bson.M{"_id": id})
}

func updateClasss(selector map[string]interface{}, update map[string]interface{}) (*mgo.ChangeInfo, error) {
	return mgoUpdateAll("class", selector, update)
}

func updateClassByID(id string, update map[string]interface{}) error {
	return mgoUpdate("class", bson.M{"_id": id}, update)
}

func searchClasss(selector map[string]interface{}) ([]map[string]interface{}, error) {
	return mgoFindAll("class", selector)
}

func getAllClasses() ([]map[string]interface{}, error) {
	return mgoFindAll("class", nil)
}

func getDistinctCourse() ([]map[string]interface{}, error) {
	courseid, err := mgoFindDistinct("class", nil, "courseid")
	if err != nil {
		return nil, err
	}
	return mgoFindAll("course", bson.M{"_id": bson.M{"$in": courseid}})
}

func removeClass(id string) (info *mgo.ChangeInfo, err error) {
	return mgoRemove("class", bson.M{"_id": id})
}

func getClass(id string) (map[string]interface{}, error) {
	return mgoFind("class", bson.M{"_id": id})
}

func getClasses() ([]map[string]interface{}, error) {
	return mgoFindAll("class", nil)
}

func searchClassesSelect(selector map[string]interface{}) ([]map[string]interface{}, error) {
	return mgoSearchSelect("class", selector)
}
func getClassesByPage(page int) ([]map[string]interface{}, error) {
	return mgoFindByPage("class", page)
}

func getClassByCourseID(courseid string) ([]map[string]interface{}, error) {
	return mgoFindAll("class", bson.M{"courseid": courseid})
}

func getCourseID(id string) (string, error) {
	class, err := mgoFindSelect("class", bson.M{"_id": id}, bson.M{"_id": 0, "courseid": 1})
	return class["courseid"].(string), err
}

func classHandlers() {
	classIntroductionHandlers()
	classTeachingSyllabusHandlers()
	classAnnouncementHandlers()
	classMaterialHandlers()
	classForumHandlers()
	classAssignmentHandlers()

	koala.Get("/class/:id", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		id := p.ParamUrl["id"]
		class, err := getClass(id)
		if err != nil {
			koala.NotFound(w)
			return
		}
		koala.Render(w, "class.html", map[string]interface{}{
			"title": courseWeb,
			"class": class,
		})
	})
}
