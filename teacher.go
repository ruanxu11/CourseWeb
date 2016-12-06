package main

type Teacher struct {
	ID                         int `bson:"_id"`
	Password                   string
	Name                       string
	Sex                        string
	Introduction               string
	Email                      string
	Phone                      int
	College                    string // 所在学院
	Department                 string // 系
	AcademicBackground         string // 学历
	AcademicTitle              string // 职称
	ResearchDirections         string // 研究方向
	StudentsCenteredEvaluation string // 教学质量评价
	SecurityQuestions          []SecurityQuestion
}

type SecurityQuestion struct {
	Question string
	Answer   string
}
