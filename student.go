package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

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
	SecurityQuestions   []SecurityQuestion
}

func addStudent(id string, name string) {
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

func addStudents() {
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
		if len(data) == 2 {
			addStudent(data[0], data[1])
		}
	}
	fmt.Println("ok")
}

func showStudentInfo(collection string, id string) (*map[string]interface{}, error) {
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
