package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"koala"
	"net/http"
	"net/url"
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
	r, _ := http.NewRequest("POST", "http://localhost:2333/login", body)

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded") //这个一定要加，不加form的值post不过去，被坑了两小时
	fmt.Printf("%+v\n\n\n", r)                                        //看下发送的结构

	resp, err := client.Do(r) //发送
	defer resp.Body.Close()   //一定要关闭resp.Body
	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(data), err)
}

func loginCheck(collection string, id string, password string) (string, bool) {
	User, err := mgoFind(collection,
		bson.M{"_id": id, "password": password})
	if err != nil {
		return "", false
	}
	if User, ok := User["name"].(string); ok {
		return User, true
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

func getSecurityQuestions(collection string, id string) ([]string, error) {
	User, err := mgoFindSelect(collection,
		bson.M{"_id": id},
		bson.M{"securityquestions.question1": 1, "securityquestions.question2": 1, "securityquestions.question3": 1, "_id": 0})
	if err != nil {
		return nil, err
	}
	securityQuestions := User["securityquestions"].(map[string]interface{})
	log.Println(securityQuestions)
	question1 := securityQuestions["question1"].(string)
	question2 := securityQuestions["question2"].(string)
	question3 := securityQuestions["question3"].(string)
	return []string{
		question1,
		question2,
		question3}, err
}

func checkSecurityQuestions(collection string, id string, newPassword string, answer1 string, answer2 string, answer3 string) error {
	_, err := mgoFind(collection,
		bson.M{"_id": id,
			"securityquestions.answer1": answer1,
			"securityquestions.answer2": answer2,
			"securityquestions.answer3": answer3})
	if err != nil {
		return err
	}
	return mgoUpdate(collection,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"password": newPassword}})
}

func updateSecurityQuestionsByPassword(collection string, id string, password string, securityQuestions map[string]string) error {
	_, err := mgoFind(collection,
		bson.M{"_id": id, "password": password})
	if err != nil {
		log.Println(err)
		return errors.New("密码错误")
	}
	err = mgoUpdate(collection,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{
			"securityquestions.question1": securityQuestions["question1"],
			"securityquestions.question2": securityQuestions["question2"],
			"securityquestions.question3": securityQuestions["question3"],
			"securityquestions.answer1":   securityQuestions["answer1"],
			"securityquestions.answer2":   securityQuestions["answer2"],
			"securityquestions.answer3":   securityQuestions["answer3"],
		}})
	if err != nil {
		log.Println(err)
		return errors.New("更新安全问题失败")
	}
	return nil
}

func updateSecurityQuestionsBySecurityQuestions(collection string, id string, oldSecurityQuestions map[string]string, securityQuestions map[string]string) error {
	_, err := mgoFind(collection,
		bson.M{"_id": id,
			"securityquestions.answer1": oldSecurityQuestions["answer1"],
			"securityquestions.answer2": oldSecurityQuestions["answer2"],
			"securityquestions.answer3": oldSecurityQuestions["answer3"]})
	if err != nil {
		log.Println(err)
		return errors.New("原安全问题回答错误")
	}
	err = mgoUpdate(collection,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{
			"securityquestions.question1": securityQuestions["question1"],
			"securityquestions.question2": securityQuestions["question2"],
			"securityquestions.question3": securityQuestions["question3"],
			"securityquestions.answer1":   securityQuestions["answer1"],
			"securityquestions.answer2":   securityQuestions["answer2"],
			"securityquestions.answer3":   securityQuestions["answer3"],
		}})
	if err != nil {
		log.Println(err)
		return errors.New("更新安全问题失败")
	}
	return nil
}

func getUserInfo(collection string, id string) (*map[string]interface{}, error) {
	if collection == "student" {
		return getStudentInfo(collection, id)
	} else if collection == "teacher" {
		return getTeacherInfo(collection, id)
	} else if collection == "teachingAssistant" {
		return getTeachingAssistantInfo(collection, id)
	}
	return nil, errors.New("类型错误")
}

func updateUserInfo(collection string, id string, phone string, email string, introduction string) error {
	return mgoUpdate(collection,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"phone": phone, "email": email, "introduction": introduction}})
}

