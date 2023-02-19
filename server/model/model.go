package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Question struct {
	Id       primitive.ObjectID `bson:"_id" json:"_id"`
	Question string             `bson:"question" json:"question"`
	Answers  []Answer           `bson:"answers" json:"answers"`
	Category string             `bson:"category" json:"category"`
}

type Answer struct {
	Value   string `bson:"value" json:"value"`
	Correct bool   `bson:"correct" json:"correct"`
}

type Game struct {
	Id            primitive.ObjectID   `bson:"_id" json:"_id"`
	Name          string               `bson:"name" json:"name"`
	Teams         []Team               `bson:"teams" json:"teams"`
	QuestionsUsed []primitive.ObjectID `bson:"questionsUsed" json:"questionsUsed"`
}

type Team struct {
	Name  string `bson:"name" json:"name"`
	Score int    `bson:"score" json:"score"`
}

type OpenDBResponse struct {
	ResponseCode int              `bson:"response_code"`
	Results      []OpenDBQuestion `bson:"results"`
}

type OpenDBQuestion struct {
	Category         string   `json:"category"`
	Question         string   `json:"question"`
	CorrectAnswer    string   `json:"correct_answer"`
	IncorrectAnswers []string `json:"incorrect_answers"`
}
