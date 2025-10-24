package servers

import (
	"fmt"
	"log"

	"github.com/Pranavp37/magic_movie_stream/internal/config"
	"github.com/Pranavp37/magic_movie_stream/internal/database"
	routes "github.com/Pranavp37/magic_movie_stream/internal/routers"
	"github.com/gin-gonic/gin"
)

func Start() {
	// properties from config
	configdatas := config.LoadConfig()

	//mongo connection
	database.MongoDBConnection(configdatas.MongoDB_url)

	// gin initialize
	r := gin.Default()

	// gin router register
	routes.RegisterHandler(r)

	addr := fmt.Sprintf("0.0.0.0:%s", configdatas.PORT)
	err := r.Run(addr)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

}
