package main

type TeachingAssistant struct {
	ID           int `bson:"_id"`
	Password     string
	Name         string
	Introduction string
}
