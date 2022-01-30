package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
	"messagewith-server/auth"
	"messagewith-server/env"
	"messagewith-server/graph"
	"messagewith-server/graph/generated"
	"messagewith-server/mails"
	"messagewith-server/middlewares"
	"messagewith-server/sessions"
	"messagewith-server/users"
)

func graphqlHandler() gin.HandlerFunc {
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func initDatabaseConnection() {
	err := mgm.SetDefaultConfig(nil, "messagewith", options.Client().ApplyURI(env.DatabaseURI))
	if err != nil {
		panic(err)
	}
}

func main() {
	env.InitEnvConstants()
	initDatabaseConnection()
	mails.InitClient()
	users.InitService()
	sessions.InitService()
	auth.InitService()

	r := gin.Default()
	r.Use(middlewares.GinContextToContextMiddleware())
	r.Use(middlewares.AuthMiddleware())
	r.POST("/query", graphqlHandler())
	r.GET("/", playgroundHandler())
	r.Run(":8000")
}
