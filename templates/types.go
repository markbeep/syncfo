package page

import "time"

type Message struct {
	Author    string
	Markdown  string
	HTML      string
	CreatedAt time.Time
	EditedAt  *time.Time
}

type Thread struct {
	ID        int64
	Title     string
	Server    string
	Channel   string
	CreatedAt time.Time
	Messages  []Message
}

type ThreadSummary struct {
	ID            int64
	Title         string
	Server        string
	Channel       string
	Author        string
	CreatedAt     time.Time
	ReplyCount    int
	LastMessageAt time.Time
}

type ListPageData struct {
	Number      int
	StablePages int
	Threads     []ThreadSummary
	Preview     []ThreadSummary
	NextNumber  int
	PrevNumber  int
}

type ThreadPageData struct {
	ThreadID   int64
	Title      string
	Server     string
	Channel    string
	CreatedAt  time.Time
	Messages   []Message
	Page       int
	TotalPages int
}
