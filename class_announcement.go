package main

import (
	"koala"
	"log"
	"net/http"
	"time"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Announcement struct {
	Time    string // 公告发布时间
	Title   string // 公告标题
	Content string // 公告内容
}

func addClassAnnouncement(id string, announcement *Announcement) error {
	announcement.Time = time.Now().Format("2006-01-02 15:04:05")
	return mgoUpdate("class",
		bson.M{"_id": id},
		bson.M{"$push": bson.M{"announcements": &announcement}})
}

func removeClassAnnouncementByTime(id string, time string) error {
	return mgoUpdate("class",
		bson.M{"_id": id},
		bson.M{"$pull": bson.M{"announcements": bson.M{"time": time}}})
}

func updateClassAnnouncementByTime(id string, time string, title string, content string) error {
	return mgoUpdate("class",
		bson.M{"_id": id, "announcements.time": time},
		bson.M{"$set": bson.M{"announcements.$.title": title, "announcements.$.content": content}})
}

func getClassAnnouncements(id string) ([]map[string]interface{}, error) {
	announcements := []map[string]interface{}{}
	q := func(c *mgo.Collection) (*mgo.ChangeInfo, error) {
		pipe := c.Pipe([]bson.M{
			{
				"$unwind": "$announcements",
			},
			{
				"$project": bson.M{
					"time":    "$announcements.time",
					"title":   "$announcements.title",
					"content": "$announcements.content",
				},
			},
			{
				"$group": bson.M{
					"_id": bson.M{
						"time":    "$time",
						"title":   "$title",
						"content": "$content",
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
			announcements = append(announcements, tag["_id"].(bson.M))
		}
		if err := iter.Close(); err != nil {
			return nil, err
		}
		return nil, nil
	}
	_, err := withCollection("class", q)
	return announcements, err
}

func classAnnouncementHandlers() {
	koala.Get("/class/:id/announcement", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		id := p.ParamUrl["id"]
		announcements, err := getClassAnnouncements(id)
		if err != nil {
			koala.NotFound(w)
			return
		}
		koala.Render(w, "class_announcement.html", map[string]interface{}{
			"title":         courseWeb,
			"id":            id,
			"announcements": announcements,
			"powers":        getPowersInClass(r, id),
		})
	})

	koala.Post("/class/:id/announcement", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		submit := p.ParamPost["submit"][0]
		id := p.ParamUrl["id"]
		time := p.ParamPost["time"][0]
		powers := getPowersInClass(r, id)
		if submit == "删除课程公告" {
			if !powers["AnnouncementRemove"] {
				koala.NotFound(w)
				return
			}
			err := removeClassAnnouncementByTime(id, time)
			if err != nil {
				log.Println(err)
				koala.Relocation(w, "/class/"+id+"/announcement", "删除课程公告失败", "error")
			} else {
				koala.Relocation(w, "/class/"+id+"/announcement", "删除课程公告成功", "success")
			}
		} else if submit == "更改课程公告" {
			if !powers["AnnouncementUpdate"] {
				koala.NotFound(w)
				return
			}
			title := p.ParamPost["title"][0]
			content := p.ParamPost["content"][0]
			err := updateClassAnnouncementByTime(id, time, title, content)
			if err != nil {
				log.Println(err)
				koala.Relocation(w, "/class/"+id+"/announcement", "更改课程公告失败", "error")
			} else {
				koala.Relocation(w, "/class/"+id+"/announcement", "更改课程公告成功", "success")
			}
		}
	})

	koala.Post("/class/:id/announcement/add", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		id := p.ParamUrl["id"]
		powers := getPowersInClass(r, id)
		if !powers["AnnouncementAdd"] {
			koala.NotFound(w)
			return
		}
		title := p.ParamPost["title"][0]
		content := p.ParamPost["content"][0]
		err := addClassAnnouncement(id, &Announcement{
			Title:   title,
			Content: content,
		})
		if err != nil {
			log.Println(err)
			koala.Relocation(w, "/class/"+id+"/announcement", "新增课程公告失败", "error")
		} else {
			koala.Relocation(w, "/class/"+id+"/announcement", "新增课程公告成功", "success")
		}
	})
}
