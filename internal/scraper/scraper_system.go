package scraper

import (
	"log"
)

type SystemScraper struct {
	BaseScraper
}

func NewScraperSystem() *SystemScraper {
	return &SystemScraper{
		BaseScraper: BaseScraper{
			re: scraperRegex,
		},
	}
}

func (s *SystemScraper) GetData() {
	data, err := fetchData(s.command())
	if err != nil {
		log.Println(err)
		return
	}

	err = s.ParseData(data, s.GetSnapshotDataRowElementsCnt())
	if err != nil {
		log.Println(err)
		return
	}
}

func (s *SystemScraper) GetSnapshot(seconds int64) ([]string, error) {
	result, err := snapshot(s.data, seconds,
		s.GetSnapshotDataRowElementsCnt(), s.GetSnapshotHeaders(), s.GetSnapshotFormat())
	if err != nil {
		return result, err
	}

	return result, nil
}

func (s *SystemScraper) GetSnapshotFormat() string {
	return "%.2f"
}

func (s *SystemScraper) GetSnapshotHeaders() string {
	return "Load average\n1m 5m 15m\n"
}

func (s *SystemScraper) GetSnapshotDataRowElementsCnt() int {
	return 3
}
