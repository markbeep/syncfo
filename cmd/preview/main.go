package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/a-h/templ"
	"github.com/markbeep/syncfo/internal/config"
	"github.com/markbeep/syncfo/internal/generation"
)

func main() {
	pub, err := config.LoadPublic()
	if err != nil {
		log.Fatalf("Failed to load public config: %v", err)
	}
	basePath := pub.Static.BasePath
	if basePath == "" {
		basePath = "/"
	}

	ctx := context.Background()

	// Clear stale generated output while keeping the committed public/style.css.
	for _, dir := range []string{"public/list", "public/thread"} {
		if err := os.RemoveAll(dir); err != nil {
			log.Fatalf("Failed to remove %s: %v", dir, err)
		}
	}

	render := func(rel string, c templ.Component) {
		path := filepath.Join("public", rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			log.Fatalf("Failed to create directory for %s: %v", path, err)
		}
		f, err := os.Create(path)
		if err != nil {
			log.Fatalf("Failed to create %s: %v", path, err)
		}
		if err := c.Render(ctx, f); err != nil {
			f.Close()
			log.Fatalf("Failed to render %s: %v", path, err)
		}
		if err := f.Close(); err != nil {
			log.Fatalf("Failed to close %s: %v", path, err)
		}
	}

	// List pages: page 0 plus stable pages 1..StablePages.
	summaries := generation.MockThreads()
	stable := len(summaries) / generation.PageSize
	for n := 0; n <= stable; n++ {
		data := generation.ResolvePage(summaries, n)
		render(fmt.Sprintf("list/%d/index.html", n), generation.ListPage(data, basePath))
	}

	// Thread pages: page 1 at /thread/<id>/index.html, later pages nested.
	for _, t := range generation.MockThreadsDetail() {
		pages, err := generation.ThreadPages(t, basePath)
		if err != nil {
			log.Fatalf("Failed to build thread %d: %v", t.ID, err)
		}
		for i, c := range pages {
			rel := fmt.Sprintf("thread/%d/index.html", t.ID)
			if i > 0 {
				rel = fmt.Sprintf("thread/%d/%d/index.html", t.ID, i+1)
			}
			render(rel, c)
		}
	}

	log.Println("Preview written to public/")
}
