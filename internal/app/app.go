package app

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/escoutdoor/study-platform/internal/config"
	"github.com/escoutdoor/study-platform/internal/handler/auth"
	"github.com/escoutdoor/study-platform/internal/handler/course"
	"github.com/escoutdoor/study-platform/internal/handler/student"
	"github.com/escoutdoor/study-platform/internal/handler/teacher"
	"github.com/escoutdoor/study-platform/internal/handler/user"
	"github.com/escoutdoor/study-platform/internal/middleware"
	"github.com/escoutdoor/study-platform/pkg/closer"
	"github.com/escoutdoor/study-platform/pkg/errwrap"
	"github.com/escoutdoor/study-platform/pkg/logger"
	"github.com/escoutdoor/study-platform/pkg/validator"
	"github.com/pressly/goose/v3"
	httpswagger "github.com/swaggo/http-swagger"
)

type App struct {
	di *di

	httpServer *http.Server
}

func New(ctx context.Context) (*App, error) {
	app := &App{di: newDi()}
	if err := app.initDeps(ctx); err != nil {
		return nil, err
	}

	if err := goose.SetDialect(string(goose.DialectPostgres)); err != nil {
		return nil, errwrap.Wrap("set migrations dialect", err)
	}

	if err := goose.UpContext(ctx, app.di.DB(ctx).Conn(), config.Config().Postgres.MigrationsDir()); err != nil {
		return nil, errwrap.Wrap("migrate up", err)
	}

	return app, nil
}

func (a *App) Run(ctx context.Context) error {
	go func() {
		logger.InfoKV(ctx, "http server is running", "address", config.Config().HttpServer.Address())
		if err := a.runHttpServer(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(ctx, "run http server: ", err)
		}
	}()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	deps := []func(ctx context.Context) error{
		a.initHttpServer,
	}

	for _, d := range deps {
		if err := d(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initHttpServer(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.Handle("/swagger/", httpswagger.WrapHandler)

	cv := validator.New()

	authMiddleware := middleware.Auth(a.di.TokenProvider())

	auth.RegisterHandlers(mux, a.di.AuthService(ctx), cv)
	user.RegisterHandlers(mux, a.di.UserService(ctx), cv, authMiddleware)
	student.RegisterHandlers(
		mux,
		a.di.StudentService(ctx),
		a.di.EnrollmentService(ctx),
		cv,
		authMiddleware,
	)
	teacher.RegisterHandlers(mux, a.di.TeacherService(ctx), cv, authMiddleware)
	course.RegisterHandlers(mux, a.di.CourseService(ctx), cv, authMiddleware)

	s := &http.Server{
		Addr:              config.Config().HttpServer.Address(),
		Handler:           middleware.Logging(mux),
		ReadTimeout:       time.Second * 5,
		ReadHeaderTimeout: time.Second * 5,
	}

	a.httpServer = s
	closer.Add(func(ctx context.Context) error {
		return a.httpServer.Shutdown(ctx)
	})

	return nil
}

func (a *App) runHttpServer() error {
	if err := a.httpServer.ListenAndServe(); err != nil {
		return errwrap.Wrap("http server listen and serve", err)
	}

	return nil
}
