package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

// config 타입은 애플리케이션이 사용할 포트 및
// 운영 환경(development, staging, production)을 정의합니다.
type config struct {
	port int
	env  string
}

// application 타입은 핸들러, 헬퍼, 미들웨어 등에
// 대한 의존성을 정의합니다.
type application struct {
	config config
	logger *slog.Logger
}

func main() {

	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "api server port")
	flag.StringVar(&cfg.env, "env", "development", "operation environment (development|staging|production)")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{
		config: cfg,
		logger: logger,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  1 * time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)

	err := srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)

	//fmt.Printf("%+v", app)ﬁ
}
