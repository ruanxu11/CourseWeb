package main

import (
	"koala"
	"net/http"
	"strconv"
)

func adminClassHandlers() {
	koala.Get("/admin/class", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		if !admincheck(w, r) {
			koala.Relocation(w, "/admin", "请先登录", "error")
			return
		}
		classes, _ := getClassesByPage(0)
		koala.Render(w, "admin_class.html", map[string]interface{}{
			"title":   "教学班管理界面",
			"classes": classes,
			"page":    0,
			"items": []string{
				"_id",
				"courseid",
				"course",
				"year",
				"semester",
			},
		})
	})

	koala.Post("/admin/class", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		if !admincheck(w, r) {
			koala.Relocation(w, "/admin", "请先登录", "error")
			return
		}
		page, _ := strconv.Atoi(p.ParamPost["page"][0])
		submit := p.ParamPost["submit"][0]
		var classes []map[string]interface{}
		switch submit {
		case "上一页":
			if page != 0 {
				page--
			}
			classes, _ = getClassesByPage(page)
		case "下一页":
			page++
			classes, _ = getClassesByPage(page)
		case "搜索":
			classes, _ = searchClassesSelect(map[string]interface{}{
				"_id":      p.ParamPost["_id"][0],
				"courseid": p.ParamPost["courseid"][0],
				"course":   p.ParamPost["course"][0],
				"year":     p.ParamPost["year"][0],
				"semester": p.ParamPost["semester"][0],
			})
		case "显示全部":
			classes, _ = getClasses()
		default:
			classes, _ = getClassesByPage(page)
		}
		koala.Render(w, "admin_class.html", map[string]interface{}{
			"title":   "教学班管理界面",
			"classes": classes,
			"page":    page,
			"items": []string{
				"_id",
				"courseid",
				"course",
				"year",
				"semester",
			},
		})
	})
}
