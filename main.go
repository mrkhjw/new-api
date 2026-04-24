package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"new-api/common"
	"new-api/middleware"
	"new-api/model"
	"new-api/router"
)

func main() {
	common.SetupLogger()
	common.SysLog("New API " + common.Version + " started")

	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	if os.Getenv("DEBUG") == "true" {
		common.DebugEnabled = true
	}

	// Initialize database
	err := model.InitDB()
	if err != nil {
		common.FatalLog("failed to initialize database: " + err.Error())
	}
	defer func() {
		err := model.CloseDB()
		if err != nil {
			common.SysError("failed to close database: " + err.Error())
		}
	}()

	// Initialize Redis (optional)
	err = common.InitRedisClient()
	if err != nil {
		// Redis is optional; log the error but continue startup
		common.SysLog("Redis not available, running without cache: " + err.Error())
	}

	// Initialize options from database
	model.InitOptionMap()

	// Initialize token encoder
	common.InitTokenEncoders()

	// Setup Gin router
	server := gin.New()
	server.Use(gin.Recovery())
	server.Use(middleware.RequestId())
	middleware.SetUpLogger(server)

	// Register all routes
	router.SetRouter(server)

	// Default to port 3000 instead of upstream default; override with PORT env var
	var port = os.Getenv("PORT")
	if port == "" {
		port = strconv.Itoa(*common.Port)
		if port == "0" {
			port = "3000"
		}
	}

	common.SysLog(fmt.Sprintf("server started on http://localhost:%s", port))

	if err := server.Run(":" + port); err != nil {
		common.FatalLog("failed to start HTTP server: " + err.Error())
	}
}
