package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"koala"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"log"

	"labix.org/v2/mgo/bson"
)

type SecurityQuestion struct {
	Question string
	Answer   string
}

func imitateLogin() {
	v := url.Values{}
	v.Set("ID", "3140102431")
	v.Set("Password", "3140102431")
	v.Set("Type", "student")
	body := ioutil.NopCloser(strings.NewReader(v.Encode())) //把form数据编下码
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://localhost:2333/login", body)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded") //这个一定要加，不加form的值post不过去，被坑了两小时
	fmt.Printf("%+v\n\n\n", req)                                        //看下发送的结构

	resp, err := client.Do(req) //发送
	defer resp.Body.Close()     //一定要关闭resp.Body
	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(data), err)
}

func loginCheck(collection string, id string, password string) (string, bool) {
	person, err := mgoFind(collection,
		bson.M{"_id": id, "password": password})
	if err != nil {
		return "", false
	}
	if person, ok := person["name"].(string); ok {
		return person, true
	}
	return "", true
}

func changePWD(collection string, id string, newPassword string, oldPassword string) error {
	if _, ok := loginCheck(collection, id, oldPassword); !ok {
		return errors.New("账号或密码错误")
	}
	return mgoUpdate(collection,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"password": newPassword}})
}

func showSecurityQuestions(collection string, id string) ([]interface{}, error) {
	person, err := mgoFind(collection,
		bson.M{"_id": id})
	securityQuestions := person["securityQuestions"].([]interface{})
	return securityQuestions, err
}

func forgetPWD(collection string, id string, newPassword string, securityQuestions []SecurityQuestion) error {
	person, _ := mgoFind(collection,
		bson.M{"_id": id, "securityQuestions": securityQuestions})
	switch person["securityQuestions"].(type) {
	case []interface{}:
	default:
		return errors.New("安全问题回答错误")
	}
	return mgoUpdate(collection,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"password": newPassword}})
}

func changeSecurityQuestionsbyOld(collection string, id string, newSecurityQuestions []SecurityQuestion, oldSecurityQuestions []SecurityQuestion) error {
	log.Println(bson.M{"_id": id, "securityQuestions": oldSecurityQuestions})
	person, _ := mgoFind(collection,
		bson.M{"_id": id, "securityQuestions": oldSecurityQuestions})
	switch person["securityQuestions"].(type) {
	case []interface{}:
	default:
		return errors.New("安全问题回答错误")
	}
	return mgoUpdate(collection,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"securityQuestions": newSecurityQuestions}})
}

func showPersonInfo(collection string, id string) (*map[string]interface{}, error) {
	if collection == "student" {
		return showStudentInfo(collection, id)
	} else if collection == "teacher" {
		return showTeacherInfo(collection, id)
	} else if collection == "teachingAssistant" {
		return showTeachingAssistantInfo(collection, id)
	}
	return nil, errors.New("类型错误")
}

func updatePersonInfo(collection string, id string, phone string, email string, introduction string) error {
	return mgoUpdate(collection,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"phone": phone, "email": email, "introduction": introduction}})
}

