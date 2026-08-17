package generation

import (
	"github.com/a-h/templ"
	"github.com/markbeep/syncfo/templates"
)

// ThreadPageCount returns the number of message pages for a thread.
func ThreadPageCount(thread page.Thread) int {
	n := len(thread.Messages)
	if n == 0 {
		return 1
	}
	return (n + PageSize - 1) / PageSize
}

// ThreadPage renders a single page of a thread's messages.
func ThreadPage(thread page.Thread, number int, basePath string) (templ.Component, error) {
	total := ThreadPageCount(thread)
	if number < 1 {
		number = 1
	}
	if number > total {
		number = total
	}

	start := (number - 1) * PageSize
	end := start + PageSize
	if end > len(thread.Messages) {
		end = len(thread.Messages)
	}

	msgs := make([]page.Message, 0, end-start)
	for _, m := range thread.Messages[start:end] {
		html, err := Markdown(m.Markdown)
		if err != nil {
			return nil, err
		}
		m.HTML = html
		msgs = append(msgs, m)
	}

	data := page.ThreadPageData{
		ThreadID:   thread.ID,
		Title:      thread.Title,
		Server:     thread.Server,
		Channel:    thread.Channel,
		CreatedAt:  thread.CreatedAt,
		Messages:   msgs,
		Page:       number,
		TotalPages: total,
	}
	return page.ThreadView(data, basePath), nil
}

// ThreadPages renders every page of a thread's messages in order.
func ThreadPages(thread page.Thread, basePath string) ([]templ.Component, error) {
	total := ThreadPageCount(thread)
	pages := make([]templ.Component, 0, total)
	for p := 1; p <= total; p++ {
		c, err := ThreadPage(thread, p, basePath)
		if err != nil {
			return nil, err
		}
		pages = append(pages, c)
	}
	return pages, nil
}
