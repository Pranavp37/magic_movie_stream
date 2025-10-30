package routes

import (
	"github.com/Pranavp37/magic_movie_stream/internal/handlers"
	"github.com/Pranavp37/magic_movie_stream/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterHandler(r *gin.Engine) {

	v1 := r.Group("api/v1")
	{
		v1.GET("/hello-dear", handlers.HelloHandler)
		v1.POST("/auth/register", handlers.RegisterUser)
		v1.POST("/auth/login", handlers.Login)
	}

	v2 := r.Group("api/v2", middleware.JwtMiddleware())
	{
		v2.GET("/user/details", handlers.GetUserDetails)
		v2.GET("/user/users/search", handlers.SearchUser)
	}

}
