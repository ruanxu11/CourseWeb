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

func mgoFindAll(collection string, selector map[string]interface{}) (searchResults []map[string]interface{}, err error) {
	query := func(c *mgo.Collection) (*mgo.ChangeInfo, error) {
		err := c.Find(selector).All(&searchResults)
		return nil, err
	}
	_, err = withCollection(collection, query)
	return searchResults, err
}

func mgoFind(collection string, selector map[string]interface{}) (searchResults map[string]interface{}, err error) {
	query := func(c *mgo.Collection) (*mgo.ChangeInfo, error) {
		err := c.Find(selector).One(&searchResults)
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

func mgoRemove(collection string, selector interface{}) (info *mgo.ChangeInfo, err error) {
	query := func(c *mgo.Collection) (info *mgo.ChangeInfo, err error) {
		err = c.Remove(selector)
		return nil, err
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

func mgoUpdate(collection string, selector interface{}, update interface{}) error {
	query := func(c *mgo.Collection) (*mgo.ChangeInfo, error) {
		err := c.Update(selector, update)
		return nil, err
	}
	_, err := withCollection(collection, query)
	return err
}
