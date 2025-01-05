package api

import (
	"context"
	"log"
	"net/http"
	"polling_websocket/pkg/config"
	"polling_websocket/pkg/di"
	"polling_websocket/pkg/honeycomb"
	"polling_websocket/pkg/interfaces/middlewares"
	"polling_websocket/pkg/interfaces/routes"

	"github.com/gin-gonic/gin"
)

var (
	app *gin.Engine
)

// Init initializes the application without starting the server.
func init() {
	InitApp()
}

// InitApp initializes the Gin application.
func InitApp() {
	log.Print("---- Initializing App ----")
	config.LoadEnvs(".")

	// Setup OpenTelemetry
	ctx := context.Background()
	tp, exp, err := honeycomb.SetupHoneyComb(ctx)
	if err != nil {
		log.Panicf("ERROR | Failed to initialize OpenTelemetry: %v", err)
	}

	// Ensure sub processes and telemetry are exported correctly.
	defer func() {
		_ = exp.Shutdown(ctx)
		_ = tp.Shutdown(ctx)
	}()

	// Initialize Gin app
	gin.SetMode(gin.DebugMode)
	app = gin.New()

	// Dependency injection and routes setup
	dependencies := di.InitDependencies()
	middlewares.Register(app, dependencies.AuthService)
	routes.Register(app, dependencies)
}

// Handler is the main function that Vercel calls to handle HTTP requests.
func Handler(w http.ResponseWriter, r *http.Request) {
	// If app is not initialized, initialize it
	if app == nil {
		InitApp()
	}
	// Use Gin to serve the HTTP request
	app.ServeHTTP(w, r)
}

func Dummy() {
	RunWebserver()
}

func RunWebserver() {
	addr := config.GetEnv("POLLING_ADDR", ":4050")
	err := app.Run(addr)
	if err != nil {
		log.Panicf("ERROR | Starting gin failed, %v", err)
	}
}
