package route

import (
	"net/http"
	"os"
	"trivia-the-game/server/ctx"
	"trivia-the-game/server/model"
	"trivia-the-game/server/server"

	"github.com/gin-gonic/gin"
)

type Router struct {
	serviceContext ctx.ServiceContext
	server         server.DefaultServer
}

func NewRouter(serviceContext ctx.ServiceContext) *Router {
	return &Router{
		serviceContext: serviceContext,
		server:         *server.NewServer(serviceContext),
	}
}

func (r *Router) Install(engine *gin.RouterGroup) {
	engine.POST("/question", r.CreateQuestion)
	engine.GET("/question", r.GetQuestion)
	engine.GET("/questions-opendb", r.GetQuestionsFromOpenDB)
	engine.POST("/game", r.CreateNewGame)
	engine.POST("/update-score", r.UpdateScore)
}

func (r *Router) CreateQuestion(ginctx *gin.Context) {
	req := model.Question{}
	bindErr := ginctx.BindJSON(&req)
	if bindErr != nil {
		ginctx.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
		return
	}
	err := r.server.CreateQuestion(req)
	if err != nil {
		ginctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ginctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (r *Router) GetQuestion(ginctx *gin.Context) {
	category, _ := ginctx.GetQuery("category")
	gameId, _ := ginctx.GetQuery("gameId")
	if gameId == "" {
		ginctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid_game_id"})
		return
	}
	res, err := r.server.GetQuestion(category, gameId)
	if err != nil {
		ginctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ginctx.JSON(http.StatusOK, gin.H{"data": res})
}

func (r *Router) GetQuestionsFromOpenDB(ginctx *gin.Context) {
	allowGetQuestionsFromOpenDB := os.Getenv("ALLOW_GET_QUESTIONS_FROM_OPENDB") == "true"
	if !allowGetQuestionsFromOpenDB {
		ginctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid_allow_get_questions_from_openDB"})
		return
	}
	err := r.server.GetQuestionsFromOpenDB()
	if err != nil {
		ginctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ginctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (r *Router) CreateNewGame(ginctx *gin.Context) {
	req := model.Game{}
	bindErr := ginctx.BindJSON(&req)
	if bindErr != nil {
		ginctx.JSON(http.StatusBadRequest, gin.H{"error": "bad_request"})
		return
	}
	err := r.server.CreateNewGame(req)
	if err != nil {
		ginctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ginctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (r *Router) GetGames(ginctx *gin.Context) {
	res, err := r.server.GetGames()
	if err != nil {
		ginctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ginctx.JSON(http.StatusOK, gin.H{"data": res})
}

func (r *Router) UpdateScore(ginctx *gin.Context) {
	req := model.UpdateScoreRequest{}
	bindErr := ginctx.BindJSON(&req)
	if bindErr != nil {
		ginctx.JSON(http.StatusBadRequest, gin.H{"error": "bad_request"})
		return
	}
	res, err := r.server.UpdateScore(req)
	if err != nil {
		ginctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ginctx.JSON(http.StatusOK, gin.H{"data": res})
}
