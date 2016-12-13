package main

import (
	"koala"
	"net/http"
	"strconv"
)

func adminStudentHandlers() {
	koala.Get("/admin/student", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		if !admincheck(w, r) {
			koala.Relocation(w, "/admin", "请先登录", "error")
			return
		}
		students, _ := getStudentsByPage(0)
		koala.Render(w, "admin_student.html", map[string]interface{}{
			"title":    "学生管理界面",
			"students": students,
			"page":     0,
			"items": []string{
				"_id",
				"name",
				"sex",
				"email",
				"phone",
				"college",
				"department",
				"major",
				"administrationclass",
			},
		})
	})

	koala.Post("/admin/student", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		if !admincheck(w, r) {
			koala.Relocation(w, "/admin", "请先登录", "error")
			return
		}
		page := 0
		submit := p.ParamPost["submit"][0]
		var students []map[string]interface{}
		switch submit {
		case "上一页":
			page, _ = strconv.Atoi(p.ParamPost["page"][0])
			if page != 0 {
				page--
			}
			students, _ = getStudentsByPage(page)
		case "下一页":
			page, _ = strconv.Atoi(p.ParamPost["page"][0])
			page++
			students, _ = getStudentsByPage(page)
		case "搜索":
			students, _ = searchStudentsSelect(map[string]interface{}{
				"_id":                 p.ParamPost["_id"][0],
				"name":                p.ParamPost["name"][0],
				"sex":                 p.ParamPost["sex"][0],
				"email":               p.ParamPost["email"][0],
				"phone":               p.ParamPost["phone"][0],
				"college":             p.ParamPost["college"][0],
				"department":          p.ParamPost["department"][0],
				"major":               p.ParamPost["major"][0],
				"administrationclass": p.ParamPost["administrationclass"][0],
			})
		case "增加":
			id := p.ParamPost["_id"][0]
			err := addStudent(&Student{
				ID:                  id,
				Password:            id[len(id)-4 : len(id)],
				Name:                p.ParamPost["name"][0],
				Sex:                 p.ParamPost["sex"][0],
				Email:               p.ParamPost["email"][0],
				Phone:               p.ParamPost["phone"][0],
				College:             p.ParamPost["college"][0],
				Department:          p.ParamPost["department"][0],
				Major:               p.ParamPost["major"][0],
				AdministrationClass: p.ParamPost["administrationclass"][0],
			})
			if err != nil {
				koala.Relocation(w, "/admin/student", "添加学生失败\n"+err.Error(), "error")
				return
			}
			koala.Relocation(w, "/admin/student", "添加学生成功", "error")
			return
		case "显示全部":
			students, _ = searchStudents(nil)
		default:
			page, _ = strconv.Atoi(p.ParamPost["page"][0])
			students, _ = getStudentsByPage(page)
		}
		koala.Render(w, "admin_student.html", map[string]interface{}{
			"title":    "学生管理界面",
			"students": students,
			"page":     page,
			"items": []string{
				"_id",
				"name",
				"sex",
				"email",
				"phone",
				"college",
				"department",
				"major",
				"administrationclass",
			},
		})
	})
}
