package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	_ "github.com/deshortone/ledger-system/docs"
	"github.com/deshortone/ledger-system/internal/api"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"
)

//	@title			Equipment Picker Calculator API
//	@version		0.1
//	@description	This will take in user requests for equipments they want

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api

func main() {
	ctx := context.Background()

	router := gin.Default()
	app, err := api.RegisterRoutes(ctx, router)
	if err != nil {
		panic(fmt.Sprintf("Failed to start due to: %s", err.Error()))
	}
	defer app.Close()

	router.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler))

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadTimeout:       25 * time.Second,
		WriteTimeout:      25 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	srv.ListenAndServe()
}
