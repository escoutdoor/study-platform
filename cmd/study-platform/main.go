package main

import (
	"context"

	_ "github.com/escoutdoor/study-platform/docs"
	"github.com/escoutdoor/study-platform/internal/app"
	"github.com/escoutdoor/study-platform/internal/config"
	"github.com/escoutdoor/study-platform/pkg/closer"
	"github.com/escoutdoor/study-platform/pkg/logger"
	"go.uber.org/zap"
)

//	@title			Study Platform API
//	@version		1.0.0
//	@description	API for managing courses, students, and teachers.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Ivan Popov
//	@contact.email	vanap387@gmail.com

//	@host		localhost:3800
//	@BasePath	/

// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				Type "Bearer " followed by your JWT token.
func main() {
	ctx := context.Background()
	if err := config.Load("env.dev"); err != nil {
		logger.Fatal(ctx, "load config:", err)
	}

	if config.Config().App.IsProd() {
		logger.SetLevel(zap.InfoLevel)
	} else {
		logger.SetLevel(zap.DebugLevel)
	}

	closer.SetShutdownTimeout(config.Config().App.GracefulShutdownTimeout())

	a, err := app.New(ctx)
	if err != nil {
		logger.FatalKV(ctx, "new application", "error", err.Error())
	}

	if err := a.Run(ctx); err != nil {
		logger.FatalKV(ctx, "run application", "error", err.Error())
	}

	closer.Wait()
}
