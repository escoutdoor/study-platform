package app

import (
	"context"

	"github.com/escoutdoor/study-platform/internal/config"
	course_repo "github.com/escoutdoor/study-platform/internal/repository/course"
	enrollment_repo "github.com/escoutdoor/study-platform/internal/repository/enrollment"
	student_repo "github.com/escoutdoor/study-platform/internal/repository/student"
	teacher_repo "github.com/escoutdoor/study-platform/internal/repository/teacher"
	user_repo "github.com/escoutdoor/study-platform/internal/repository/user"
	auth_svc "github.com/escoutdoor/study-platform/internal/service/auth"
	course_svc "github.com/escoutdoor/study-platform/internal/service/course"
	enrollment_svc "github.com/escoutdoor/study-platform/internal/service/enrollment"
	student_svc "github.com/escoutdoor/study-platform/internal/service/student"
	teacher_svc "github.com/escoutdoor/study-platform/internal/service/teacher"
	user_svc "github.com/escoutdoor/study-platform/internal/service/user"
	"github.com/escoutdoor/study-platform/internal/util/token"
	"github.com/escoutdoor/study-platform/pkg/closer"
	"github.com/escoutdoor/study-platform/pkg/database"
	"github.com/escoutdoor/study-platform/pkg/database/pq"
	"github.com/escoutdoor/study-platform/pkg/database/txmanager"
	"github.com/escoutdoor/study-platform/pkg/logger"
)

type di struct {
	db        database.DB
	txManager database.TxManager

	tokenProvider *token.TokenProvider

	userRepository       *user_repo.Repository
	studentRepository    *student_repo.Repository
	teacherRepository    *teacher_repo.Repository
	courseRepository     *course_repo.Repository
	enrollmentRepository *enrollment_repo.Repository

	authService       *auth_svc.Service
	userService       *user_svc.Service
	studentService    *student_svc.Service
	teacherService    *teacher_svc.Service
	courseService     *course_svc.Service
	enrollmentService *enrollment_svc.Service
}

func newDi() *di {
	return &di{}
}

func (d *di) DB(ctx context.Context) database.DB {
	if d.db == nil {
		db, err := pq.New(ctx, config.Config().Postgres.Dsn())
		if err != nil {
			logger.Fatal(ctx, "new database connection: ", err)
		}

		if err := db.Ping(ctx); err != nil {
			logger.Fatal(ctx, "ping new database connection: ", err)
		}

		d.db = db
		closer.Add(func(ctx context.Context) error {
			db.Close()
			return nil
		})
	}

	return d.db
}

func (d *di) TxManager(ctx context.Context) database.TxManager {
	if d.txManager == nil {
		d.txManager = txmanager.NewTransactionManager(d.DB(ctx))
	}

	return d.txManager
}

func (d *di) TokenProvider() *token.TokenProvider {
	if d.tokenProvider == nil {
		d.tokenProvider = token.NewTokenProvider(
			config.Config().JwtToken.AccessTokenSecretKey(),
			config.Config().JwtToken.RefreshTokenSecretKey(),
			config.Config().JwtToken.AccessTokenTTL(),
			config.Config().JwtToken.RefreshTokenTTL(),
		)
	}

	return d.tokenProvider
}

func (d *di) UserRepository(ctx context.Context) *user_repo.Repository {
	if d.userRepository == nil {
		d.userRepository = user_repo.New(d.DB(ctx))
	}

	return d.userRepository
}

func (d *di) StudentRepository(ctx context.Context) *student_repo.Repository {
	if d.studentRepository == nil {
		d.studentRepository = student_repo.New(d.DB(ctx))
	}

	return d.studentRepository
}

func (d *di) TeacherRepository(ctx context.Context) *teacher_repo.Repository {
	if d.teacherRepository == nil {
		d.teacherRepository = teacher_repo.New(d.DB(ctx))
	}

	return d.teacherRepository
}

func (d *di) CourseRepository(ctx context.Context) *course_repo.Repository {
	if d.courseRepository == nil {
		d.courseRepository = course_repo.New(d.DB(ctx))
	}

	return d.courseRepository
}

func (d *di) EnrollmentRepository(ctx context.Context) *enrollment_repo.Repository {
	if d.enrollmentRepository == nil {
		d.enrollmentRepository = enrollment_repo.New(d.DB(ctx))
	}

	return d.enrollmentRepository
}

func (d *di) AuthService(ctx context.Context) *auth_svc.Service {
	if d.authService == nil {
		d.authService = auth_svc.New(
			d.UserRepository(ctx),
			d.StudentRepository(ctx),
			d.TxManager(ctx),
			d.TokenProvider(),
		)
	}

	return d.authService
}

func (d *di) UserService(ctx context.Context) *user_svc.Service {
	if d.userService == nil {
		d.userService = user_svc.New(d.UserRepository(ctx))
	}

	return d.userService
}

func (d *di) StudentService(ctx context.Context) *student_svc.Service {
	if d.studentService == nil {
		d.studentService = student_svc.New(d.StudentRepository(ctx))
	}

	return d.studentService
}

func (d *di) TeacherService(ctx context.Context) *teacher_svc.Service {
	if d.teacherService == nil {
		d.teacherService = teacher_svc.New(d.TeacherRepository(ctx))
	}

	return d.teacherService
}

func (d *di) CourseService(ctx context.Context) *course_svc.Service {
	if d.courseService == nil {
		d.courseService = course_svc.New(d.CourseRepository(ctx))
	}

	return d.courseService
}

func (d *di) EnrollmentService(ctx context.Context) *enrollment_svc.Service {
	if d.enrollmentService == nil {
		d.enrollmentService = enrollment_svc.New(d.EnrollmentRepository(ctx), d.CourseRepository(ctx))
	}

	return d.enrollmentService
}
