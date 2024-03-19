package scraper

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type CPUScraper struct {
	BaseScraper
}

func (s *CPUScraper) Init() {
	s.re = scraperRegex
}

func (s *CPUScraper) FetchData() (string, error) {
	c := "top -bn 1 | awk '/Cpu/' | grep -oE \"[0-9,\\.]+ (us|sy|id)\" | awk '{printf $1\" \"}' | sed 's/,/./g'"
	cmd := exec.Command("bash", "-c", c)

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error executing command %s: %s", c, err.Error())
	}

	return string(output), err
}

func (s *CPUScraper) ParseData(data string) error {
	matches := s.re.FindStringSubmatch(data)

	if len(matches) != 4 {
		return fmt.Errorf("error parsing data")
	}
	parsedData := fmt.Sprintf("%s %s %s", matches[1], matches[2], matches[3])
	s.data = append(s.data, parsedData)

	return nil
}

func (s *CPUScraper) GetData() {
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

	fmt.Printf("scraper cpu data len: %d\n", len(s.data))
}

func (s *CPUScraper) GetSnapshot(seconds int64) ([]string, error) {
	var result []string
	var sb strings.Builder

	items := getLastN(s.data, seconds)
	if len(items) == 0 {
		return result, errors.New("empty scraper data")
	}

	us := make([]float64, 0)
	sy := make([]float64, 0)
	id := make([]float64, 0)

	for _, i := range items {
		f := strings.Fields(i)
		if len(f) != 3 {
			return result, errors.New("data parse error")
		}
		float, err := strconv.ParseFloat(f[0], 64)
		if err != nil {
			return result, errors.New("data parse error")
		}
		us = append(us, float)

		float, err = strconv.ParseFloat(f[1], 64)
		if err != nil {
			return result, errors.New("data parse error")
		}
		sy = append(sy, float)

		float, err = strconv.ParseFloat(f[2], 64)
		if err != nil {
			return result, errors.New("data parse error")
		}
		id = append(id, float)
	}

	sb.WriteString(s.GetSnapshotHeaders())
	sb.WriteString(fmt.Sprintf("%.1f %.1f %.1f",
		calculateAverage(us),
		calculateAverage(sy),
		calculateAverage(id),
	))

	result = append(result, sb.String())

	return result, nil
}

func (s *CPUScraper) GetSnapshotHeaders() string {
	return "Cpu Average\nuser_mode system_mode idle\n"
}
