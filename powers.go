package main

import (
	"errors"
	"koala"
	"log"
	"net/http"

	"strconv"

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

// 权限缺省值 class >> course >> teacher >> default
type PowerTable struct {
	ForumView              int `bson:"ForumView"`        // 浏览讨论区
	ForumPost              int `bson:"ForumPost"`        // 发帖
	ForumReply             int `bson:"ForumReply"`       // 回帖
	ForumPostRemove        int `bson:"ForumPostRemove"`  // 删帖
	MaterialAdd            int `bson:"MaterialAdd"`      // 上传资料
	MaterialRemove         int `bson:"MaterialRemove"`   // 删除资料
	MaterialDownload       int `bson:"MaterialDownload"` // 下载资料
	AssignmentAdd          int `bson:"AssignmentAdd"`    // 增加作业
	AssignmentRemove       int `bson:"AssignmentRemove"` // 删除作业
	AssignmentUpdate       int `bson:"AssignmentUpdate"` // 更新作业
	AssignmentView         int `bson:"AssignmentView"`   // 查看作业
	AssignmentCheck        int `bson:"AssignmentCheck"`  // 批改作业
	AssignmentDo           int `bson:"AssignmentDo"`
	AnnouncementAdd        int `bson:"AnnouncementAdd"`        // 增加公告
	AnnouncementRemove     int `bson:"AnnouncementRemove"`     // 删除公告
	AnnouncementUpdate     int `bson:"AnnouncementUpdate"`     // 更改公告
	AnnouncementView       int `bson:"AnnouncementView"`       // 查看公告
	TeachingSyllabusUpdate int `bson:"TeachingSyllabusUpdate"` // 更改课程大纲
	IntroductionUpdate     int `bson:"IntroductionUpdate"`     // 更改课程介绍
	StudentList            int `bson:"StudentList"`            // 查看学生名单
	MakeTeam               int `bson:"MakeTeam"`
	PowersControl          int `bson:"PowersControl"`
}

var defaultPowerTables = PowerTable{
	// admin|teacher|teachingAssistant|student|otherTeacher|otherTeachingAssistant|otherStudent|others
	ForumView:              0xff, // "11111111",
	ForumPost:              0xf0, // "11110000",
	ForumReply:             0xf0, // "11110000",
	ForumPostRemove:        0xe0, // "11100000",
	MaterialAdd:            0xc0, // "11000000",
	MaterialRemove:         0xc0, // "11000000",
	MaterialDownload:       0xf0, // "11110000",
	AssignmentAdd:          0xe0, // "11100000",
	AssignmentRemove:       0xe0, // "11100000",
	AssignmentUpdate:       0xe0, // "11100000",
	AssignmentView:         0xf0, // "11110000",
	AssignmentCheck:        0xc0, // "11000000",
	AssignmentDo:           0x10, // "00010000",
	AnnouncementAdd:        0xc0, // "11000000",
	AnnouncementRemove:     0xc0, // "11000000",
	AnnouncementUpdate:     0xc0, // "11000000",
	AnnouncementView:       0xf0, // "11110000",
	TeachingSyllabusUpdate: 0xc0, // "11000000",
	IntroductionUpdate:     0xc0, // "11000000",
	StudentList:            0xf0, // "11110000",
	MakeTeam:               0xe0, // "11100000",
	PowersControl:          0xc0, // "11000000",
}

var defaultPowerTable = map[string]int{
	// admin|teacher|teachingAssistant|student|otherTeacher|otherTeachingAssistant|otherStudent|others
	"ForumView":              0xff, // "11111111",
	"ForumPost":              0xf0, // "11110000",
	"ForumReply":             0xf0, // "11110000",
	"ForumPostRemove":        0xe0, // "11100000",
	"MaterialAdd":            0xc0, // "11000000",
	"MaterialRemove":         0xc0, // "11000000",
	"MaterialDownload":       0xf0, // "11110000",
	"AssignmentAdd":          0xe0, // "11100000",
	"AssignmentRemove":       0xe0, // "11100000",
	"AssignmentUpdate":       0xe0, // "11100000",
	"AssignmentView":         0xf0, // "11110000",
	"AssignmentCheck":        0xc0, // "11000000",
	"AssignmentDo":           0x10, // "00010000",
	"AnnouncementAdd":        0xc0, // "11000000",
	"AnnouncementRemove":     0xc0, // "11000000",
	"AnnouncementUpdate":     0xc0, // "11000000",
	"AnnouncementView":       0xf0, // "11110000",
	"TeachingSyllabusUpdate": 0xc0, // "11000000",
	"IntroductionUpdate":     0xc0, // "11000000",
	"StudentList":            0xf0, // "11110000",
	"MakeTeam":               0xe0, // "11100000",
	"PowersControl":          0xc0, // "11000000",
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

func getGroup(classid string, collection string, id string) int {
	var err error
	switch collection {
	case "admin":
		return 128
	case "student":
		_, err = mgoFind("class", bson.M{"_id": classid, "students.id": id})
		if err != nil {
			log.Println(err)
			return 2
		}
		return 16
	case "teacher":
		_, err = mgoFind("class", bson.M{"_id": classid, "teachers.id": id})
		if err != nil {
			log.Println(err)
			return 8
		}
		return 64
	case "teachingAssistant":
		_, err = mgoFind("class", bson.M{"_id": classid, "teachingassistantid": id})
		if err != nil {
			log.Println(err)
			return 4
		}
		return 32
	default:
		return 1
	}
}

func getClassPowerTable(id string) (map[string]interface{}, error) {
	result, err := mgoFind("class", bson.M{"_id": id})
	if err != nil {
		return nil, err
	}
	switch result["powerTable"].(type) {
	case map[string]interface{}:
	default:
		return nil, errors.New("错误的powerTable类型")
	}
	powerTable := result["powerTable"].(map[string]interface{})
	return powerTable, nil
}

func getCoursePowerTable(id string) (map[string]interface{}, error) {
	result, err := mgoFind("course", bson.M{"_id": id})
	if err != nil {
		return nil, err
	}
	switch result["powerTable"].(type) {
	case map[string]interface{}:
	default:
		return nil, errors.New("错误的powerTable类型")
	}
	powerTable := result["powerTable"].(map[string]interface{})
	return powerTable, nil
}

func getPowerTable(classid string) (map[string]int, error) {
	powerTable := make(map[string]int)
	ClassPowerTable, _ := getClassPowerTable(classid)
	courseid, err := getCourseID(classid)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	CoursePowerTable, _ := getCoursePowerTable(courseid)
	for k, v := range defaultPowerTable {
		if _, ok := ClassPowerTable[k]; ok {
			switch ClassPowerTable[k].(type) {
			case int:
			default:
				return nil, errors.New("错误的powerTable类型")
			}
			powerTable[k] = ClassPowerTable[k].(int)
		} else if _, ok := CoursePowerTable[k]; ok {
			switch CoursePowerTable[k].(type) {
			case int:
			default:
				return nil, errors.New("错误的powerTable类型")
			}
			powerTable[k] = CoursePowerTable[k].(int)
		} else {
			powerTable[k] = v
		}
	}
	return powerTable, nil
}

func getPowers(r *http.Request, classid string) map[string]bool {
	powers := make(map[string]bool)
	group := 1
	if koala.ExistSession(r, "sessionID") {
		session := koala.PeekSession(r, "sessionID")
		group = getGroup(classid, session.Values["collection"].(string), session.Values["id"].(string))
	}
	powerTable, err := getPowerTable(classid)
	if err != nil {
		return nil
	}
	for k, v := range powerTable {
		powers[k] = !(v&group == 0)
	}
	return powers
}

func updatePowerTable(collection string, id string, powerTable map[string]int) error {
	return mgoUpdate(collection,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"powerTable": powerTable}})
}

