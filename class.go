package main

import (
	"fmt"
	"io/ioutil"
	"koala"
	"log"
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

func getClass(id string) (map[string]interface{}, error) {
	return mgoFind("class", bson.M{"_id": id})
}

func getClasses() ([]map[string]interface{}, error) {
	return mgoFindAll("class", nil)
}

func getClassesByID(idList []interface{}) ([]map[string]interface{}, error) {
	return mgoFindAll("class", bson.M{"_id": bson.M{"$in": idList}})
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

func getClassName(id string) (string, error) {
	class, err := mgoFindSelect("class", bson.M{"_id": id}, bson.M{"_id": 0, "course": 1})
	return class["course"].(string), err
}

func getClassTA(id string) (interface{}, error) {
	class, err := mgoFindSelect("class", bson.M{"_id": id}, bson.M{"_id": 0, "teachingassistantid": 1})
	return class["teachingassistantid"], err
}

func addClass(class *Class) error {
	class.ID = koala.HashString(time.Now().Format(time.UnixDate))
	for i := 0; i < len(class.Teachers); i++ {
		err := mgoUpdate("teacher",
			bson.M{"_id": class.Teachers[i].ID},
			bson.M{"$push": bson.M{"classes": class.ID}})
		if err != nil {
			return err
		}
	}
	return mgoInsert("class", &class)
}

func removeClass(id string) (info *mgo.ChangeInfo, err error) {
	class, err := getClass(id)
	if err != nil {
		return nil, err
	}
	switch class["teachers"].(type) {
	case []TeacherInClass:
		teachers := class["teachers"].([]TeacherInClass)
		for i := 0; i < len(teachers); i++ {
			mgoUpdate("teacher",
				bson.M{"_id": teachers[i].ID},
				bson.M{"$pull": bson.M{"classes": id}})
		}
	default:
	}
	switch class["students"].(type) {
	case []StudentInClass:
		students := class["teachers"].([]StudentInClass)
		for i := 0; i < len(students); i++ {
			mgoUpdate("student",
				bson.M{"_id": students[i].ID},
				bson.M{"$pull": bson.M{"classes": id}})
		}
	default:
	}
	switch class["teachingassistantid"].(type) {
	case string:
		teachingassistantid := class["teachingassistantid"]
		mgoUpdate("teachingAssistant",
			bson.M{"_id": teachingassistantid},
			bson.M{"$pull": bson.M{"classes": id}})
	default:
	}
	return mgoRemove("class", bson.M{"_id": id})
}

func updateClassTeachingAssistant(classid string, ID string, Name string) error {
	oldTAID, err := getClassTA(classid)
	if err == nil && oldTAID != nil {
		mgoUpdate("teachingAssistant",
			bson.M{"_id": oldTAID},
			bson.M{"$pull": bson.M{"classes": classid}})
	}
	mgoUpdate("teachingAssistant",
		bson.M{"_id": ID},
		bson.M{"$push": bson.M{"classes": classid}})
	return mgoUpdate("class",
		bson.M{"_id": classid},
		bson.M{"$set": bson.M{"teachingassistantid": ID, "teachingassistant": Name}})
}

func classHandlers() {
	classIntroductionHandlers()
	classTeachingSyllabusHandlers()
	classAnnouncementHandlers()
	classMaterialHandlers()
	classForumHandlers()
	classAssignmentHandlers()
	classStudentHandlers()
	classAssignmentDoHandlers()
	koala.Get("/class/:id", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		id := p.ParamUrl["id"]
		class, err := getClass(id)
		if err != nil {
			koala.NotFound(w)
			return
		}
		koala.Render(w, "class.html", map[string]interface{}{
			"title":  courseWeb,
			"class":  class,
			"admin":  admincheck(w, r),
			"powers": getPowersInClass(r, id),
		})
	})

	koala.Post("/class/:id/teachingAssistant/update", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		if !admincheck(w, r) {
			koala.NotFound(w)
			return
		}
		id := p.ParamUrl["id"]
		ID := p.ParamPost["ID"][0]
		Name := p.ParamPost["Name"][0]
		err := updateClassTeachingAssistant(id, ID, Name)
		if err != nil {
			log.Println(err)
			koala.Relocation(w, "/class/"+id, "修改助教失败", "error")
		} else {
			koala.Relocation(w, "/class/"+id, "修改助教成功", "success")
		}
	})
}
