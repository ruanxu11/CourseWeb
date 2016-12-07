package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Student struct {
	ID                  int `bson:"_id"`
	Password            string
	Name                string
	Email               string
	Phone               int
	College             string // 所在学院
	Department          string // 系
	Major               string // 专业
	AdministrationClass string
	行政班                 string `bson:"行政班"`
}

func addStudent(sid string, name string) {
	id, _ := strconv.Atoi(sid)
	password := sid
	student := Student{
		ID:       id,
		Password: password,
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
