package server

import (
	"context"
	"errors"
	"fmt"
	"trivia-the-game/server/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *DefaultServer) CreateNewGame(game model.Game) error {
	op := "NewGame"
	db := s.serviceContext.GetDB()
	col := db.Collection(gamesCol)
	game.Id = primitive.NewObjectID()
	game.QuestionsUsed = []primitive.ObjectID{}
	for _, team := range game.Teams {
		team.Score = 0
	}
	_, err := col.InsertOne(context.Background(), game)
	if err != nil {
		s.serviceContext.Logger().Errorln(op+".InsertOne", err.Error())
		return err
	}
	return nil
}

func (s *DefaultServer) GetGames() ([]model.Game, error) {
	op := "GetGames"
	db := s.serviceContext.GetDB()
	col := db.Collection(gamesCol)
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"_id": -1})
	cur, err := col.Find(context.Background(), bson.M{}, findOptions)
	if err != nil {
		s.serviceContext.Logger().Errorln(op+".Find", err.Error())
		return []model.Game{}, err
	}
	games := []model.Game{}
	for cur.Next(context.Background()) {
		game := model.Game{}
		err := cur.Decode(&game)
		if err != nil {
			s.serviceContext.Logger().Errorln(op+".Decode", err.Error())
			return []model.Game{}, err
		}
		games = append(games, game)
	}
	return games, nil
}

func (s *DefaultServer) GetGameById(id primitive.ObjectID) (model.Game, error) {
	op := "GetGameById"
	game := model.Game{}
	db := s.serviceContext.GetDB()
	col := db.Collection(gamesCol)
	err := col.FindOne(context.Background(), bson.M{"_id": id}).Decode(&game)
	if err != nil {
		s.serviceContext.Logger().Errorln(op+".FindOne", err.Error())
		return model.Game{}, err
	}

	return game, nil
}

func (s *DefaultServer) ExcludeQuestionFromGame(gameId, questionId string) error {
	op := "ExcludeQuestionFromGame"
	gameObjId, err := primitive.ObjectIDFromHex(gameId)
	if err != nil {
		s.serviceContext.Logger().Errorln(op+".ObjectIDFromHex.gameId", err.Error())
	}
	questionObjId, err := primitive.ObjectIDFromHex(questionId)
	if err != nil {
		s.serviceContext.Logger().Errorln(op+".ObjectIDFromHex.questionId", err.Error())
	}
	db := s.serviceContext.GetDB()
	col := db.Collection(gamesCol)
	upsert := true
	filter := bson.M{"_id": gameObjId}
	update := bson.M{"$addToSet": bson.M{"questionsUsed": questionObjId}}
	options := &options.UpdateOptions{Upsert: &upsert}
	_, err = col.UpdateOne(context.TODO(), filter, update, options)
	if err != nil {
		s.serviceContext.Logger().Errorln(op+".UpdateOne", err.Error())
	}
	s.serviceContext.Logger().Debugln("ExcludeQuestionFromGame.GameID_", gameId, "_QuestionID_", questionId)
	return nil
}

func (s *DefaultServer) UpdateScore(req model.UpdateScoreRequest) (model.Game, error) {
	op := "UpdateScore"
	gameObjId, err := primitive.ObjectIDFromHex(req.GameId)
	if err != nil {
		s.serviceContext.Logger().Errorln(op+".ObjectIDFromHex.gameId", err.Error())
		return model.Game{}, err
	}
	db := s.serviceContext.GetDB()
	col := db.Collection(gamesCol)
	filter := bson.M{"_id": gameObjId}
	game := model.Game{}
	err = col.FindOne(context.Background(), filter).Decode(&game)
	if err != nil {
		s.serviceContext.Logger().Errorln(op+".Decode", err.Error())
		return model.Game{}, err
	}
	for i, team := range game.Teams {
		if team.Name == req.TeamName {
			if req.Operation == "subtract" {
				if team.Score == 0 {
					s.serviceContext.Logger().Errorln(op + ".subtractZero")
					return model.Game{}, errors.New("subtract_zero")
				}
				game.Teams[i].Score = team.Score - 1
			} else if req.Operation == "addition" {
				game.Teams[i].Score = team.Score + 1
			}
			break
		}
	}
	fmt.Println("new game:", game)
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	updatedGame := model.Game{}
	update := bson.M{"$set": game}
	err = col.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&updatedGame)
	if err != nil {
		s.serviceContext.Logger().Errorln(op+".FindOneAndUpdate", err.Error())
		return model.Game{}, err
	}
	return updatedGame, nil
}
