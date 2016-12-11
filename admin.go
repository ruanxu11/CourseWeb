package main

import (
	"encoding/json"
	"io/ioutil"
	"koala"
	"net/http"
	"os"
)

func getIDandPassword() map[string]string {
	fi, err := os.Open("./admin")
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	data := make(map[string]string)
	err = json.Unmarshal(fd, &data)
	return data
}

func adminIDandPasswordCheck(id string, password string) bool {
	data := getIDandPassword()
	if id == data["id"] && password == data["password"] {
		return true
	}
	return false
}

func admincheck(w http.ResponseWriter, r *http.Request) bool {
	data := getIDandPassword()
	if !koala.CheckSession(r, map[string]interface{}{"password": data["password"], "id": data["id"]}) {
		return false
	}
	return true
}

func adminHandlers() {
	adminCourseHandlers()
	adminClassHandlers()
	adminStudentHandlers()
	adminTeachingAssistantHandlers()
	adminTeacherHandlers()

	koala.Get("/admin", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		koala.Render(w, "index_admin.html", map[string]interface{}{
			"title": "管理员登录",
		})
	})

	koala.Post("/admin/logout", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		if koala.ExistSession(r, "sessionID") {
			if session := koala.PeekSession(r, "sessionID"); session != nil {
				session.Destory()
			}
		}
		koala.Relocation(w, "/admin", "注销成功", "success")
	})

	koala.Post("/admin", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		id := p.ParamPost["id"][0]
		password := p.ParamPost["password"][0]
		if ok := adminIDandPasswordCheck(id, password); ok {
			session := koala.GetSession(r, w, "sessionID")
			session.Values["collection"] = "admin"
			session.Values["id"] = id
			session.Values["password"] = password
			koala.Relocation(w, "/admin/main", "欢迎管理员大佬登录", "success")
		} else {
			koala.Relocation(w, "/admin", "帐号或密码错误", "error")
		}
	})

	koala.Get("/admin/main", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		if !admincheck(w, r) {
			koala.Relocation(w, "/admin", "请先登录", "error")
			return
		}
		koala.Render(w, "admin.html", map[string]interface{}{
			"title": "管理员主界面",
		})
	})

}
