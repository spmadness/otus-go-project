package app

import (
	"context"
	"fmt"
	"github.com/spmadness/otus-go-project/internal/scraper"
	"sync"
	"time"
)

type App struct {
	scrapers map[scraper.MetricType]scraper.Scraper
}

func (a *App) GatherMetrics(ctx context.Context) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, s := range a.Scrapers() {
		wg.Add(1)
		s := s
		go func(scraper scraper.Scraper) {
			defer wg.Done()
			scraper.Init()

			ticker := time.NewTicker(time.Second)
			for {
				select {
				case <-ctx.Done():
					ticker.Stop()
					return
				case <-ticker.C:
					if !scraper.HasConnections() {
						scraper.ClearData()
						continue
					}
					mu.Lock()
					scraper.GetData()
					mu.Unlock()
				}
			}
		}(s)
	}
	wg.Wait()
}

func (a *App) Scrapers() map[scraper.MetricType]scraper.Scraper {
	return a.scrapers
}

func (a *App) Scraper(code scraper.MetricType) (scraper.Scraper, error) {
	s, ok := a.scrapers[code]
	if !ok {
		return nil, fmt.Errorf("no active scraper with code %d", code)
	}

	return s, nil
}

func New(scrapers map[scraper.MetricType]scraper.Scraper) *App {
	return &App{
		scrapers: scrapers,
	}
}
