package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"routerKL"
	"utilKL"
)

func renderHTML(w http.ResponseWriter, file string, data interface{}) {
	t, err := template.New(file).ParseFiles("static/views/" + file)
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, data)
}

func RouteStart() {
	// files
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("static/js"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("static/css"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("static/images"))))
	http.Handle("/views/", http.StripPrefix("/views/", http.FileServer(http.Dir("static/views"))))
	http.Handle("/template/", http.StripPrefix("/template/", http.FileServer(http.Dir("static/template"))))

	routerKL.Handle("/login", func(w http.ResponseWriter, r *http.Request) {
		collection := routerKL.Param(r, "collection")
		id := routerKL.Param(r, "id")
		password := routerKL.Param(r, "password")
		if name, ok := loginCheck(collection, id, password); ok {
			session := utilKL.GetSession(r, w, "sessionID")
			if session.IsNew {
				session.Values["collection"] = collection
				session.Values["id"] = id
				session.Values["password"] = password
				session.Values["name"] = name
				log.Println(session)
				w.Write([]byte("登陆成功\n欢迎," + name + "同学\n"))
			} else {
				w.Write([]byte("您已经登陆了," + session.Values["name"].(string) + "同学, 请不要再次登陆\n"))
			}
		} else {
			w.Write([]byte("账号或密码错误"))
		}
	})

	routerKL.Handle("/logout", func(w http.ResponseWriter, r *http.Request) {
		if !utilKL.ExistSession(r, "sessionID") {
			w.Write([]byte("你根本没有登录啊逗比"))
		} else if session := utilKL.PeekSession(r, "sessionID"); session != nil {
			session.Destory()
			w.Write([]byte("注销成功"))
		}
	})

	routerKL.Handle("/changePWD", func(w http.ResponseWriter, r *http.Request) {
		collection := routerKL.Param(r, "collection")
		id := routerKL.Param(r, "id")
		newPassword := routerKL.Param(r, "newPassword")
		oldPassword := routerKL.Param(r, "oldPassword")
		err := changePWD(collection, id, newPassword, oldPassword)
		if err != nil {
			w.Write([]byte("修改密码失败\n" + err.Error()))
		} else {
			w.Write([]byte("修改密码成功"))
			if session := utilKL.PeekSession(r, "sessionID"); session != nil {
				w.Write([]byte("请重新登录"))
				session.Destory()
			}
		}
	})

	routerKL.Handle("/showSecurityQuestions", func(w http.ResponseWriter, r *http.Request) {
		collection := routerKL.Param(r, "collection")
		id := routerKL.Param(r, "id")
		securityQuestions, err := showSecurityQuestions(collection, id)
		if err != nil {
			w.Write([]byte("查看安全问题失败\n" + err.Error()))
		} else {
			json, _ := json.Marshal(securityQuestions)
			w.Write([]byte("安全问题\n"))
			w.Write([]byte(json))
		}
	})

	routerKL.Handle("/forgetPWD", func(w http.ResponseWriter, r *http.Request) {
		collection := routerKL.Param(r, "collection")
		id := routerKL.Param(r, "id")
		newPassword := routerKL.Param(r, "newPassword")
		ssecurityQuestionNum := routerKL.Param(r, "securityQuestionNum")
		securityQuestionNum, _ := strconv.Atoi(ssecurityQuestionNum)
		securityQuestions := make([]SecurityQuestion, securityQuestionNum)
		for i := 0; i < securityQuestionNum; i++ {
			securityQuestions[i].Question = routerKL.ParamX(r, "question", i)
			securityQuestions[i].Answer = routerKL.ParamX(r, "answer", i)
		}
		err := forgetPWD(collection, id, newPassword, securityQuestions)
		if err != nil {
			w.Write([]byte("修改密码失败\n" + err.Error()))
		} else {
			w.Write([]byte("修改密码成功\n"))
			if session := utilKL.PeekSession(r, "sessionID"); session != nil {
				w.Write([]byte("请重新登录"))
				session.Destory()
			}
		}
	})

	routerKL.Handle("/changeSecurityQuestionsbyOld", func(w http.ResponseWriter, r *http.Request) {
		collection := routerKL.Param(r, "collection")
		id := routerKL.Param(r, "id")
		newSSecurityQuestionNum := routerKL.Param(r, "newSecurityQuestionNum")
		newSecurityQuestionNum, _ := strconv.Atoi(newSSecurityQuestionNum)
		newSecurityQuestions := make([]SecurityQuestion, newSecurityQuestionNum)
		for i := 0; i < newSecurityQuestionNum; i++ {
			newSecurityQuestions[i].Question = routerKL.ParamX(r, "newQuestion", i)
			newSecurityQuestions[i].Answer = routerKL.ParamX(r, "newAnswer", i)
		}
		oldSSecurityQuestionNum := routerKL.Param(r, "oldSecurityQuestionNum")
		oldSecurityQuestionNum, _ := strconv.Atoi(oldSSecurityQuestionNum)
		oldSecurityQuestions := make([]SecurityQuestion, oldSecurityQuestionNum)
		for i := 0; i < oldSecurityQuestionNum; i++ {
			oldSecurityQuestions[i].Question = routerKL.ParamX(r, "oldQuestion", i)
			oldSecurityQuestions[i].Answer = routerKL.ParamX(r, "oldAnswer", i)
		}
		err := changeSecurityQuestionsbyOld(collection, id, newSecurityQuestions, oldSecurityQuestions)
		if err != nil {
			w.Write([]byte("修改密码失败\n" + err.Error()))
		} else {
			w.Write([]byte("修改密码成功\n"))
			if session := utilKL.PeekSession(r, "sessionID"); session != nil {
				session.Destory()
			}
			w.Write([]byte("请重新登录"))
		}
	})

	routerKL.Handle("/", func(w http.ResponseWriter, r *http.Request) {
		fileData, _ := ioutil.ReadFile("./readme.md")
		readme := strings.Replace(string(fileData), "\n", "<br>", -1)
		readme = strings.Replace(readme, "\t", "    ", -1)
		readme = strings.Replace(readme, " ", "&nbsp", -1)
		args := map[string]template.HTML{
			"title":  template.HTML("浙江大学课程网站系统"),
			"readme": template.HTML(string(readme)),
		}
		renderHTML(w, "index.html", args)
	})

	routerKL.RunWithLog(":2333")
}

// searchResults, _ := mgoFindAll("students", nil, 0, 0)
// json, _ := json.Marshal(searchResults)
// w.Write([]byte(json))
// var err error
// if err != nil {
// 	json, _ := json.Marshal(err)
// 	w.Write([]byte("err\n"))
// 	w.Write([]byte(json))
// 	log.Println(err)
// } else {
// 	w.Write([]byte("登陆成功"))
// 	w.Write([]byte("您已经登陆了," + Name + "同学"))
// }
