package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"koala"
	"log"
	"net/http"
	"net/url"
	"strconv"
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
	Introduction        string      // 课程简介 如果为空则用course的简介
	TeachingSyllabus    string      // 教学大纲 如果为空则用course的大纲
	ClassRooms          []ClassRoom // 上课教室
	Teachers            []TeacherInClass
	TeachingAssistantID int
	TeachingAssistant   string
	Students            []StudentInClass // 每个学生的学号
	Assignments         []Assignment     // 作业
	Forum               []Post           // 讨论区
	Materials           []Material       // 课程资料
	Announcements       []Announcement   // 课程公告
}

type ClassRoom struct {
	Time    string
	Positon string
}

type StudentInClass struct {
	ID   string
	Name []string
	Team int
}

type TeacherInClass struct {
	ID   string
	Name string
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
	req, _ := http.NewRequest("POST", "http://localhost:2333/addClass", body)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded") //这个一定要加，不加form的值post不过去，被坑了两小时
	fmt.Printf("%+v\n\n\n", req)                                        //看下发送的结构

	resp, err := client.Do(req) //发送
	defer resp.Body.Close()     //一定要关闭resp.Body
	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(data), err)
}

func addClass(CourseID string, Course string, Year string, Semester string, ClassRooms []ClassRoom, Teachers []TeacherInClass) error {
	class := &Class{
		ID:         koala.HashString(time.Now().Format(time.UnixDate)),
		CourseID:   CourseID,
		Course:     Course,
		Year:       Year,
		Semester:   Semester,
		ClassRooms: ClassRooms,
		Teachers:   Teachers,
	}
	return mgoInsert("class", &class)
}

func showAllClasses() ([]map[string]interface{}, error) {
	return mgoFindAll("class", nil)
}

func removeClass(id string) (info *mgo.ChangeInfo, err error) {
	return mgoRemove("class", bson.M{"_id": id})
}

func showClass(id string) (map[string]interface{}, error) {
	return mgoFind("class", bson.M{"_id": id})
}

func showClassIntroduction(id string) (string, error) {
	class, err := mgoFind("class", bson.M{"_id": id})
	if err != nil {
		return "", err
	}
	log.Println("class.introduction")
	switch class["introduction"].(type) {
	case string:
		introduction := class["introduction"].(string)
		if introduction != "" {
			return introduction, nil
		}
	case nil:
	default:
		return "", errors.New("error type of class.introduction")
	}
	log.Println("courseid")
	courseid := class["courseid"].(string)
	switch class["courseid"].(type) {
	case string:
	default:
		return "", errors.New("error type of courseid")
	}
	course, err := mgoFind("course", bson.M{"_id": courseid})
	log.Println("course.introduction")
	log.Println(courseid)
	switch course["introduction"].(type) {
	case string:
		introduction := course["introduction"].(string)
		return introduction, nil
	case nil:
	default:
		return "", errors.New("error type of Course.introduction")
	}
	return "", nil
}

func updateClassIntroduction(id string, introduction string) error {
	return mgoUpdate("class",
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"introduction": introduction}})
}

func showClassTeachingSyllabus(id string) (string, error) {
	class, err := mgoFind("class", bson.M{"_id": id})
	if err != nil {
		return "", err
	}
	log.Println("class.teachingsyllabus")
	switch class["teachingsyllabus"].(type) {
	case string:
		teachingsyllabus := class["teachingsyllabus"].(string)
		if teachingsyllabus != "" {
			return teachingsyllabus, nil
		}
	case nil:
	default:
		return "", errors.New("error type of class.teachingsyllabus")
	}
	log.Println("courseid")
	courseid := class["courseid"].(string)
	switch class["courseid"].(type) {
	case string:
	default:
		return "", errors.New("error type of courseid")
	}
	course, err := mgoFind("course", bson.M{"_id": courseid})
	log.Println("course.teachingsyllabus")
	log.Println(courseid)
	switch course["teachingsyllabus"].(type) {
	case string:
		teachingsyllabus := course["teachingsyllabus"].(string)
		return teachingsyllabus, nil
	case nil:
	default:
		return "", errors.New("error type of Course.teachingsyllabus")
	}
	return "", nil
}

func updateClassTeachingSyllabus(id string, teachingsyllabus string) error {
	return mgoUpdate("class",
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"teachingsyllabus": teachingsyllabus}})
}

