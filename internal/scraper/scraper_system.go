package scraper

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type SystemScraper struct {
	BaseScraper
}

func (s *SystemScraper) Init() {
	s.re = scraperRegex
}

func (s *SystemScraper) FetchData() (string, error) {
	c := "top -bn 1 | sed -n 's/^.*load average: //p' | sed 's/, / /g' | sed 's/,/./g'"
	cmd := exec.Command("bash", "-c", c)

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error executing command %s: %s", c, err.Error())
	}

	return string(output), err
}

func (s *SystemScraper) ParseData(data string) error {
	matches := s.re.FindStringSubmatch(data)

	if len(matches) != 4 {
		return fmt.Errorf("error parsing data")
	}
	parsedData := fmt.Sprintf("%s %s %s", matches[1], matches[2], matches[3])
	s.data = append(s.data, parsedData)

	return nil
}

func (s *SystemScraper) GetData() {
	data, err := s.FetchData()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = s.ParseData(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("scraper system data len: %d\n", len(s.data))
}

func (s *SystemScraper) GetSnapshot(seconds int64) ([]string, error) {
	var result []string
	var sb strings.Builder

	items := getLastN(s.data, seconds)
	if len(items) == 0 {
		return result, errors.New("empty scraper data")
	}

	avg1 := make([]float64, 0)
	avg5 := make([]float64, 0)
	avg15 := make([]float64, 0)

	for _, i := range items {
		f := strings.Fields(i)
		if len(f) != 3 {
			return result, errors.New("data parse error")
		}
		float, err := strconv.ParseFloat(strings.TrimRight(f[0], "."), 64)
		if err != nil {
			return result, errors.New("data parse error")
		}
		avg1 = append(avg1, float)

		float, err = strconv.ParseFloat(strings.TrimRight(f[1], "."), 64)
		if err != nil {
			return result, errors.New("data parse error")
		}
		avg5 = append(avg5, float)

		float, err = strconv.ParseFloat(f[2], 64)
		if err != nil {
			return result, errors.New("data parse error")
		}
		avg15 = append(avg15, float)
	}

	sb.WriteString(s.GetSnapshotHeaders())
	sb.WriteString(fmt.Sprintf("%.2f %.2f %.2f",
		calculateAverage(avg1),
		calculateAverage(avg5),
		calculateAverage(avg15),
	))

	result = append(result, sb.String())

	return result, nil
}

func (s *SystemScraper) GetSnapshotHeaders() string {
	return "Load average\n1m 5m 15m\n"
}