func accountHandler() {
	koala.Handle("/login/:collection/:id", func(p *koala.Params, w http.ResponseWriter, req *http.Request) {
		collection := p.ParamUrl["collection"]
		id := p.ParamUrl["id"]
		password := p.Param["password"][0]
		if name, ok := loginCheck(collection, id, password); ok {
			session := koala.GetSession(req, w, "sessionID")
			if session.IsNew {
				session.Values["collection"] = collection
				session.Values["id"] = id
				session.Values["password"] = password
				session.Values["name"] = name
				log.Println(session)
				w.Write([]byte("登陆成功\n欢迎," + name + "\n"))
			} else {
				w.Write([]byte("您已经登陆了," + name + ", 请不要再次登陆\n"))
			}
		} else {
			w.Write([]byte("账号或密码错误"))
		}
	})

	koala.Handle("/logout", func(p *koala.Params, w http.ResponseWriter, req *http.Request) {
		if !koala.ExistSession(req, "sessionID") {
			w.Write([]byte("你根本没有登录啊逗比"))
		} else if session := koala.PeekSession(req, "sessionID"); session != nil {
			session.Destory()
			w.Write([]byte("注销成功"))
		}
	})

	koala.Handle("/user/:collection/:id/password/change", func(p *koala.Params, w http.ResponseWriter, req *http.Request) {
		collection := p.ParamUrl["collection"]
		id := p.ParamUrl["id"]
		newPassword := p.Param["newPassword"][0]
		oldPassword := p.Param["oldPassword"][0]
		err := changePWD(collection, id, newPassword, oldPassword)
		if err != nil {
			w.Write([]byte("修改密码失败\n" + err.Error()))
		} else {
			w.Write([]byte("修改密码成功"))
			if session := koala.PeekSession(req, "sessionID"); session != nil {
				w.Write([]byte("请重新登录"))
				session.Destory()
			}
		}
	})

	koala.Handle("/user/:collection/:id/security/questions", func(p *koala.Params, w http.ResponseWriter, req *http.Request) {
		collection := p.ParamUrl["collection"]
		id := p.ParamUrl["id"]
		securityQuestions, err := showSecurityQuestions(collection, id)
		if err != nil {
			w.Write([]byte("查看安全问题失败\n" + err.Error()))
		} else {
			json, _ := json.Marshal(securityQuestions)
			w.Write([]byte("安全问题\n"))
			w.Write([]byte(json))
		}
	})

	koala.Handle("/user/:collection/:id/password/forget", func(p *koala.Params, w http.ResponseWriter, req *http.Request) {
		collection := p.ParamUrl["collection"]
		id := p.ParamUrl["id"]
		newPassword := p.Param["newPassword"][0]
		ssecurityQuestionNum := p.Param["securityQuestionNum"][0]
		securityQuestionNum, _ := strconv.Atoi(ssecurityQuestionNum)
		securityQuestions := make([]SecurityQuestion, securityQuestionNum)
		for i := 0; i < securityQuestionNum; i++ {
			securityQuestions[i].Question = p.Param["question"][i]
			securityQuestions[i].Answer = p.Param["answer"][i]
		}
		err := forgetPWD(collection, id, newPassword, securityQuestions)
		if err != nil {
			w.Write([]byte("修改密码失败\n" + err.Error()))
		} else {
			w.Write([]byte("修改密码成功\n"))
			if session := koala.PeekSession(req, "sessionID"); session != nil {
				w.Write([]byte("请重新登录"))
				session.Destory()
			}
		}
	})

	koala.Handle("/user/:collection/:id/security/questions/update", func(p *koala.Params, w http.ResponseWriter, req *http.Request) {
		collection := p.ParamUrl["collection"]
		id := p.ParamUrl["id"]
		newSSecurityQuestionNum := p.Param["newSecurityQuestionNum"][0]
		newSecurityQuestionNum, _ := strconv.Atoi(newSSecurityQuestionNum)
		newSecurityQuestions := make([]SecurityQuestion, newSecurityQuestionNum)
		for i := 0; i < newSecurityQuestionNum; i++ {
			newSecurityQuestions[i].Question = p.Param["newQuestion"][i]
			newSecurityQuestions[i].Answer = p.Param["newAnswer"][i]
		}
		oldSSecurityQuestionNum := p.Param["oldSecurityQuestionNum"][0]
		oldSecurityQuestionNum, _ := strconv.Atoi(oldSSecurityQuestionNum)
		oldSecurityQuestions := make([]SecurityQuestion, oldSecurityQuestionNum)
		for i := 0; i < oldSecurityQuestionNum; i++ {
			oldSecurityQuestions[i].Question = p.Param["oldQuestion"][i]
			oldSecurityQuestions[i].Answer = p.Param["oldAnswer"][i]
		}
		err := changeSecurityQuestionsbyOld(collection, id, newSecurityQuestions, oldSecurityQuestions)
		if err != nil {
			w.Write([]byte("修改安全问题失败\n" + err.Error()))
		} else {
			w.Write([]byte("修改安全问题成功成功\n"))
		}
	})

	koala.Handle("/user/:collection/:id", func(p *koala.Params, w http.ResponseWriter, req *http.Request) {
		collection := p.ParamUrl["collection"]
		id := p.ParamUrl["id"]
		personInfo, err := showPersonInfo(collection, id)
		if err != nil {
			w.Write([]byte("查看个人信息失败\n" + err.Error()))
		} else {
			json, _ := json.Marshal(personInfo)
			w.Write([]byte("个人信息\n"))
			w.Write([]byte(json))
		}
	})

	koala.Handle("/user/:collection/:id/update", func(p *koala.Params, w http.ResponseWriter, req *http.Request) {
		if session := koala.PeekSession(req, "sessionID"); session != nil {
			id := session.Values["id"].(string)
			collection := session.Values["collection"].(string)
			urlID := p.ParamUrl["id"]
			urlCollection := p.ParamUrl["collection"]
			if id == urlID && collection == urlCollection {
				phone := p.Param["phone"][0]
				email := p.Param["email"][0]
				introduction := p.Param["introduction"][0]
				err := updatePersonInfo(collection, id, phone, email, introduction)
				if err != nil {
					w.Write([]byte("修改个人信息失败\n" + err.Error()))
				} else {
					w.Write([]byte("修改个人信息成功\n"))
					personInfo, err := showPersonInfo(collection, id)
					if err != nil {
						w.Write([]byte("查看个人信息失败\n" + err.Error()))
					} else {
						json, _ := json.Marshal(personInfo)
						w.Write([]byte("个人信息\n"))
						w.Write([]byte(json))
					}
				}
				return
			}
		}
		w.Write([]byte("请先登录"))
	})
}
