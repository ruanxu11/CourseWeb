package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Student struct {
	ID                  string `bson:"_id"`
	Password            string
	Name                string
	Sex                 string
	Email               string
	Phone               string
	College             string // 所在学院
	Department          string // 系
	Major               string // 专业
	AdministrationClass string // 行政班
	Introduction        string
	SecurityQuestions   map[string]string
}

func initStudents() {
	f, err := os.Open("student.txt") //打开文件
	defer f.Close()                  //打开文件出错处理
	if err != nil {
		log.Println(err)
	}
	buff := bufio.NewReader(f) //读入缓存
	for {
		line, err := buff.ReadString('\n') //以'\n'为结束符读入一行
		if err != nil || io.EOF == err {
			break
		}
		line = line[0 : len(line)-1]
		data := strings.Split(line, " ")
		if len(data) == 3 {
			addStudent(&Student{
				ID:       data[0],
				Password: data[0],
				Name:     data[1],
				College:  data[2],
				SecurityQuestions: map[string]string{
					"question1": "你的学号是多少呀",
					"answer1":   data[0],
					"question2": "你叫什么呀",
					"answer2":   data[1],
					"question3": "你是什么学院的呀",
					"answer3":   data[2],
				},
			})
		}
	}
	fmt.Println("ok")
}

func addStudentByIDandName(id string, name string) {
	student := Student{
		ID:       id,
		Password: id,
		Name:     name,
	}
	err := mgoInsert("student", &student)
	if err != nil {
		log.Println("failed")
		log.Println(err)
	} else {
		log.Println("success")
	}
}

func addStudent(student *Student) error {
	return mgoInsert("student", &student)
}

func removeStudents(selector map[string]interface{}) (*mgo.ChangeInfo, error) {
	return mgoRemoveAll("student", selector)
}

func removeStudentByID(id string) (*mgo.ChangeInfo, error) {
	return mgoRemove("student", bson.M{"_id": id})
}

func updateStudents(selector map[string]interface{}, update map[string]interface{}) (*mgo.ChangeInfo, error) {
	return mgoUpdateAll("student", selector, update)
}

func updateStudentByID(id string, update map[string]interface{}) error {
	return mgoUpdate("student", bson.M{"_id": id}, update)
}

func searchStudents(selector map[string]interface{}) ([]map[string]interface{}, error) {
	return mgoFindAll("student", selector)
}

func getStudentInfo(collection string, id string) (*map[string]interface{}, error) {
	person, err := mgoFind(collection,
		bson.M{"_id": id})
	personInfo := make(map[string]interface{})
	personInfo["id"] = person["_id"]
	personInfo["name"] = person["name"]
	personInfo["sex"] = person["sex"]
	personInfo["email"] = person["email"]
	personInfo["phone"] = person["phone"]
	personInfo["college"] = person["college"]
	personInfo["department"] = person["department"]
	personInfo["major"] = person["major"]
	personInfo["administrationclass"] = person["administrationclass"]
	personInfo["introduction"] = person["introduction"]
	return &personInfo, err
}
