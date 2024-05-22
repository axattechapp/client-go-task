package router

import (
	"client_task/graph/generated"
	graph "client_task/graph/resolvers"
	sqlc "client_task/pkg/common/db/sqlc"
	"client_task/pkg/controllers"
	webhook "client_task/pkg/controllers/webhooks"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func SetupRouter(db *sqlc.Queries, ctx context.Context) *gin.Engine {
	server := gin.Default()

	UsersController := controllers.NewUsersController(db, ctx)
	UserRoutes := NewRouteUser(*UsersController)
	UserAuthRoutes := NewRouteUserAuth(*UsersController)

	ProfileController := controllers.NewProfilesController(db, ctx)
	ProfileRoutes := NewProfileRoutes(*ProfileController)

	CareerController := controllers.NewCareersController(db, ctx)
	CareerRoutes := NewCareerRoutes(*CareerController)

	SkillsController := controllers.NewSkillsController(db, ctx)
	SkillsRoutes := NewSkillsRoutes(*SkillsController)

	WebhookController := webhook.NewWebhookController(db, ctx)
	WebhookRoutes := NewWebhookRoutes(WebhookController)

	router := server.Group("/api")

	router.GET("/healthcheck", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "The contact API is working fine"})
	})

	UserRoutes.UserRoutes(router)
	UserAuthRoutes.UserAuthRoutes(router)
	ProfileRoutes.ProfileRoutes(router)
	CareerRoutes.CareerRoutes(router)
	SkillsRoutes.SkillsRoutes(router)
	WebhookRoutes.RegisterRoutes(router)

	server.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": fmt.Sprintf("The specified route %s not found", ctx.Request.URL)})
	})

	router2 := server.Group("/graphql")

	server.POST("/query", graphqlHandler(db))
	router2.GET("/", playgroundHandler())

	return server
}

func graphqlHandler(db *sqlc.Queries) gin.HandlerFunc {
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{DB: db}}))
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