// admin|teacher|teachingAssistant|student|otherTeacher|otherTeachingAssistant|otherStudent|others
func getPowerGroup(x int, power string) map[string]interface{} {
	result := make(map[string]interface{})
	result["power"] = power
	result["admin"] = ((x>>7)&1 == 1)
	result["teacher"] = ((x>>6)&1 == 1)
	result["teachingAssistant"] = ((x>>5)&1 == 1)
	result["student"] = ((x>>4)&1 == 1)
	result["otherTeacher"] = ((x>>3)&1 == 1)
	result["otherTeachingAssistant"] = ((x>>2)&1 == 1)
	result["otherStudent"] = ((x>>1)&1 == 1)
	result["others"] = (x&1 == 1)
	return result
}

func getPowerTableBool(collection string, id string) (map[string]interface{}, []string, error) {
	var notReferedPowers []string
	for k := range defaultPowerTable {
		notReferedPowers = append(notReferedPowers, k)
	}
	result, err := mgoFind(collection, bson.M{"_id": id})
	if err != nil {
		return nil, notReferedPowers, err
	}
	switch result["powerTable"].(type) {
	case map[string]interface{}:
	default:
		return nil, notReferedPowers, errors.New("错误的powerTable类型")
	}
	powerTable := result["powerTable"].(map[string]interface{})
	ClassPowerTable := make(map[string]interface{})
	for k, v := range powerTable {
		for index, val := range notReferedPowers {
			if k == val {
				notReferedPowers = append(notReferedPowers[:index], notReferedPowers[index+1:]...)
			}
		}
		switch v.(type) {
		case int:
		default:
			return nil, notReferedPowers, errors.New("错误的powerTable类型")
		}
		ClassPowerTable[k] = getPowerGroup(v.(int), k)
	}
	return ClassPowerTable, notReferedPowers, nil
}

