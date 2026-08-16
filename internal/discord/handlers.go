package discord

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/markbeep/syncfo/internal/config"
	"go.uber.org/zap"
)

type Handlers struct {
	cmdIDs map[string]string
	appID  string
}

func (h *Handlers) SetupHandlers(s *discordgo.Session, logger *zap.Logger, sec config.Secret, db *sql.DB) {
	var once sync.Once
	appIDCh := make(chan string, 1)

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		logger.Info("Bot is ready", zap.String("user", r.User.Username), zap.String("appID", r.Application.ID))
		once.Do(func() {
			appIDCh <- r.Application.ID
			close(appIDCh)
		})
	})

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			switch i.ApplicationCommandData().Name {
			case "setup":
				if err := HandleTrackedCommand(ctx, s, i, db); err != nil {
					logger.Error("Failed to handle tracked command", zap.Error(err))
					logger.Error("Failed to handle message", zap.Error(err))
					s.InteractionRespond(i.Interaction, errResponse("Failed to build interactive setup message"))
				}
			}

		case discordgo.InteractionMessageComponent:
			parts := strings.Split(i.MessageComponentData().CustomID, ":")
			switch parts[0] {
			case PrevPagePrefix, NextPagePrefix, EnabledPrefix, DisabledPrefix:
				if err := HandleInteraction(ctx, s, i, db); err != nil {
					logger.Error("Failed to handle message", zap.Error(err))
					s.InteractionRespond(i.Interaction, errResponse("Failed to handle message"))
				}
			}
		}
	})

	go func() {
		h.appID = <-appIDCh
		rcmd, err := s.ApplicationCommandCreate(h.appID, sec.Discord.TestGuild, &discordgo.ApplicationCommand{
			Name:        "setup",
			Description: "Setup what channels and forums should be followed",
		})
		if err != nil {
			logger.Fatal("Failed to create application command", zap.Error(err))
		}
		h.cmdIDs = make(map[string]string)
		h.cmdIDs[rcmd.ID] = rcmd.Name
	}()
}

func (h *Handlers) Cleanup(s *discordgo.Session, testGuild string) error {
	if h.appID == "" {
		return nil
	}
	for cmdID := range h.cmdIDs {
		if err := s.ApplicationCommandDelete(h.appID, testGuild, cmdID); err != nil {
			return fmt.Errorf("failed to delete application command %s: %w", cmdID, err)
		}
	}
	return nil
}
