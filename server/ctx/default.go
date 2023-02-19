package ctx

import (
	"context"
	"time"

	"os"
	"strconv"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ServiceContext interface {
	GetMongoClient() *mongo.Client
	GetEnvValue(string, bool) string
	GetIntEnvValue(string, bool) int
	GetDB() *mongo.Database
	Logger() *logrus.Logger
}

type DefaultServiceContext struct {
	// connection pools, integration clients, environment, etc. should be here
	name          string
	mongodbClient *mongo.Client
	dbo           *mongo.Database
	logger        *logrus.Logger
}

func NewDefaultServiceContext() *DefaultServiceContext {

	mylogger := logrus.New()
	mylogger.SetLevel(logrus.DebugLevel)
	// Force a connection to verify our connection string
	ctx := &DefaultServiceContext{
		logger: mylogger,
	}
	return ctx
}
func (ctx *DefaultServiceContext) WithMongo() *DefaultServiceContext {
	mongodbClient := ctx.getMongoClient(os.Getenv("MONGO_CONNECTION_URL"))
	ctx.mongodbClient = mongodbClient
	ctx.dbo = mongodbClient.Database(os.Getenv("MONGO_DB_NAME"))
	return ctx
}
func (ctx *DefaultServiceContext) getMongoClient(uri string) *mongo.Client {
	maxPools, err := strconv.ParseUint(os.Getenv("MAX_POOL_SIZE"), 10, 64)
	minPools, err := strconv.ParseUint(os.Getenv("MIN_POOL_SIZE"), 10, 64)
	idleTime, err := strconv.ParseUint(os.Getenv("MAX_IDLE_SECONDS"), 10, 64)
	mongodbClient, err := mongo.NewClient(
		options.Client().ApplyURI(uri),
		options.Client().SetMaxPoolSize(maxPools),
		options.Client().SetMinPoolSize(minPools),
		options.Client().SetMaxConnIdleTime(time.Duration(idleTime)*time.Second),
	)
	err = mongodbClient.Connect(context.Background())
	if err != nil {
		ctx.Logger().Panic("No mongo connection :", err)
	}
	err = mongodbClient.Ping(context.Background(), nil)
	if err != nil {
		ctx.Logger().Panic("failed to ping mongo :", err)
	}
	return mongodbClient
}
func (ctx *DefaultServiceContext) GetMongoClient() *mongo.Client {
	if ctx.mongodbClient == nil {
		ctx.mongodbClient = ctx.WithMongo().mongodbClient
	}
	return ctx.mongodbClient
}

func (ctx *DefaultServiceContext) GetDB() *mongo.Database {
	if ctx.dbo == nil {
		ctx.dbo = ctx.WithMongo().dbo
	}
	return ctx.dbo
}

func (ctx *DefaultServiceContext) Logger() *logrus.Logger {
	return ctx.logger
}
func (ctx *DefaultServiceContext) Shutdown() {
	if ctx.mongodbClient != nil {
		ctx.mongodbClient.Disconnect(context.Background())
	}
}

func (ctx *DefaultServiceContext) GetEnvValue(envName string, isENVRequired bool) (envValue string) {
	envValue = os.Getenv(envName)
	if envValue == "" && isENVRequired {
		panic("Missing required env var: " + envName)
	}
	return
}

func (ctx *DefaultServiceContext) GetIntEnvValue(envName string, isENVRequired bool) (intEnvValue int) {
	envValue := os.Getenv(envName)
	if envValue == "" && isENVRequired {
		ctx.Logger().Panic("Missing required env var: " + envName)
	}
	intEnvValue, err := strconv.Atoi(envValue)
	if err != nil {
		ctx.Logger().Panicf("Error converting env var: %s to Int, %s", envName, err.Error())
	}
	return intEnvValue
}

// func (ctx *DefaultServiceContext) WithRedis() *DefaultServiceContext {
// 	redisClient := redis.NewClient(&redis.Options{
// 		Addr:         os.Getenv("REDIS_URL"),
// 		ReadTimeout:  time.Minute,
// 		DialTimeout:  time.Minute,
// 		WriteTimeout: time.Minute,
// 	})
// 	ctx.redisClient = redisClient
// 	return ctx
// }

// func (ctx *DefaultServiceContext) WithMysql() *DefaultServiceContext {
// 	db, err := sql.Open("mysql", os.Getenv("MYSQL_URL"))
// 	if err != nil {
// 		ctx.logger.Panic(err)
// 	}

// 	ctx.mysql = db
// 	return ctx
// }
// func (ctx *DefaultServiceContext) GetMysql() *sql.DB {
// 	if ctx.mysql == nil {
// 		ctx.mysql = ctx.WithMysql().mysql
// 	}
// 	return ctx.mysql
// }
