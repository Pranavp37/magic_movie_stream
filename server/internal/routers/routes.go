package routes

import (
	"github.com/Pranavp37/magic_movie_stream/internal/handlers"
	"github.com/Pranavp37/magic_movie_stream/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterHandler(r *gin.Engine) {

	// for texting porpuse
	r.GET("/", handlers.HelloHandler)

	v1 := r.Group("api/v1")
	{

		v1.POST("/auth/register", handlers.RegisterUser)
		v1.POST("/auth/login", handlers.Login)
		v1.GET("/ws", handlers.WebSocketConnection)
		v1.GET("/sse", handlers.ChatListSSE)
	}

	v2 := r.Group("api/v2", middleware.JwtMiddleware())
	{
		v2.GET("/user/details", handlers.GetUserDetails)
		v2.GET("/user/users/search", handlers.SearchUser)
		v2.GET("/chat/chatList", handlers.ChatListHistory)
	}

}
