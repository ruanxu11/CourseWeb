package main

import (
	"koala"
	"net/http"
	"strconv"
)

func adminTeachingAssistantHandlers() {
	koala.Get("/admin/teachingAssistant", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		if !admincheck(w, r) {
			koala.Relocation(w, "/admin", "请先登录", "error")
			return
		}
		teachingAssistants, _ := getTeachingAssistantsByPage(0)
		koala.Render(w, "admin_teachingAssistant.html", map[string]interface{}{
			"title":              "助教管理界面",
			"teachingAssistants": teachingAssistants,
			"page":               0,
			"items": []string{
				"_id",
				"name",
				"sex",
				"email",
				"phone",
			},
		})
	})

	koala.Post("/admin/teachingAssistant", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		if !admincheck(w, r) {
			koala.Relocation(w, "/admin", "请先登录", "error")
			return
		}
		page := 0
		submit := p.ParamPost["submit"][0]
		var teachingAssistants []map[string]interface{}
		switch submit {
		case "上一页":
			page, _ = strconv.Atoi(p.ParamPost["page"][0])
			if page != 0 {
				page--
			}
			teachingAssistants, _ = getTeachingAssistantsByPage(page)
		case "下一页":
			page, _ = strconv.Atoi(p.ParamPost["page"][0])
			page++
			teachingAssistants, _ = getTeachingAssistantsByPage(page)
		case "搜索":
			teachingAssistants, _ = searchTeachingAssistantsSelect(map[string]interface{}{
				"_id":   p.ParamPost["_id"][0],
				"name":  p.ParamPost["name"][0],
				"sex":   p.ParamPost["sex"][0],
				"email": p.ParamPost["email"][0],
				"phone": p.ParamPost["phone"][0],
			})
		case "增加":
			err := addTeachingAssistant(&TeachingAssistant{
				ID:       p.ParamPost["_id"][0],
				Password: p.ParamPost["_id"][0],
				Name:     p.ParamPost["name"][0],
				Sex:      p.ParamPost["sex"][0],
				Email:    p.ParamPost["email"][0],
				Phone:    p.ParamPost["phone"][0],
			})
			if err != nil {
				koala.Relocation(w, "/admin/teachingAssistant", "添加助教失败", "error")
				return
			}
			koala.Relocation(w, "/admin/teachingAssistant", "添加助教成功", "error")
			return
		case "显示全部":
			teachingAssistants, _ = searchTeachingAssistants(nil)
		default:
			page, _ = strconv.Atoi(p.ParamPost["page"][0])
			teachingAssistants, _ = getTeachingAssistantsByPage(page)
		}
		koala.Render(w, "admin_teachingAssistant.html", map[string]interface{}{
			"title":              "助教管理界面",
			"teachingAssistants": teachingAssistants,
			"page":               page,
			"items": []string{
				"_id",
				"name",
				"sex",
				"email",
				"phone",
			},
		})
	})
}
