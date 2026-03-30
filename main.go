package main

import (
	"log"
	"net/http"
	"os"
	"syscall"
	"time"
	"context"
	"os/signal"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/handler/router"
)

func main() {
	err := realMain()
	if err != nil {
		log.Println("ERROR:", err)
	}
	log.Println("program exited")

}

func realMain() error {
	// config values
	const (
		defaultPort   = ":8010"
		defaultDBPath = ".sqlite3/todo.db"
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = defaultDBPath
	}

	// set time zone
	var err error
	time.Local, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}

	// set up sqlite3
	todoDB, err := db.NewDB(dbPath)
	if err != nil {
		return err
	}
	defer todoDB.Close()

	// NOTE: 新しいエンドポイントの登録はrouter.NewRouterの内部で行うようにする
	log.Println("starting server...")

	mux := router.NewRouter(todoDB)

	log.Println("calling ListenAndServe...")

	// TODO: サーバーをlistenする

	server := &http.Server{
		Addr:	port,
		Handler: mux,
	}

	ctx,stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Println("server started on", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("listen error:", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Println("shutdown error:", err)
	}

	log.Println("server gracefully stopped")
	return err
}
