package main

import "labix.org/v2/mgo"

var mgoSession *mgo.Session

func getSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		mongoDBDialInfo := &mgo.DialInfo{
			Addrs:    []string{"115.159.151.237:27017"},
			Username: "CourseWeb",
			Password: "ruanxu1123",
			Database: "CourseWeb",
			Source:   "CourseWeb",
		}
		mgoSession, err = mgo.DialWithInfo(mongoDBDialInfo) //连接数据库
		if err != nil {
			panic(err) // no, not really
		}
	}
	return mgoSession.Clone()
}

func withCollection(collection string, query func(*mgo.Collection) (*mgo.ChangeInfo, error)) (*mgo.ChangeInfo, error) {
	session := getSession()
	defer session.Close()
	c := session.DB("CourseWeb").C(collection)
	return query(c)
}

func mgoFindAll(collection string, selector map[string]interface{}, skip int, limit int) (searchResults []map[string]interface{}, err error) {

	query := func(c *mgo.Collection) (*mgo.ChangeInfo, error) {
		var err error
		if limit < 0 {
			err = c.Find(selector).Skip(skip).All(&searchResults)
		} else {
			err = c.Find(selector).Skip(skip).Limit(limit).All(&searchResults)
		}
		return nil, err
	}
	_, err = withCollection(collection, query)
	return searchResults, err
}

func mgoFind(collection string, selector map[string]interface{}, skip int, limit int) (searchResults map[string]interface{}, err error) {

	query := func(c *mgo.Collection) (*mgo.ChangeInfo, error) {
		var err error
		if limit < 0 {
			err = c.Find(selector).Skip(skip).One(&searchResults)
		} else {
			err = c.Find(selector).Skip(skip).Limit(limit).One(&searchResults)
		}
		return nil, err
	}
	_, err = withCollection(collection, query)
	return searchResults, err
}

func mgoInsert(collection string, docs interface{}) error {
	query := func(c *mgo.Collection) (*mgo.ChangeInfo, error) {
		err := c.Insert(docs)
		return nil, err
	}
	_, err := withCollection(collection, query)
	return err
}

func mgoRemoveAll(collection string, selector interface{}) (info *mgo.ChangeInfo, err error) {
	query := func(c *mgo.Collection) (info *mgo.ChangeInfo, err error) {
		info, err = c.RemoveAll(selector)
		return info, err
	}
	return withCollection(collection, query)
}

func mgoUpdateAll(collection string, selector interface{}, update interface{}) (info *mgo.ChangeInfo, err error) {
	query := func(c *mgo.Collection) (info *mgo.ChangeInfo, err error) {
		info, err = c.UpdateAll(selector, update)
		return info, err
	}
	return withCollection(collection, query)
}
