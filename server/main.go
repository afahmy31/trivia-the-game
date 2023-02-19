package main

import (
	"os"
	"trivia-the-game/server/ctx"
	"trivia-the-game/server/route"

	"github.com/gin-gonic/gin"
)

func main() {
	serviceContext := ctx.NewDefaultServiceContext().WithMongo()
	httpServer(serviceContext)
}

func httpServer(c *ctx.DefaultServiceContext) {
	engine := gin.Default()
	router := route.NewRouter(c)
	routeGroup := engine.Group("")
	router.Install(routeGroup)
	c.Logger().Info("Http server")
	engine.Run(":" + os.Getenv("PORT"))
}

/**
ENV VARS :
export PORT=8000
export MONGO_DB_NAME=trivia-the-game
export MONGO_CONNECTION_URL="mongodb://localhost:27017"
export ALLOW_GET_QUESTIONS_FROM_OPENDB="false"
**/
