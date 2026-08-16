package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/markbeep/syncfo/internal/config"
	"github.com/markbeep/syncfo/internal/database"
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

	totalTime := time.Duration(0)
	totalCount := 0

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		active := sync.Map{}

		for {
			select {
			case <-ctx.Done():
				logger.Info("Received interrupt signal, stopping separate...")
				return
			case <-time.Tick(2 * time.Second):
			}
			queries := database.New(db)
			ts, err := queries.ListUnbuiltThreads(ctx, database.ListUnbuiltThreadsParams{Limit: 100})
			if err != nil {
				logger.Fatal("Failed to list unbuilt threads", zap.Error(err))
			}
			for _, t := range ts {
				if _, loaded := active.LoadOrStore(t, struct{}{}); loaded {
					continue
				}

				go func(threadID int64) {
					for {
						msgs, err := queries.ListUnbuiltMessages(ctx, database.ListUnbuiltMessagesParams{ThreadID: threadID, Limit: 100})
						if err != nil {
							logger.Fatal("Failed to list unbuilt messages", zap.Error(err))
						}

						// add to whatever datastructure
						time.Sleep(20 * time.Millisecond)

						for _, m := range msgs {
							queries.SetMessageBuilt(ctx, m.ID)
						}

						logger.Debug("Built messages", zap.Int64("threadId", threadID), zap.Int("messageCount", len(msgs)))

						if len(msgs) == 100 {
							break
						}
					}
					active.Delete(threadID)
				}(t)

			}
		}
	}()

out:
	for {
		wg := sync.WaitGroup{}
		select {
		case <-ctx.Done():
			break out
		default:
		}
		for range 100 {
			wg.Go(func() {
				start := time.Now()
				if err := bench(ctx, logger, db); err != nil {
					if ctx.Err() != nil {
						return
					}
					logger.Fatal("Benchmark failed", zap.Error(err))
				}
				totalTime += time.Since(start)
				totalCount++
			})

		}
		wg.Wait()
	}

	if totalCount == 0 {
		logger.Info("No benchmarks finished yet")
	} else {
		logger.Info("Total time", zap.Duration("totalTime", totalTime), zap.Int("totalCount", totalCount), zap.Duration("avgTime", totalTime/time.Duration(totalCount)))
	}

	// discord, err := discordgo.New("Bot " + sec.Discord.Token)
	// if err != nil {
	// 	logger.Fatal("Failed to create Discord session", zap.Error(err))
	// }

}

func bench(ctx context.Context, logger *zap.Logger, db *sql.DB) error {
	start := time.Now()

	// tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	// if err != nil {
	// 	return fmt.Errorf("failed to begin transaction: %w", err)
	// }
	// defer tx.Rollback()
	queries := database.New(db) //.WithTx(tx)

	threadId := rand.Int64()
	guildId := rand.Int64()

	t, err := queries.InsertThread(ctx, database.InsertThreadParams{ID: threadId, GuildID: guildId})
	if err != nil {
		return fmt.Errorf("failed to insert thread (threadId=%d, guildId=%d, elapsed=%s): %w", threadId, guildId, time.Since(start), err)
	}
	logger.Debug("Inserted thread", zap.Any("thread", t))

	innerWg := sync.WaitGroup{}
	for range rand.IntN(100) + 1 {
		innerWg.Go(func() {
			messageId := rand.Int64()
			_, err := queries.InsertMessage(ctx, database.InsertMessageParams{ID: messageId, ThreadID: t.ID})
			if err != nil {
				logger.Fatal("Failed to insert message", zap.Error(err), zap.Int64("messageId", messageId), zap.Int64("threadId", t.ID), zap.Duration("elapsed", time.Since(start)))
			}
		})
	}
	innerWg.Wait()
	logger.Debug("Inserted message", zap.Int64("threadId", t.ID), zap.Duration("elapsed", time.Since(start)))

	// if err := tx.Commit(); err != nil {
	// 	return fmt.Errorf("failed to commit transaction: %w", err)
	// }

	q := database.New(db)
	msgs, err := q.ListMessagesByThread(ctx, database.ListMessagesByThreadParams{ThreadID: t.ID, Limit: 100})
	if err != nil {
		return fmt.Errorf("failed to list messages by thread (threadId=%d, elapsed=%s): %w", t.ID, time.Since(start), err)
	}
	logger.Debug("Listed messages by thread", zap.Int64("threadId", t.ID), zap.Int("messageCount", len(msgs)))

	return nil
}
