package discord

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/markbeep/syncfo/internal/database"
)

const PAGE_SIZE = 20

const (
	PrevPagePrefix = "prev_page"
	NextPagePrefix = "next_page"
	EnabledPrefix  = "enable"
	DisabledPrefix = "disable"
)

var cachedChannels = make(map[string]struct {
	chs []*discordgo.Channel
	ts  time.Time
})

func HandleTrackedCommand(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, db *sql.DB) error {
	rows, err := buildTrackedRows(ctx, s, db, i.GuildID, 0)
	if err != nil {
		return fmt.Errorf("failed to build tracked page: %w", err)
	}
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Components: rows,
			Content:    "Tracked channels. Red=untracked, Green=tracked",
			Flags:      discordgo.MessageFlagsEphemeral,
			CustomID:   "tracked_channels",
		},
	}); err != nil {
		return fmt.Errorf("failed to respond to interaction: %w", err)
	}
	return nil
}

func HandleInteraction(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, db *sql.DB) error {
	customID := i.MessageComponentData().CustomID
	parts := strings.Split(customID, ":")
	action := parts[0]
	target := parts[1]
	page, err := strconv.Atoi(parts[2])
	if err != nil {
		return fmt.Errorf("failed to parse page number: %w", err)
	}

	switch action {
	case EnabledPrefix, DisabledPrefix:
		chID, err := strconv.ParseInt(target, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse channel ID: %w", err)
		}
		queries := database.New(db)
		if action == EnabledPrefix {
			if err := queries.SetChannelUntracked(ctx, chID); err != nil {
				return fmt.Errorf("failed to set channel untracked: %w", err)
			}
		} else {
			gID, err := strconv.ParseInt(i.GuildID, 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse guild ID: %w", err)
			}
			if err := queries.SetChannelTracked(ctx, database.SetChannelTrackedParams{ID: chID, GuildID: gID}); err != nil {
				return fmt.Errorf("failed to set channel tracked: %w", err)
			}
		}
	case PrevPagePrefix, NextPagePrefix:
	}

	rows, err := buildTrackedRows(ctx, s, db, i.GuildID, page)
	if err != nil {
		return fmt.Errorf("failed to build tracked page: %w", err)
	}
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Components: rows,
			Content:    "Tracked channels. Red=untracked, Green=tracked",
			CustomID:   "tracked_channels",
		},
	}); err != nil {
		return fmt.Errorf("failed to respond to interaction: %w", err)
	}

	return nil
}

func buildTrackedRows(ctx context.Context, s *discordgo.Session, db *sql.DB, guildID string, page int) ([]discordgo.MessageComponent, error) {
	channels, err := sortedChannels(ctx, s, guildID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch guild channels: %w", err)
	}
	if len(channels) == 0 {
		return nil, fmt.Errorf("no channels found in guild %s", guildID)
	}

	pages := (len(channels) + PAGE_SIZE - 1) / PAGE_SIZE
	page = min(pages, page)
	channels = channels[page*PAGE_SIZE : min(len(channels), (page+1)*PAGE_SIZE)]

	var ids []int64
	for _, ch := range channels {
		id, err := strconv.ParseInt(ch.ID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse channel ID %s: %w", ch.ID, err)
		}
		ids = append(ids, id)
	}

	queries := database.New(db)
	tracked := make(map[string]struct{})
	res, err := queries.ListTrackedChannels(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to list tracked DB channels: %w", err)
	}
	for _, ch := range res {
		tracked[strconv.FormatInt(ch, 10)] = struct{}{}
	}

	var rows []discordgo.MessageComponent
	for chunk := range slices.Chunk(channels, 5) {
		var rowComponents []discordgo.MessageComponent
		for _, ch := range chunk {
			style := discordgo.DangerButton
			prefix := DisabledPrefix
			if _, ok := tracked[ch.ID]; ok {
				style = discordgo.SuccessButton
				prefix = EnabledPrefix
			}
			rowComponents = append(rowComponents, discordgo.Button{
				Label:    "#" + ch.Name,
				Style:    style,
				CustomID: fmt.Sprintf("%s:%s:%d", prefix, ch.ID, page),
			})
		}
		rows = append(rows, &discordgo.ActionsRow{Components: rowComponents})
	}

	rows = append(rows, &discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{Label: "◀ Prev", CustomID: fmt.Sprintf("%s::%d", PrevPagePrefix, page-1), Disabled: page == 0},
			discordgo.Button{Label: "Next ▶", CustomID: fmt.Sprintf("%s::%d", NextPagePrefix, page+1), Disabled: page >= pages-1},
		},
	})

	return rows, nil
}

func sortedChannels(ctx context.Context, s *discordgo.Session, guildID string) ([]*discordgo.Channel, error) {
	if cached, ok := cachedChannels[guildID]; ok && time.Since(cached.ts) < 5*time.Minute {
		return cached.chs, nil
	}
	channels, err := s.GuildChannels(guildID, discordgo.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	invalid := map[discordgo.ChannelType]struct{}{
		discordgo.ChannelTypeGuildCategory:  {},
		discordgo.ChannelTypeGuildDirectory: {},
		discordgo.ChannelTypeGuildStore:     {},
	}
	channels = slices.DeleteFunc(channels, func(s *discordgo.Channel) bool {
		_, ok := invalid[s.Type]
		return ok
	})
	slices.SortFunc(channels, func(a, b *discordgo.Channel) int {
		al, bl := strings.ToLower(a.Name), strings.ToLower(b.Name)
		if al > bl {
			return 1
		}
		if al < bl {
			return -1
		}
		return 0
	})
	cachedChannels[guildID] = struct {
		chs []*discordgo.Channel
		ts  time.Time
	}{chs: channels, ts: time.Now()}
	return channels, nil
}

func errResponse(message string) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}
}
