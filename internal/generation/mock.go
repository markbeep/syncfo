package generation

import (
	"fmt"
	"time"

	page "github.com/markbeep/syncfo/templates"
)

var mockServers = []string{"syncfo", "golang", "rust", "kubernetes"}
var mockChannels = []string{"support", "general", "help", "showcase"}
var mockBase = time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)

func mockMessageCount(i int) int {
	return 1 + (i % 25)
}

// MockThreads returns 45 thread summaries ordered oldest -> newest, exercising
// stable pages plus /list/0 blending (45 = 2 stable pages + 5 newest).
func MockThreads() []page.ThreadSummary {
	const n = 45
	out := make([]page.ThreadSummary, 0, n)
	for i := 0; i < n; i++ {
		replies := mockMessageCount(i) - 1
		created := mockBase.Add(time.Duration(i) * 24 * time.Hour)
		out = append(out, page.ThreadSummary{
			ID:            int64(1001 + i),
			Title:         fmt.Sprintf("Thread %d: question about %s", i+1, mockChannels[i%len(mockChannels)]),
			Server:        mockServers[i%len(mockServers)],
			Channel:       mockChannels[i%len(mockChannels)],
			Author:        fmt.Sprintf("user%d", i),
			CreatedAt:     created,
			ReplyCount:    replies,
			LastMessageAt: created.Add(time.Duration(replies) * 30 * time.Minute),
		})
	}
	return out
}

// MockThreadsDetail returns the full thread structs (with messages) for the
// summaries produced by MockThreads. IDs line up so list links resolve.
func MockThreadsDetail() []page.Thread {
	summaries := MockThreads()
	out := make([]page.Thread, 0, len(summaries))
	for i, s := range summaries {
		count := mockMessageCount(i)
		msgs := make([]page.Message, 0, count)
		for j := range count {
			author := s.Author
			if j > 0 {
				author = fmt.Sprintf("user%d", i+j)
			}
			msgs = append(msgs, page.Message{
				Author:    author,
				Markdown:  mockMarkdown(i, j),
				CreatedAt: s.CreatedAt.Add(time.Duration(j) * 30 * time.Minute),
			})
		}
		out = append(out, page.Thread{
			ID:        s.ID,
			Title:     s.Title,
			Server:    s.Server,
			Channel:   s.Channel,
			CreatedAt: s.CreatedAt,
			Messages:  msgs,
		})
	}
	return out
}

func mockMarkdown(i, j int) string {
	if j == 0 {
		return fmt.Sprintf(`## Question %d

Here is the **original post** for thread %d.

Some *italic* text and ~~strikethrough~~ text.

`+"```go"+`
func main() {
	fmt.Println("hello")
}
`+"```"+`

- item one
- item two
- item three

> a blockquote with **bold**

See [the docs](https://example.com) for more.
`, i+1, i+1)
	}
	return fmt.Sprintf(`Reply %d to thread %d with **bold**, *italic* and ~~gone~~ text.

`+"```"+`
code block in reply
`+"```"+`
`, j, i+1)
}