func classHandlers() {
	koala.Handle("/add/class", func(p *koala.Params, w http.ResponseWriter, req *http.Request) {
		// if !koala.ExistSession(req, "sessionID") {
		// 	w.Write([]byte("请先登录"))
		// } else {
		CourseID := p.Param["CourseID"][0]
		Course := p.Param["Course"][0]
		Year := p.Param["Year"][0]
		Semester := p.Param["Semester"][0]
		sClassRoomNum := p.Param["ClassRoomNum"][0]
		ClassRoomNum, _ := strconv.Atoi(sClassRoomNum)
		ClassRooms := make([]ClassRoom, ClassRoomNum)
		for i := 0; i < ClassRoomNum; i++ {
			ClassRooms[i].Time = p.Param["ClassRoomTime"][i]
			ClassRooms[i].Positon = p.Param["ClassRoomPositon"][i]
		}
		sTeacherNum := p.Param["TeacherNum"][0]
		TeacherNum, _ := strconv.Atoi(sTeacherNum)
		Teachers := make([]TeacherInClass, TeacherNum)
		for i := 0; i < TeacherNum; i++ {
			Teachers[i].ID = p.Param["TeacherID"][i]
			Teachers[i].Name = p.Param["TeacherName"][i]
		}
		err := addClass(CourseID, Course, Year, Semester, ClassRooms, Teachers)
		if err != nil {
			w.Write([]byte("添加教学班失败\n" + err.Error()))
		} else {
			w.Write([]byte("添加教学班成功\n"))
			classes, err := showAllClasses()
			if err != nil {
				w.Write([]byte("查看教学班失败\n" + err.Error()))
			} else {
				json, _ := json.Marshal(classes)
				w.Write([]byte("教学班\n"))
				w.Write([]byte(json))
			}
		}
		// }
	})

	koala.Handle("/classes", func(p *koala.Params, w http.ResponseWriter, req *http.Request) {
		// if !koala.ExistSession(req, "sessionID") {
		// 	w.Write([]byte("请先登录"))
		// } else {
		classes, err := showAllClasses()
		if err != nil {
			w.Write([]byte("查看教学班失败\n" + err.Error()))
		} else {
			json, _ := json.Marshal(classes)
			w.Write([]byte("教学班\n"))
			w.Write([]byte(json))
		}
		// }
	})

	koala.Handle("/remove/class", func(p *koala.Params, w http.ResponseWriter, req *http.Request) {
		// if !koala.ExistSession(req, "sessionID") {
		// 	w.Write([]byte("请先登录"))
		// } else {
		id := p.Param["id"][0]
		_, err := removeClass(id)
		if err != nil {
			w.Write([]byte("删除教学班失败\n" + err.Error()))
		} else {
			w.Write([]byte("删除教学班成功\n"))
		}
		// }
	})

	koala.Handle("/class/:id", func(p *koala.Params, w http.ResponseWriter, req *http.Request) {
		// if !koala.ExistSession(req, "sessionID") {
		// 	w.Write([]byte("请先登录"))
		// } else {
		id := p.ParamUrl["id"]
		classes, err := showClass(id)
		if err != nil {
			w.Write([]byte("查看教学班失败\n" + err.Error()))
		} else {
			json, _ := json.Marshal(classes)
			w.Write([]byte("教学班\n"))
			w.Write([]byte(json))
		}
		// }
	})

	koala.Handle("/class/:id/introduction", func(p *koala.Params, w http.ResponseWriter, req *http.Request) {
		// if !koala.ExistSession(req, "sessionID") {
		// 	w.Write([]byte("请先登录"))
		// } else {
		id := p.ParamUrl["id"]
		introduction, err := showClassIntroduction(id)
		if err != nil {
			w.Write([]byte("查看课程介绍失败\n" + err.Error()))
		} else {
			json, _ := json.Marshal(introduction)
			w.Write([]byte("课程介绍\n"))
			w.Write([]byte(json))
		}
		// }
	})

	koala.Handle("/class/:id/teachingSyllabus", func(p *koala.Params, w http.ResponseWriter, req *http.Request) {
		// if !koala.ExistSession(req, "sessionID") {
		// 	w.Write([]byte("请先登录"))
		// } else {
		id := p.ParamUrl["id"]
		teachingSyllabus, err := showClassTeachingSyllabus(id)
		if err != nil {
			w.Write([]byte("查看课程大纲失败\n" + err.Error()))
		} else {
			json, _ := json.Marshal(teachingSyllabus)
			w.Write([]byte("课程大纲\n"))
			w.Write([]byte(json))
		}
		// }
	})

}
