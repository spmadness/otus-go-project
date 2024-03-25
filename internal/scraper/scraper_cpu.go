package scraper

import (
	"log"
)

type CPUScraper struct {
	BaseScraper
}

func NewScraperCPU() *CPUScraper {
	return &CPUScraper{
		BaseScraper: BaseScraper{
			re: scraperRegex,
		},
	}
}

func (s *CPUScraper) GetSnapshotHeaders() string {
	return "Cpu Average\nuser_mode system_mode idle\n"
}

func (s *CPUScraper) GetSnapshotFormat() string {
	return "%.1f"
}

func (s *CPUScraper) GetSnapshotDataRowElementsCnt() int {
	return 3
}

func (s *CPUScraper) GetData() {
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

func (s *CPUScraper) GetSnapshot(seconds int64) ([]string, error) {
	result, err := snapshot(s.data, seconds,
		s.GetSnapshotDataRowElementsCnt(), s.GetSnapshotHeaders(), s.GetSnapshotFormat())
	if err != nil {
		return result, err
	}

	return result, nil
}
