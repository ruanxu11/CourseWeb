package main

import (
	"koala"
	"net/http"
	"strconv"
)

func adminCourseHandlers() {
	koala.Get("/admin/course", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		if !admincheck(w, r) {
			koala.Relocation(w, "/admin", "请先登录", "error")
			return
		}
		courses, _ := getCoursesByPage(0)
		koala.Render(w, "admin_course.html", map[string]interface{}{
			"title":   "课程管理界面",
			"courses": courses,
			"page":    0,
			"items": []string{
				"_id",
				"name",
				"college",
				"type",
			},
		})
	})

	koala.Post("/admin/course", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		if !admincheck(w, r) {
			koala.Relocation(w, "/admin", "请先登录", "error")
			return
		}
		page, _ := strconv.Atoi(p.ParamPost["page"][0])
		submit := p.ParamPost["submit"][0]
		var courses []map[string]interface{}
		switch submit {
		case "上一页":
			if page != 0 {
				page--
			}
			courses, _ = getCoursesByPage(page)
		case "下一页":
			page++
			courses, _ = getCoursesByPage(page)
		case "搜索":
			courses, _ = searchCoursesSelect(map[string]interface{}{
				"_id":     p.ParamPost["_id"][0],
				"name":    p.ParamPost["name"][0],
				"college": p.ParamPost["college"][0],
				"type":    p.ParamPost["type"][0],
			})
		case "显示全部":
			courses, _ = searchCourses(nil)
		default:
			courses, _ = getCoursesByPage(page)
		}
		koala.Render(w, "admin_course.html", map[string]interface{}{
			"title":   "课程管理界面",
			"courses": courses,
			"page":    page,
			"items": []string{
				"_id",
				"name",
				"college",
				"type",
			},
		})
	})
}