func powersHandlers() {
	koala.Get("/class/:id/powers", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		id := p.ParamUrl["id"]
		powers := getPowers(r, id)
		if !powers["PowersControl"] {
			koala.NotFound(w)
			return
		}
		PowerTable, notReferedPowers, _ := getPowerTableBool("class", id)
		koala.Render(w, "powers_class.html", map[string]interface{}{
			"title":            courseWeb,
			"id":               id,
			"admin":            admincheck(w, r),
			"PowerTable":       PowerTable,
			"globalPowers":     globalPowers,
			"notReferedPowers": notReferedPowers,
			"powers":           powers,
		})
	})

	koala.Post("/class/:id/powers/update", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		id := p.ParamUrl["id"]
		powers := getPowers(r, id)
		if !powers["PowersControl"] {
			koala.NotFound(w)
			return
		}
		powerTable := make(map[string]int)
		for _, v := range p.ParamPost["power"] {
			powerGroup := p.ParamPost[v]
			powerValue := 0
			for _, v := range powerGroup {
				val, err := strconv.Atoi(v)
				if err != nil {
					log.Println(err)
					koala.Relocation(w, "/class/"+id+"/powers", "参数错误", "error")
					return
				}
				powerValue += val
			}
			powerTable[v] = powerValue
		}
		err := updatePowerTable("class", id, powerTable)
		if err != nil {
			koala.Relocation(w, "/class/"+id+"/powers", "班级权限更新失败", "error")
			return
		}
		koala.Relocation(w, "/class/"+id+"/powers", "班级权限更新成功", "error")
	})

	koala.Get("/course/:id/powers", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		id := p.ParamUrl["id"]
		if !admincheck(w, r) {
			koala.Relocation(w, "/admin", "请先登录", "error")
			return
		}
		PowerTable, notReferedPowers, _ := getPowerTableBool("course", id)
		koala.Render(w, "powers_course.html", map[string]interface{}{
			"title":            courseWeb,
			"id":               id,
			"PowerTable":       PowerTable,
			"globalPowers":     globalPowers,
			"notReferedPowers": notReferedPowers,
		})
	})

	koala.Post("/course/:id/powers/update", func(p *koala.Params, w http.ResponseWriter, r *http.Request) {
		id := p.ParamUrl["id"]
		if !admincheck(w, r) {
			koala.Relocation(w, "/admin", "请先登录", "error")
			return
		}
		powerTable := make(map[string]int)
		for _, v := range p.ParamPost["power"] {
			powerGroup := p.ParamPost[v]
			powerValue := 0
			for _, v := range powerGroup {
				val, err := strconv.Atoi(v)
				if err != nil {
					log.Println(err)
					koala.Relocation(w, "/course/"+id+"/powers", "参数错误", "error")
					return
				}
				powerValue += val
			}
			powerTable[v] = powerValue
		}
		err := updatePowerTable("course", id, powerTable)
		if err != nil {
			koala.Relocation(w, "/course/"+id+"/powers", "课程权限更新失败", "error")
			return
		}
		koala.Relocation(w, "/course/"+id+"/powers", "课程权限更新成功", "error")
	})
}
