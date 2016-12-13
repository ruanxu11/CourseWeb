package main

import (
	"koala"
	"net/http"
	"strconv"
)

func adminTeacherHandlers() {
	koala.Get("/admin/teacher", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		if !admincheck(w, r) {
			koala.Relocation(w, "/admin", "请先登录", "error")
			return
		}
		teachers, _ := getTeachersByPage(0)
		koala.Render(w, "admin_teacher.html", map[string]interface{}{
			"title":    "老师管理界面",
			"teachers": teachers,
			"page":     0,
			"items": []string{
				"_id",
				"name",
				"sex",
				"email",
				"phone",
				"college",
				"department",
				"academicbackground",
				"academictitle",
				"researchdirections",
			},
		})
	})

	koala.Post("/admin/teacher", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		if !admincheck(w, r) {
			koala.Relocation(w, "/admin", "请先登录", "error")
			return
		}
		page := 0
		submit := p.ParamPost["submit"][0]
		var teachers []map[string]interface{}
		switch submit {
		case "上一页":
			page, _ = strconv.Atoi(p.ParamPost["page"][0])
			if page != 0 {
				page--
			}
			teachers, _ = getTeachersByPage(page)
		case "下一页":
			page, _ = strconv.Atoi(p.ParamPost["page"][0])
			page++
			teachers, _ = getTeachersByPage(page)
		case "搜索":
			teachers, _ = searchTeachersSelect(map[string]interface{}{
				"_id":                p.ParamPost["_id"][0],
				"name":               p.ParamPost["name"][0],
				"sex":                p.ParamPost["sex"][0],
				"email":              p.ParamPost["email"][0],
				"phone":              p.ParamPost["phone"][0],
				"college":            p.ParamPost["college"][0],
				"department":         p.ParamPost["department"][0],
				"academicbackground": p.ParamPost["academicbackground"][0],
				"academictitle":      p.ParamPost["academictitle"][0],
				"researchdirections": p.ParamPost["researchdirections"][0],
			})
		case "增加":
			err := addTeacher(&Teacher{
				ID:                 p.ParamPost["_id"][0],
				Password:           p.ParamPost["_id"][0],
				Name:               p.ParamPost["name"][0],
				Sex:                p.ParamPost["sex"][0],
				Email:              p.ParamPost["email"][0],
				Phone:              p.ParamPost["phone"][0],
				College:            p.ParamPost["college"][0],
				Department:         p.ParamPost["department"][0],
				AcademicBackground: p.ParamPost["academicbackground"][0],
				AcademicTitle:      p.ParamPost["academictitle"][0],
				ResearchDirections: p.ParamPost["researchdirections"][0],
			})
			if err != nil {
				koala.Relocation(w, "/admin/teacher", "添加老师失败", "error")
				return
			}
			koala.Relocation(w, "/admin/teacher", "添加老师成功", "error")
			return
		case "显示全部":
			teachers, _ = searchTeachers(nil)
		default:
			page, _ = strconv.Atoi(p.ParamPost["page"][0])
			teachers, _ = getTeachersByPage(page)
		}
		koala.Render(w, "admin_teacher.html", map[string]interface{}{
			"title":    "老师管理界面",
			"teachers": teachers,
			"page":     page,
			"items": []string{
				"_id",
				"name",
				"sex",
				"email",
				"phone",
				"college",
				"department",
				"academicbackground",
				"academictitle",
				"researchdirections",
			},
		})
	})
}
