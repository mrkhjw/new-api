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

	// Default to port 3001; override with PORT env var or common.Port config.
	// Personal note: using 3001 to avoid conflicts with common dev servers
	// (3000 is often React/Next.js, 8080 is often other backend services).
	// If common.Port is explicitly set to something other than 0, respect that value.
	// Fallback chain: PORT env var -> common.Port config -> default (3001)
	var port = os.Getenv("PORT")
	if port == "" {
		configPort := strconv.Itoa(*common.Port)
		if configPort != "0" {
			port = configPort
		} else {
			port = "3001"
		}
	}

	common.SysLog(fmt.Sprintf("server started on http://localhost:%s", port))
	fmt.Printf("\n  ➜  Local: http://localhost:%s\n\n", port)

	if err := server.Run(":" + port); err != nil {
		common.FatalLog("failed to start HTTP server: " + err.Error())
	}
}
