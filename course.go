package main

import (
	"errors"
	"io/ioutil"
	"koala"
	"log"
	"net/http"
	"regexp"
	"strconv"

	iconv "github.com/djimenez/iconv-go"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Course struct {
	ID               string `bson:"_id"` // 课程代码
	Name             string // 课程名称
	College          string // 开课学院
	Credit           string // 学分
	HoursPerWeek     string // 周学时
	Type             string // 课程类型
	PreviousCourse   string // 预修要求
	Introduction     string // 课程简介
	TeachingSyllabus string // 教学大纲
}

func initCourses(no int) {
	id := strconv.Itoa(no)
	url := "http://10.202.78.13/html_kc/" + id + ".html"
	log.Println(id)
	r, err := http.Get(url)
	if err != nil {
		return
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	input := body
	out := make([]byte, len(input))
	out = out[:]
	iconv.Convert(input, out, "gb2312", "utf-8")
	str := string(out)

	reg, _ := regexp.Compile("[\\s\\S]*<span id=\"kczwmc\">(.*)</span>[\\s\\S]*")
	if s := reg.FindString(str); len(s) == 0 {
		return
	}
	Name := reg.ReplaceAllString(str, "$1")
	log.Println(Name)

	reg, _ = regexp.Compile("[\\s\\S]*<span id=\"kkxy\">(.*)</span>[\\s\\S]*")
	College := reg.ReplaceAllString(str, "$1")
	log.Println(College)

	reg, _ = regexp.Compile("[\\s\\S]*<span id=\"xf\">(.*)</span>[\\s\\S]*")
	Credit := reg.ReplaceAllString(str, "$1")
	log.Println(Credit)

	reg, _ = regexp.Compile("[\\s\\S]*<span id=\"zxs\">(.*)</span>[\\s\\S]*")
	HoursPerWeek := reg.ReplaceAllString(str, "$1")
	log.Println(HoursPerWeek)

	reg, _ = regexp.Compile("[\\s\\S]*<span id=\"kcgs\">(.*)</span>[\\s\\S]*")
	Type := reg.ReplaceAllString(str, "$1")
	log.Println(Type)

	reg, _ = regexp.Compile("[\\s\\S]*<span id=\"yxyq\">(.*)</span>[\\s\\S]*")
	PreviousCourse := reg.ReplaceAllString(str, "$1")
	log.Println(PreviousCourse)

	reg, _ = regexp.Compile("[\\s\\S]*<textarea name=\"kcjj\" readonly=\"readonly\" id=\"kcjj\" style=\"border-style:None;height:150px;width:100%;\">([\\s\\S]*?)</textarea>[\\s\\S]*")
	Introduction := ""
	if s := reg.FindString(str); len(s) == 0 {
		log.Println("not found")
	} else {
		Introduction = reg.ReplaceAllString(str, "$1")
		log.Println(Introduction)
	}

	reg, _ = regexp.Compile("[\\s\\S]*<textarea name=\"jxdg\" readonly=\"readonly\" id=\"jxdg\" style=\"border-style:None;height:150px;width:100%;\">([\\s\\S]*?)</textarea>[\\s\\S]*")
	TeachingSyllabus := ""
	if s := reg.FindString(str); len(s) == 0 {
		log.Println("not found")
	} else {
		TeachingSyllabus = reg.ReplaceAllString(str, "$1")
		log.Println(TeachingSyllabus)
	}

	course := &Course{
		ID:               id,
		Name:             Name,
		College:          College,
		HoursPerWeek:     HoursPerWeek,
		Type:             Type,
		PreviousCourse:   PreviousCourse,
		Introduction:     Introduction,
		TeachingSyllabus: TeachingSyllabus,
	}
	err = addCourse(course)
	if err != nil {
		log.Println(err)
	}
}

func addCourse(course *Course) error {
	return mgoInsert("course", &course)
}

func removeCourses(selector map[string]interface{}) (*mgo.ChangeInfo, error) {
	return mgoRemoveAll("course", selector)
}

func removeCourseByID(id string) (*mgo.ChangeInfo, error) {
	return mgoRemove("course", bson.M{"_id": id})
}

func updateCourses(selector map[string]interface{}, update map[string]interface{}) (*mgo.ChangeInfo, error) {
	return mgoUpdateAll("course", selector, update)
}

func updateCourseByID(id string, update map[string]interface{}) error {
	return mgoUpdate("course", bson.M{"_id": id}, update)
}

func searchCourses(selector map[string]interface{}) ([]map[string]interface{}, error) {
	return mgoFindAll("course", selector)
}

func getCourse(id string) (map[string]interface{}, error) {
	return mgoFind("course", bson.M{"_id": id})
}

func getCoursesByPage(page int) ([]map[string]interface{}, error) {
	return mgoFindByPage("course", page)
}

func searchCoursesSelect(selector map[string]interface{}) ([]map[string]interface{}, error) {
	return mgoSearchSelect("course", selector)
}

func getCourseIntroduction(id string) (string, error) {
	course, err := mgoFind("course", bson.M{"_id": id})
	if err != nil {
		return "", err
	}
	switch course["introduction"].(type) {
	case string:
		introduction := course["introduction"].(string)
		return introduction, err
	case nil:
	default:
		return "", errors.New("error type of Course.introduction")
	}
	return "", err
}

func updateCourseIntroduction(id string, introduction string) error {
	return mgoUpdate("course",
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"introduction": introduction}})
}

