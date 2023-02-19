package server

import (
	"trivia-the-game/server/ctx"
)

const (
	questionsCol = "questions"
	gamesCol     = "games"
)

type DefaultServer struct {
	serviceContext ctx.ServiceContext
}

func NewServer(serviceContext ctx.ServiceContext) *DefaultServer {
	return &DefaultServer{
		serviceContext: serviceContext,
	}
}
