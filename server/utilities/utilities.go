package utilities

import (
	"math/rand"
	"trivia-the-game/server/model"
)

func ShuffleAnswers(slice []model.Answer) []model.Answer {
	for i := range slice {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}
