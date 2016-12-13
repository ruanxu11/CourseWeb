package main

import (
	"koala"
	"log"
	"net/http"

	"labix.org/v2/mgo/bson"
)

type PowersInClass struct {
	Forum                  []string `bson:"_id"`                    // 浏览讨论区
	ForumPost              []string `bson:"ForumPost"`              // 发帖
	ForumReply             []string `bson:"ForumReply"`             // 回帖
	ForumPostRemove        []string `bson:"ForumPostRemove"`        // 删帖
	MaterialAdd            []string `bson:"MaterialAdd"`            // 上传资料
	MaterialRemove         []string `bson:"MaterialRemove"`         // 删除资料
	MaterialDownload       []string `bson:"MaterialDownload"`       // 下载资料
	AssignmentAdd          []string `bson:"AssignmentAdd"`          // 增加作业
	AssignmentRemove       []string `bson:"AssignmentRemove"`       // 删除作业
	AssignmentUpdate       []string `bson:"AssignmentUpdate"`       // 更新作业
	AssignmentView         []string `bson:"AssignmentView"`         // 查看作业
	AssignmentCheck        []string `bson:"AssignmentCheck"`        // 批改作业
	AnnouncementAdd        []string `bson:"AnnouncementAdd"`        // 增加公告
	AnnouncementRemove     []string `bson:"AnnouncementRemove"`     // 删除公告
	AnnouncementUpdate     []string `bson:"AnnouncementUpdate"`     // 更改公告
	AnnouncementView       []string `bson:"AnnouncementView"`       // 查看公告
	TeachingSyllabusUpdate []string `bson:"TeachingSyllabusUpdate"` // 更改课程大纲
	IntroductionUpdate     []string `bson:"IntroductionUpdate"`     // 更改课程介绍
	StudentList            []string `bson:"StudentList"`            // 查看学生名单
}

type Powers struct {
	Forum                  bool // 浏览讨论区
	ForumPost              bool // 发帖
	ForumReply             bool // 回帖
	ForumPostRemove        bool // 删帖
	MaterialAdd            bool // 上传资料
	MaterialRemove         bool // 删除资料
	MaterialDownload       bool // 下载资料
	AssignmentAdd          bool // 增加作业
	AssignmentRemove       bool // 删除作业
	AssignmentUpdate       bool // 更新作业
	AssignmentView         bool // 查看作业
	AssignmentCheck        bool // 批改作业
	AnnouncementAdd        bool // 增加公告
	AnnouncementRemove     bool // 删除公告
	AnnouncementUpdate     bool // 更改公告
	AnnouncementView       bool // 查看公告
	TeachingSyllabusUpdate bool // 更改课程大纲
	IntroductionUpdate     bool // 更改课程介绍
	StudentList            bool // 查看学生名单
}

func initPowersInclass(powers *PowersInClass) {
	mgoUpdateAll("class", nil,
		bson.M{"$set": bson.M{"class.powers": &powers}})
}

func getTypeInClass(classid string, collection string, id string) string {
	var err error
	if collection == "student" {
		_, err = mgoFind("class", bson.M{"_id": classid, "students.id": id})
		if err != nil {
			log.Println(err)
			return "otherStudent"
		}
		return "student"
	} else if collection == "teacher" {
		_, err = mgoFind("class", bson.M{"_id": classid, "teachers.id": id})
		if err != nil {
			log.Println(err)
			return "otherTeacher"
		}
		return "teacher"
	} else if collection == "teachingAssistant" {
		_, err = mgoFind("class", bson.M{"_id": classid, "teachingassistantid": id})
		if err != nil {
			log.Println(err)
			return "otherTeachingAssistant"
		}
		return "teachingAssistant"
	}
	return "unknown"
}

func getPowersInClass(r *http.Request, classid string) map[string]bool {
	powers := make(map[string]bool)
	typesInClass := "others"
	if koala.ExistSession(r, "sessionID") {
		session := koala.PeekSession(r, "sessionID")
		if session.Values["collection"] == "admin" {
			for k := range globalPowers {
				powers[k] = true
			}
			return powers
		}
		typesInClass = getTypeInClass(classid, session.Values["collection"].(string), session.Values["id"].(string))
	}
	for k, v := range globalPowers {
		flag := false
		for _, types := range v {
			if typesInClass == types {
				flag = true
				break
			}
		}
		powers[k] = flag
	}
	return powers
}
