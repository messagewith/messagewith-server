package main

import (
	"awesomeProject/graph"
	"awesomeProject/graph/generated"
	"awesomeProject/utils"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
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

func initDatabaseConnection() gin.HandlerFunc {
	err := mgm.SetDefaultConfig(nil, "awesome_project", options.Client().ApplyURI(os.Getenv("MESSAGEWITH_DATABASE_URI")))

	if err != nil {
		panic(err)
	}

	return func(c *gin.Context) {}
}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		panic(fmt.Errorf("create .env file"))
	}

	r := gin.Default()
	r.Use(initDatabaseConnection())
	r.Use(utils.GinContextToContextMiddleware())
	r.POST("/query", graphqlHandler())
	r.GET("/", playgroundHandler())
	r.Run(":8000")
}
