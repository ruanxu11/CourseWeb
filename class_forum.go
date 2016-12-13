package main

import (
	"koala"
	"log"
	"net/http"
	"time"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Post struct {
	ID               string
	ReplieNum        int
	LastModifyTime   string
	Poster           string
	PosterID         string
	PosterCollection string
	CreateTime       string
	Topic            string
	Content          string
	Replies          []Reply
}

type Reply struct {
	ID               string
	Poster           string
	PosterID         string
	PosterCollection string
	Time             string
	Content          string
}

func addClassPost(id string, post *Post) error {
	t := time.Now().Format("2006-01-02 15:04:05")
	post.ID = koala.HashString(time.Now().Format(time.UnixDate) + post.PosterID + post.Content)
	post.CreateTime = t
	post.LastModifyTime = t
	return mgoUpdate("class",
		bson.M{"_id": id},
		bson.M{"$push": bson.M{"forum": &post}})
}

func addClassPostReply(id string, postid string, reply *Reply) error {
	reply.ID = koala.HashString(time.Now().Format(time.UnixDate) + reply.PosterID + reply.Content)
	reply.Time = time.Now().Format("2006-01-02 15:04:05")
	return mgoUpdate("class",
		bson.M{"_id": id, "forum": bson.M{"$elemMatch": bson.M{"id": postid}}},
		bson.M{"$push": bson.M{"forum.$.replies": &reply}})
}

func removeClassPostByID(id string, postid string) error {
	return mgoUpdate("class",
		bson.M{"_id": id},
		bson.M{"$pull": bson.M{"forum": bson.M{"id": postid}}})
}

func removeClassPostReplyByID(id string, postid string, replyid string) error {
	return mgoUpdate("class",
		bson.M{"_id": id, "forum": bson.M{"$elemMatch": bson.M{"id": postid}}},
		bson.M{"$pull": bson.M{"forum.$.replies": bson.M{"id": replyid}}})
}

func getClassForum(id string) ([]map[string]interface{}, error) {
	forum := []map[string]interface{}{}
	q := func(c *mgo.Collection) (*mgo.ChangeInfo, error) {
		pipe := c.Pipe([]bson.M{
			{
				"$unwind": "$forum",
			},
			{
				"$project": bson.M{
					"id":               "$forum.id",
					"replienum":        "$forum.replienum",
					"lastmodifytime":   "$forum.lastmodifytime",
					"poster":           "$forum.poster",
					"posterid":         "$forum.posterid",
					"postercollection": "$forum.postercollection",
					"createtime":       "$forum.createtime",
					"topic":            "$forum.topic",
					"content":          "$forum.content",
					"replies":          "$forum.replies",
				},
			},
			{
				"$group": bson.M{
					"_id": bson.M{
						"id":               "$id",
						"replienum":        "$replienum",
						"lastmodifytime":   "$lastmodifytime",
						"poster":           "$poster",
						"posterid":         "$posterid",
						"postercollection": "$postercollection",
						"createtime":       "$createtime",
						"topic":            "$topic",
						"content":          "$content",
						"replies":          "$replies",
					},
				},
			},
			{
				"$sort": bson.M{
					"_id.createtime": -1,
				},
			},
		})
		iter := pipe.Iter()
		tag := bson.M{}
		for iter.Next(&tag) {
			forum = append(forum, tag["_id"].(bson.M))
		}
		if err := iter.Close(); err != nil {
			return nil, err
		}
		return nil, nil
	}
	_, err := withCollection("class", q)
	return forum, err
}

func getClassPost(id string, postid string) (post map[string]interface{}, err error) {
	q := func(c *mgo.Collection) (*mgo.ChangeInfo, error) {
		pipe := c.Pipe([]bson.M{
			{
				"$match": bson.M{
					"_id": id,
				},
			},
			{
				"$unwind": "$forum",
			},
			{
				"$project": bson.M{
					"id":               "$forum.id",
					"poster":           "$forum.poster",
					"posterid":         "$forum.posterid",
					"postercollection": "$forum.postercollection",
					"createtime":       "$forum.createtime",
					"topic":            "$forum.topic",
					"content":          "$forum.content",
					"replies":          "$forum.replies",
				},
			},
			{
				"$group": bson.M{
					"_id": bson.M{
						"id":               "$id",
						"poster":           "$poster",
						"posterid":         "$posterid",
						"postercollection": "$postercollection",
						"createtime":       "$createtime",
						"topic":            "$topic",
						"content":          "$content",
						"replies":          "$replies",
					},
				},
			},
			{
				"$match": bson.M{
					"_id.id": postid,
				},
			},
		})
		iter := pipe.Iter()
		tag := bson.M{}
		iter.Next(&tag)
		post = tag["_id"].(bson.M)
		if err := iter.Close(); err != nil {
			return nil, err
		}
		return nil, nil
	}
	_, err = withCollection("class", q)
	return post, err
}

func classForumHandlers() {
	koala.Get("/class/:id/forum", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		id := p.ParamUrl["id"]
		forum, err := getClassForum(id)
		if err != nil {
			koala.NotFound(w)
			return
		}
		koala.Render(w, "class_forum.html", map[string]interface{}{
			"title":  courseWeb,
			"id":     id,
			"forum":  forum,
			"powers": getPowersInClass(r, id),
		})
	})

	koala.Post("/class/:id/forum/add", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		id := p.ParamUrl["id"]
		powers := getPowersInClass(r, id)
		if !powers["ForumPost"] {
			koala.NotFound(w)
			return
		}
		topic := p.ParamPost["topic"][0]
		content := p.ParamPost["content"][0]
		var poster, PosterID, PosterCollection string
		if koala.ExistSession(r, "sessionID") {
			session := koala.GetSession(r, w, "sessionID")
			poster = session.Values["name"].(string)
			PosterID = session.Values["id"].(string)
			PosterCollection = session.Values["collection"].(string)
		}
		err := addClassPost(id, &Post{
			ReplieNum:        0,
			Poster:           poster,
			PosterID:         PosterID,
			PosterCollection: PosterCollection,
			Topic:            topic,
			Content:          content,
		})
		if err != nil {
			log.Println(err)
			koala.Relocation(w, "/class/"+id+"/forum", "发帖失败", "error")
		} else {
			koala.Relocation(w, "/class/"+id+"/forum", "发帖成功", "success")
		}
	})

	koala.Get("/class/:id/forum/remove/:postid", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		id := p.ParamUrl["id"]
		powers := getPowersInClass(r, id)
		if !powers["ForumPostRemove"] {
			koala.NotFound(w)
			return
		}
		postid := p.ParamUrl["postid"]
		err := removeClassPostByID(id, postid)
		if err != nil {
			log.Println(err)
			koala.Relocation(w, "/class/"+id+"/forum", "删除帖子失败", "error")
		} else {
			koala.Relocation(w, "/class/"+id+"/forum", "删除帖子成功", "success")
		}
	})

	koala.Get("/class/:id/forum/post/:postid", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		id := p.ParamUrl["id"]
		postid := p.ParamUrl["postid"]
		post, err := getClassPost(id, postid)
		log.Println(post)
		if err != nil {
			koala.NotFound(w)
			return
		}
		koala.Render(w, "class_forum_post.html", map[string]interface{}{
			"title":  courseWeb,
			"id":     id,
			"post":   post,
			"powers": getPowersInClass(r, id),
		})
	})

	koala.Post("/class/:id/forum/post/:postid/add", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		id := p.ParamUrl["id"]
		powers := getPowersInClass(r, id)
		if !powers["ForumReply"] {
			koala.NotFound(w)
			return
		}
		postid := p.ParamUrl["postid"]
		content := p.ParamPost["content"][0]
		var poster, PosterID, PosterCollection string
		if koala.ExistSession(r, "sessionID") {
			session := koala.GetSession(r, w, "sessionID")
			poster = session.Values["name"].(string)
			PosterID = session.Values["id"].(string)
			PosterCollection = session.Values["collection"].(string)
		}
		err := addClassPostReply(id, postid, &Reply{
			Poster:           poster,
			PosterID:         PosterID,
			PosterCollection: PosterCollection,
			Content:          content,
		})
		if err != nil {
			log.Println(err)
			koala.Relocation(w, "/class/"+id+"/forum/post/"+postid, "回复失败", "error")
		} else {
			koala.Relocation(w, "/class/"+id+"/forum/post/"+postid, "回复成功", "success")
		}
	})

	koala.Get("/class/:id/forum/post/:postid/remove/:replyid", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		id := p.ParamUrl["id"]
		powers := getPowersInClass(r, id)
		if !powers["ForumPostRemove"] {
			koala.NotFound(w)
			return
		}
		postid := p.ParamUrl["postid"]
		replyid := p.ParamUrl["replyid"]
		err := removeClassPostReplyByID(id, postid, replyid)
		if err != nil {
			log.Println(err)
			koala.Relocation(w, "/class/"+id+"/forum/post/"+postid, "删除发言失败", "error")
		} else {
			koala.Relocation(w, "/class/"+id+"/forum/post/"+postid, "删除发言成功", "success")
		}
	})
}
