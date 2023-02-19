package server

import (
	"context"
	"encoding/json"
	"net/http"
	"trivia-the-game/server/model"
	"trivia-the-game/server/utilities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *DefaultServer) GetQuestion(category string, gameId string) (model.Question, error) {
	op := "GetQuestion"
	excludedIds := []primitive.ObjectID{}
	if gameId != "" {
		gameObjId, err := primitive.ObjectIDFromHex(gameId)
		if err != nil {
			return model.Question{}, err
		}
		game, err := s.GetGameById(gameObjId)
		if err != nil {
			return model.Question{}, err
		}
		excludedIds = append(excludedIds, game.QuestionsUsed...)
	}
	db := s.serviceContext.GetDB()
	col := db.Collection(questionsCol)
	pipeline := []bson.D{}

	if category != "" {
		categoryFilter := bson.D{{
			"$match",
			bson.M{"category": category},
		}}
		pipeline = append(pipeline, categoryFilter)
	}
	if gameId != "" {
		excludedIdsFilter := bson.D{{
			"$match",
			bson.M{
				"_id": bson.M{"$nin": excludedIds},
			},
		},
		}
		pipeline = append(pipeline, excludedIdsFilter)
	}
	pipeline = append(pipeline, bson.D{{"$sample", bson.D{{"size", 1}}}})

	cur, err := col.Aggregate(context.Background(), pipeline)
	if err != nil {
		s.serviceContext.Logger().Errorln(op+".Aggregate", err.Error())
		return model.Question{}, err
	}
	question := model.Question{}
	for cur.Next(context.Background()) {
		err := cur.Decode(&question)
		if err != nil {
			s.serviceContext.Logger().Errorln(op+".Decode", err.Error())
			return model.Question{}, err
		}
		break
	}
	go s.ExcludeQuestionFromGame(gameId, question.Id.Hex())
	return question, nil
}

func (s *DefaultServer) CreateQuestion(question model.Question) error {
	op := "CreateQuestion"
	db := s.serviceContext.GetDB()
	col := db.Collection(questionsCol)
	question.Id = primitive.NewObjectID()
	_, err := col.InsertOne(context.Background(), question)
	if err != nil {
		s.serviceContext.Logger().Errorln(op+".InsertOne", err.Error())
		return err
	}
	return nil
}

func (s *DefaultServer) GetQuestionsFromOpenDB() error {
	op := "GetQuestionsFromOpenDB"
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, "https://opentdb.com/api.php?amount=50&type=multiple", nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		s.serviceContext.Logger().Errorln(op+".clientDo", err.Error())
		return err
	}
	var resJson model.OpenDBResponse
	json.NewDecoder(res.Body).Decode(&resJson)

	questions := []model.Question{}
	for i := range resJson.Results {
		opendbQuestion := resJson.Results[i]
		answers := []model.Answer{}
		answers = append(answers, model.Answer{
			Value:   opendbQuestion.CorrectAnswer,
			Correct: true,
		})
		for _, inCorrectAnswer := range opendbQuestion.IncorrectAnswers {
			answers = append(answers, model.Answer{Value: inCorrectAnswer, Correct: false})
		}
		answersShuffled := utilities.ShuffleAnswers(answers)

		question := model.Question{
			Id:       primitive.NewObjectID(),
			Question: opendbQuestion.Question,
			Category: opendbQuestion.Category,
			Answers:  answersShuffled,
		}
		questions = append(questions, question)
	}

	db := s.serviceContext.GetDB()
	col := db.Collection(questionsCol)
	inteface := []interface{}{}
	for _, q := range questions {
		inteface = append(inteface, q)
	}
	_, err = col.InsertMany(context.Background(), inteface)
	if err != nil {
		s.serviceContext.Logger().Errorln(op+".InsertMany", err.Error())
		return err
	}

	return nil
}