func getCourseTeachingSyllabus(id string) (string, error) {
	course, err := mgoFind("course", bson.M{"_id": id})
	if err != nil {
		return "", err
	}
	switch course["teachingsyllabus"].(type) {
	case string:
		teachingsyllabus := course["teachingsyllabus"].(string)
		return teachingsyllabus, err
	case nil:
	default:
		return "", errors.New("error type of Course.teachingsyllabus")
	}
	return "", err
}

func updateCourseTeachingSyllabus(id string, teachingsyllabus string) error {
	return mgoUpdate("course",
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"teachingsyllabus": teachingsyllabus}})
}

func courseHandlers() {
	koala.Get("/course/:id", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		id := p.ParamUrl["id"]
		course, err := getCourse(id)
		if err != nil {
			koala.NotFound(w)
			return
		}
		classes, err := getClassByCourseID(id)
		if err != nil {
			koala.NotFound(w)
			return
		}
		admin := admincheck(w, r)
		koala.Render(w, "course.html", map[string]interface{}{
			"title":   courseWeb,
			"course":  course,
			"classes": classes,
			"admin":   admin,
		})
	})

	koala.Post("/course/:id/class/add", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		if !admincheck(w, r) {
			koala.Relocation(w, "/admin", "请先登录", "error")
			return
		}
		CourseID := p.ParamUrl["id"]
		sClassRoomNum := p.ParamPost["ClassRoomNum"][0]
		ClassRoomNum, _ := strconv.Atoi(sClassRoomNum)
		ClassRooms := make([]ClassRoom, ClassRoomNum)
		for i := 0; i < ClassRoomNum; i++ {
			ClassRooms[i].Time = p.ParamPost["ClassRoomTime"][i]
			ClassRooms[i].Position = p.ParamPost["ClassRoomPosition"][i]
		}
		sTeacherNum := p.ParamPost["TeacherNum"][0]
		TeacherNum, _ := strconv.Atoi(sTeacherNum)
		Teachers := make([]TeacherInClass, TeacherNum)
		for i := 0; i < TeacherNum; i++ {
			Teachers[i].ID = p.ParamPost["TeacherID"][i]
			Teachers[i].Name = p.ParamPost["TeacherName"][i]
		}
		class := &Class{
			CourseID:   CourseID,
			Course:     p.ParamPost["Course"][0],
			Year:       p.ParamPost["Year"][0],
			Semester:   p.ParamPost["Semester"][0],
			ClassRooms: ClassRooms,
			Teachers:   Teachers,
		}
		err := addClass(class)
		if err != nil {
			koala.Relocation(w, "/course/"+CourseID, "添加教学班失败", "error")
			return
		}
		koala.Relocation(w, "/course/"+CourseID, "添加教学班成功", "success")
	})

	koala.Get("/course/:courseid/class/remove/:id", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		if !admincheck(w, r) {
			koala.Relocation(w, "/", "管理员账户才有权限", "error")
			return
		}
		id := p.ParamUrl["id"]
		CourseID := p.ParamUrl["courseid"]
		_, err := removeClass(id)
		if err != nil {
			koala.Relocation(w, "/course/"+CourseID, "删除教学班失败", "error")
			return
		}
		koala.Relocation(w, "/course/"+CourseID, "删除教学班成功", "success")
	})
}
