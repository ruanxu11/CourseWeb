package main

import (
	"errors"
	"strconv"

	"log"

	"labix.org/v2/mgo/bson"
)

type SecurityQuestion struct {
	Question string
	Answer   string
}

func loginCheck(collection string, sid string, password string) (string, bool) {
	id, _ := strconv.Atoi(sid)
	person, err := mgoFind(collection,
		bson.M{"_id": id, "password": password})
	if err != nil {
		return "", false
	}
	if person, ok := person["name"].(string); ok {
		return person, true
	}
	return "", true
}

func changePWD(collection string, sid string, newPassword string, oldPassword string) error {
	id, _ := strconv.Atoi(sid)
	if _, ok := loginCheck(collection, sid, oldPassword); !ok {
		return errors.New("账号或密码错误")
	}
	return mgoUpdate(collection,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"password": newPassword}})
}

func showSecurityQuestions(collection string, sid string) ([]interface{}, error) {
	id, _ := strconv.Atoi(sid)
	person, err := mgoFind(collection,
		bson.M{"_id": id})
	securityQuestions := person["securityQuestions"].([]interface{})
	return securityQuestions, err
}

func forgetPWD(collection string, sid string, newPassword string, securityQuestions []SecurityQuestion) error {
	id, _ := strconv.Atoi(sid)
	person, _ := mgoFind(collection,
		bson.M{"_id": id, "securityQuestions": securityQuestions})
	switch person["securityQuestions"].(type) {
	case []interface{}:
	default:
		return errors.New("安全问题回答错误")
	}
	return mgoUpdate(collection,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"password": newPassword}})
}

func changeSecurityQuestionsbyOld(collection string, sid string, newSecurityQuestions []SecurityQuestion, oldSecurityQuestions []SecurityQuestion) error {
	id, _ := strconv.Atoi(sid)
	log.Println(bson.M{"_id": id, "securityQuestions": oldSecurityQuestions})
	person, _ := mgoFind(collection,
		bson.M{"_id": id, "securityQuestions": oldSecurityQuestions})
	switch person["securityQuestions"].(type) {
	case []interface{}:
	default:
		return errors.New("安全问题回答错误")
	}
	return mgoUpdate(collection,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"securityQuestions": newSecurityQuestions}})
}
