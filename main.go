package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
	"messagewith-server/auth"
	"messagewith-server/chats"
	"messagewith-server/env"
	"messagewith-server/graph"
	"messagewith-server/graph/generated"
	"messagewith-server/graph/model"
	"messagewith-server/mails"
	"messagewith-server/middlewares"
	"messagewith-server/sessions"
	"messagewith-server/users"
	"net/http"
	"time"
)

func getGraphqlHandler() gin.HandlerFunc {
	srv := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		ChatMessages:  []*model.Chat{},
		ChatObservers: map[string]chan []*model.Chat{},
	}}))

	srv.AddTransport(&transport.POST{})
	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		KeepAlivePingInterval: 10 * time.Second,
	})
	srv.SetQueryCache(lru.New(1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})

	return func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	srv := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
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
	mails.InitService()
	users.InitService()
	sessions.InitService()
	auth.InitService()
	chats.InitService()

	r := gin.Default()
	r.Use(middlewares.GinContextToContextMiddleware())
	r.Use(middlewares.AuthMiddleware())

	graphqlHandler := getGraphqlHandler()

	r.POST("/query", graphqlHandler)
	r.GET("/query", graphqlHandler)
	r.GET("/", playgroundHandler())
	r.Run(":8000")
}
