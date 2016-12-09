package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"labix.org/v2/mgo/bson"

	iconv "github.com/djimenez/iconv-go"
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

func getCourses(id int) {
	no := strconv.Itoa(id)
	url := "http://10.202.78.13/html_kc/" + no + ".html"
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
		ID:               no,
		Name:             Name,
		College:          College,
		HoursPerWeek:     HoursPerWeek,
		Type:             Type,
		PreviousCourse:   PreviousCourse,
		Introduction:     Introduction,
		TeachingSyllabus: TeachingSyllabus,
	}
	err = mgoInsert("course", &course)
	if err != nil {
		log.Println("failed")
		log.Println(err)
	} else {
		log.Println("success")
	}
}

func showCourseIntroduction(id string) (string, error) {
	course, err := mgoFind("course", bson.M{"_id": id})
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

func showCourseTeachingSyllabus(id string) (string, error) {
	course, err := mgoFind("course", bson.M{"_id": id})
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
