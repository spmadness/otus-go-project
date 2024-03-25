package scraper

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type App struct {
	scrapers map[MetricType]Scraper
}

func (a *App) GatherMetrics(ctx context.Context) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, s := range a.Scrapers() {
		wg.Add(1)
		s := s
		go func(scraper Scraper) {
			defer wg.Done()

			ticker := time.NewTicker(time.Second)
			for {
				select {
				case <-ctx.Done():
					ticker.Stop()
					return
				case <-ticker.C:
					if !scraper.HasConnections() {
						mu.Lock()
						scraper.ClearData()
						mu.Unlock()
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

func (a *App) Scrapers() map[MetricType]Scraper {
	return a.scrapers
}

func (a *App) Scraper(code MetricType) (Scraper, error) {
	s, ok := a.scrapers[code]
	if !ok {
		return nil, fmt.Errorf("no active scraper with code %d", code)
	}

	return s, nil
}

func New(scrapers map[MetricType]Scraper) *App {
	return &App{
		scrapers: scrapers,
	}
}
