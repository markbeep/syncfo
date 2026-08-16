package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/markbeep/syncfo/internal/config"
	"github.com/markbeep/syncfo/internal/discord"
	"go.uber.org/zap"
	_ "modernc.org/sqlite"
)

func main() {
	pub, sec, err := config.Parse()
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}
	_ = sec

	var logger *zap.Logger
	if pub.Logging.Debug {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	dbPragmas := "?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=synchronous(NORMAL)&_pragma=foreign_keys(1)&_pragma=journal_size_limit(67108864)&_pragma=mmap_size(268435456)"
	db, err := sql.Open("sqlite", fmt.Sprintf("file:%s%s", pub.Database.Path, dbPragmas))
	if err != nil {
		logger.Fatal("Failed to open database", zap.Error(err))
	}
	defer db.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	s, err := discordgo.New("Bot " + sec.Discord.Token)
	if err != nil {
		logger.Fatal("Failed to create Discord session", zap.Error(err))
	}

	h := discord.Handlers{}
	h.SetupHandlers(s, logger, sec, db)

	if err := s.Open(); err != nil {
		logger.Fatal("Failed to open Discord session", zap.Error(err))
	}
	defer s.Close()

	<-ctx.Done()
	logger.Info("Shutting down gracefully...")
}