func userHandlers() {
	koala.Post("/login", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		collection := p.ParamPost["collection"][0]
		id := p.ParamPost["id"][0]
		password := p.ParamPost["password"][0]
		if name, ok := loginCheck(collection, id, password); ok {
			session := koala.GetSession(r, w, "sessionID")
			session.Values["collection"] = collection
			session.Values["id"] = id
			session.Values["password"] = password
			session.Values["name"] = name
			url := "/user/" + collection + "/" + id
			koala.Relocation(w, url, "欢迎，"+name, "success")
		} else {
			koala.Relocation(w, "/", "帐号或密码错误", "error")
		}
	})

	koala.Post("/logout", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		if koala.ExistSession(r, "sessionID") {
			if session := koala.PeekSession(r, "sessionID"); session != nil {
				session.Destory()
			}
		}
		koala.Relocation(w, "/", "注销成功", "success")
	})

	koala.Get("/user/:collection/:id", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		collection := p.ParamUrl["collection"]
		id := p.ParamUrl["id"]
		userInfo, err := getUserInfo(collection, id)
		if err != nil {
			koala.NotFound(w)
			return
		}
		self := false
		login := false
		if session := koala.PeekSession(r, "sessionID"); session != nil {
			login = true
			if koala.CheckSession(r, map[string]interface{}{"collection": collection, "id": id}) {
				self = true
			}
		}
		koala.Render(w, collection+".html", map[string]interface{}{
			"title":      (*userInfo)["name"].(string) + "的个人主页",
			"self":       self,
			"login":      login,
			"userInfo":   userInfo,
			"collection": collection,
		})
	})

	koala.Post("/user/:collection/:id/password/forget", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		collection := p.ParamUrl["collection"]
		id := p.ParamUrl["id"]
		answer1 := p.ParamPost["answer1"][0]
		answer2 := p.ParamPost["answer2"][0]
		answer3 := p.ParamPost["answer3"][0]
		newPassword := p.ParamPost["newPassword"][0]
		err := checkSecurityQuestions(collection, id, newPassword, answer1, answer2, answer3)
		if err != nil {
			log.Println(err)
			koala.Relocation(w, "/forget/password/id", "修改密码失败", "error")
		} else {
			koala.Relocation(w, "/", "修改密码成功", "success")
		}
	})

	koala.Post("/user/:collection/:id/update", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		id := p.ParamUrl["id"]
		collection := p.ParamUrl["collection"]
		if !koala.CheckSession(r, map[string]interface{}{"collection": collection, "id": id}) {
			koala.Relocation(w, "/", "请先登录", "error")
			return
		}
		phone := p.ParamPost["phone"][0]
		email := p.ParamPost["email"][0]
		introduction := p.ParamPost["introduction"][0]
		err := updateUserInfo(collection, id, phone, email, introduction)
		if err != nil {
			log.Println(err)
			koala.Relocation(w, "/user/"+collection+"/"+id, "修改个人信息失败", "error")
		} else {
			koala.Relocation(w, "/user/"+collection+"/"+id, "修改个人信息成功", "success")
		}
	})

	koala.Post("/user/:collection/:id/password/update", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		collection := p.ParamUrl["collection"]
		id := p.ParamUrl["id"]
		oldPassword := p.ParamPost["oldPassword"][0]
		newPassword := p.ParamPost["newPassword"][0]
		newPasswordR := p.ParamPost["newPasswordR"][0]
		if newPassword != newPasswordR {
			koala.Relocation(w, "/user/"+collection+"/"+id, "两次输入的密码不一致", "error")
			return
		}
		err := changePWD(collection, id, newPassword, oldPassword)
		if err != nil {
			log.Println(err)
			koala.Relocation(w, "/user/"+collection+"/"+id, "修改密码失败", "error")
		} else {
			koala.Relocation(w, "/", "修改密码成功, 请重新登录", "success")
			if session := koala.PeekSession(r, "sessionID"); session != nil {
				session.Destory()
			}
		}
	})

	koala.Get("/user/:collection/:id/security/questions/update/by/password", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		id := p.ParamUrl["id"]
		collection := p.ParamUrl["collection"]
		if !koala.CheckSession(r, map[string]interface{}{"collection": collection, "id": id}) {
			koala.Relocation(w, "/", "请先登录", "error")
			return
		}
		koala.Render(w, "changeSecurityQuestionsByPassword.html", map[string]interface{}{
			"title":      "修改安全问题",
			"id":         id,
			"collection": collection,
		})
	})

	koala.Post("/user/:collection/:id/security/questions/update/by/password", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		collection := p.ParamUrl["collection"]
		id := p.ParamUrl["id"]
		err := updateSecurityQuestionsByPassword(collection, id, p.ParamPost["password"][0], map[string]string{
			"question1": p.ParamPost["question1"][0],
			"question2": p.ParamPost["question2"][0],
			"question3": p.ParamPost["question3"][0],
			"answer1":   p.ParamPost["answer1"][0],
			"answer2":   p.ParamPost["answer2"][0],
			"answer3":   p.ParamPost["answer3"][0],
		})
		if err != nil {
			koala.Relocation(w, "/user/"+collection+"/"+id+"/security/questions/update/by/password", err.Error(), "error")
		} else {
			koala.Relocation(w, "/user/"+collection+"/"+id, "修改安全问题成功", "success")
		}
	})

	koala.Get("/user/:collection/:id/security/questions/update/by/security/questions", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		id := p.ParamUrl["id"]
		collection := p.ParamUrl["collection"]
		securityQuestions, err := getSecurityQuestions(collection, id)
		if err != nil {
			log.Println(err)
			koala.Relocation(w, "/", "不存在该账号", "error")
		}
		if !koala.CheckSession(r, map[string]interface{}{"collection": collection, "id": id}) {
			koala.Relocation(w, "/", "请先登录", "error")
			return
		}
		koala.Render(w, "changeSecurityQuestionsBySecurityQuestions.html", map[string]interface{}{
			"title":      "修改安全问题",
			"id":         id,
			"collection": collection,
			"question1":  securityQuestions[0],
			"question2":  securityQuestions[1],
			"question3":  securityQuestions[2],
		})
	})

	koala.Post("/user/:collection/:id/security/questions/update/by/security/questions", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		collection := p.ParamUrl["collection"]
		id := p.ParamUrl["id"]
		err := updateSecurityQuestionsBySecurityQuestions(collection, id, map[string]string{
			"answer1": p.ParamPost["oldAnswer1"][0],
			"answer2": p.ParamPost["oldAnswer2"][0],
			"answer3": p.ParamPost["oldAnswer3"][0],
		}, map[string]string{
			"question1": p.ParamPost["question1"][0],
			"question2": p.ParamPost["question2"][0],
			"question3": p.ParamPost["question3"][0],
			"answer1":   p.ParamPost["answer1"][0],
			"answer2":   p.ParamPost["answer2"][0],
			"answer3":   p.ParamPost["answer3"][0],
		})
		if err != nil {
			koala.Relocation(w, "/user/"+collection+"/"+id+"/security/questions/update/by/security/questions", err.Error(), "error")
		} else {
			koala.Relocation(w, "/user/"+collection+"/"+id, "修改安全问题成功", "success")
		}
	})
}
