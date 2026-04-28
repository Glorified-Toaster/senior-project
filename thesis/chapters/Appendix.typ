= Appendices

== Appendix A: Source Code

=== JWT Authentication Middleware

```go
// internal/adapters/inbound/http/middleware/auth_middleware.go

func (auth *AuthMiddleware) AuthenticationMiddleware() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        authHeader, err := ctx.Cookie("auth_token")
        if err != nil || authHeader == "" {
            handleUnauthorized(ctx)
            return
        }

        claims, err := auth.jwt.ValidateToken(authHeader)
        if err != nil {
            auth.logger.LogErrorWithLevel("error", "HTTP_SERVER_ERROR",
                "TOKEN_VALIDATION_ERROR", "Token validation failed", err)
            ctx.SetCookie("auth_token", "", -1, "/", "", false, true)
            handleUnauthorized(ctx)
            return
        }

        if claims != nil && claims.ExpiresAt < time.Now().Unix() {
            ctx.SetCookie("auth_token", "", -1, "/", "", false, true)
            handleUnauthorized(ctx)
            return
        }

        setClaimsInContext(ctx, claims)
        ctx.Next()
    }
}
```

=== Graceful Shutdown

```go
func (s *Server) startWithGracefulShutdown(certFile, keyFile string) {
	errChan := make(chan error, 1)
	go func() {
		var err error

		// check if certFile and keyFile are provided
		if certFile != "" && keyFile != "" {
			err = s.server.ListenAndServeTLS(certFile, keyFile)
		} else {
			s.logger.LogErrorWithLevel("warn", logger.FailedToStartWithTLS.Type, logger.FailedToStartWithTLS.Code, logger.FailedToStartWithTLS.Msg, errors.New("unable to get TLS cert"))
			err = s.server.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("failed to start HTTP server: %v", err)
		}
	}()
	s.waitForShutdownSignal(errChan)
}
```

=== Application Transaction Pattern

```go
// Example: CreateExam in internal/application/exam_application.go

func (app *Application) CreateExam(ctx context.Context, 
    arg ports.CreateExamParams) (domain.Exam, error) {
    var exam domain.Exam
    err := app.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
        var err error
        exam, err = app.examRepo.Create(txCtx, arg)
        return err
    })
    return exam, err
}

```
=== Dependency Injection Pattern (func main)

```go
// cmd/api/main.go
func main() {
	cfg, err := configInstance.GetConfig()
	if err != nil {
		log.Fatalf("Failed to get configuration: %v", err)
	}

	// Initialize Zap
	zapLogger, err := logger.InitZapLogger(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer func() {
		err := zapLogger.Sync()
		if err != nil && !errors.Is(err, syscall.EINVAL) {
			fmt.Fprintf(os.Stderr, "Failed to sync logger: %v\n", err)
		}
	}()

	// Initialize Logger
	zlog := logger.New(zapLogger)

	DBConfig := &database.DBConfig{
		// mapping from YAML config to DBConfig struct
	}

	pool, err := database.NowConnection(*DBConfig)
	if err != nil {
		zlog.LogErrorWithLevel("fatal", logger.DatabaseError, logger.MongoFailedToConnect.Code, "Database failed to connect", err)
		return
	}
	query := sqlc.New(pool)
	userRepo := repository.NewUserRepository(query)
	subjectRepo := repository.NewSubjectRepository(query)
	examRepo := repository.NewExamRepository(query)
	questionRepo := repository.NewQuestionRepository(query)
	txManager := database.NewPostgresTxManager(pool.Pool)
	localDisk := storage.NewLocalDiskAdapter()
	app := application.NewApplication(userRepo, subjectRepo, examRepo, questionRepo, localDisk, txManager, pool, zlog)

	// init validator
	validate := validator.New()
	// init jwt
	jwt := helpers.NewJWT(cfg)
	// init auth middleware
	authMiddleware := middleware.NewAuthMiddleware(jwt, zlog)
	// pass cache, repo, validator, jwt to controllers
	userCtrl := handler.NewUserHandler(app, validate, jwt, cfg, zlog)

	// initialize the server
	srv := server.NewServer(userCtrl, authMiddleware, zlog, cfg)

	zlog.LogInfo(logger.ServerStartOK.Type, logger.ServerStartOK.Msg, zap.String("server_address", net.JoinHostPort(cfg.HTTPServer.Addr, cfg.HTTPServer.Port)))

	// start the server over TLS
	srv.StartOverTLS(cfg)
}
```

=== Database functions

```sql
-- Function to check if the number of choices is less than 4
CREATE OR REPLACE FUNCTION check_choices_count()
RETURNS TRIGGER AS $$
DECLARE
    choice_count INT;
BEGIN
    SELECT COUNT(*) INTO choice_count 
    FROM choices 
    WHERE question_id = NEW.question_id;
    
    IF choice_count >= 4 THEN
        RAISE EXCEPTION 'Maximum 4 choices allowed per question';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
```

== Appendix B: User Interface Screens

#figure(
  rect(image("../assets/Admin-Login.png", width: 100%)),
  caption: "Login Page"
)

#figure(
  image("../assets/Admin-Dashboard.png", width: 100%),
  caption: "Admin Dashboard Page"
)

#figure(
  image("../assets/Edit-Exam.png", width: 100%),
  caption: "Edit Exam Page"
)

#figure(
  image("../assets/Edit-Subject-dashboard.png", width: 100%),
  caption: "Edit Subject Dashboard Page"
)

#figure(
  image("../assets/All-Users.png", width: 100%),
  caption: "All Users Page"
)

#figure(
  image("../assets/Student-Dashboard.png", width: 100%),
  caption: "Student Dashboard Page"
)