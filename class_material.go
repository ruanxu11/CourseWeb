package main

import (
	"io"
	"koala"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Material struct {
	Time     string // 资料发布时间
	Filename string // 资料名称
	Path     string // 资料存储路径
	Suffix   string // 资料类型
}

func addClassMaterial(id string, material *Material) error {
	material.Time = time.Now().Format("2006-01-02 15:04:05")
	return mgoUpdate("class",
		bson.M{"_id": id},
		bson.M{"$push": bson.M{"materials": &material}})
}

func removeClassMaterialByTime(id string, time string) error {
	return mgoUpdate("class",
		bson.M{"_id": id},
		bson.M{"$pull": bson.M{"materials": bson.M{"time": time}}})
}

func getClassMaterials(id string) ([]map[string]interface{}, error) {
	materials := []map[string]interface{}{}
	q := func(c *mgo.Collection) (*mgo.ChangeInfo, error) {
		pipe := c.Pipe([]bson.M{
			{
				"$unwind": "$materials",
			},
			{
				"$project": bson.M{
					"time":     "$materials.time",
					"filename": "$materials.filename",
					"path":     "$materials.path",
					"suffix":   "$materials.suffix",
				},
			},
			{
				"$group": bson.M{
					"_id": bson.M{
						"time":     "$time",
						"filename": "$filename",
						"path":     "$path",
						"suffix":   "$suffix",
					},
				},
			},
			{
				"$sort": bson.M{
					"_id.time": -1,
				},
			},
		})
		iter := pipe.Iter()
		tag := bson.M{}
		for iter.Next(&tag) {
			materials = append(materials, tag["_id"].(bson.M))
		}
		if err := iter.Close(); err != nil {
			return nil, err
		}
		return nil, nil
	}
	_, err := withCollection("class", q)
	return materials, err
}

func classMaterialHandlers() {
	koala.Get("/class/:id/material", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		id := p.ParamUrl["id"]
		materials, err := getClassMaterials(id)
		if err != nil {
			koala.NotFound(w)
			return
		}
		power := false
		log.Println("typesInClass:")
		if koala.ExistSession(r, "sessionID") {
			session := koala.GetSession(r, w, "sessionID")
			typesInClass := getTypeInClass(id, session.Values["collection"].(string), session.Values["id"].(string))
			if typesInClass == "teacher" {
				power = true
			}
		}
		koala.Render(w, "class_material.html", map[string]interface{}{
			"title":      courseWeb,
			"id":         id,
			"materials":  materials,
			"permission": true,
			"power":      power,
		})
	})

	koala.Post("/class/:id/material/add", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		id := p.ParamUrl["id"]
		power := false
		if koala.ExistSession(r, "sessionID") {
			session := koala.GetSession(r, w, "sessionID")
			typesInClass := getTypeInClass(id, session.Values["collection"].(string), session.Values["id"].(string))
			if typesInClass == "teacher" {
				power = true
			}
		}
		if !power {
			koala.NotFound(w)
			return
		}
		file, handle, err := r.FormFile("material")
		if err != nil {
			log.Println(err)
			koala.Relocation(w, "/class/"+id+"/material", "新增课程资料失败", "error")
			return
		}
		filename := handle.Filename
		suffix := path.Ext(filename)
		log.Println(suffix)
		courseid, err := getCourseID(id)
		if err != nil {
			log.Println(err)
			koala.Relocation(w, "/class/"+id+"/material", "不存在这门课程", "error")
			return
		}
		mgofilepath := "/material/" + courseid + "/" + filename
		filepath := "./static/upload/material/" + courseid + "/" + filename
		os.MkdirAll(path.Dir(filepath), 0777)
		f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Println(err)
			koala.Relocation(w, "/class/"+id+"/material", "新增课程资料失败", "error")
			return
		}
		_, err = io.Copy(f, file)
		if err != nil {
			log.Println(err)
			koala.Relocation(w, "/class/"+id+"/material", "新增课程资料失败", "error")
			return
		}
		defer f.Close()
		defer file.Close()
		err = addClassMaterial(id, &Material{
			Filename: filename,
			Path:     mgofilepath,
			Suffix:   suffix,
		})
		if err != nil {
			log.Println(err)
			koala.Relocation(w, "/class/"+id+"/material", "新增课程资料失败", "error")
			return
		}
		koala.Relocation(w, "/class/"+id+"/material", "新增课程资料成功", "success")
	})

	koala.Get("/class/:id/material/:time/remove", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		id := p.ParamUrl["id"]
		power := false
		time := p.ParamUrl["time"]
		if koala.ExistSession(r, "sessionID") {
			session := koala.GetSession(r, w, "sessionID")
			typesInClass := getTypeInClass(id, session.Values["collection"].(string), session.Values["id"].(string))
			if typesInClass == "teacher" {
				power = true
			}
		}
		if !power {
			koala.NotFound(w)
			return
		}
		err := removeClassMaterialByTime(id, time)
		if err != nil {
			log.Println(err)
			koala.Relocation(w, "/class/"+id+"/material", "删除课程资料失败", "error")
		} else {
			koala.Relocation(w, "/class/"+id+"/material", "删除课程资料成功", "success")
		}
	})
}
