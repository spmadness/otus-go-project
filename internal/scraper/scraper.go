package scraper

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type MetricType int32

const (
	System MetricType = iota
	CPU
)

var (
	ErrEmptyParseCommand     = errors.New("no parse command given")
	ErrParseInitialData      = errors.New("wrong initial data")
	ErrEmptyScraperData      = errors.New("empty scraper data")
	ErrParseValues           = errors.New("wrong data values")
	ErrSecondsValue          = errors.New("seconds must be positive integer")
	ErrStringOutputEmptyData = errors.New("empty data for string output processing")
)

var scraperRegex = regexp.MustCompile(`^\.?(\d+\.\d+) \.?(\d+\.\d+) \.?(\d+\.\d+)$`)

type Scraper interface {
	GetData()
	ClearData()
	ParseData(data string, dataRowCnt int) error
	GetSnapshot(seconds int64) ([]string, error)
	GetSnapshotHeaders() string
	GetSnapshotFormat() string
	GetSnapshotDataRowElementsCnt() int
	HasConnections() bool
	AddConnection()
	RemoveConnection()
	command() *exec.Cmd
}

type DataParse interface {
	ParseData(data string) error
}

type BaseScraper struct {
	data []string

	connections uint
	mu          sync.Mutex
	re          *regexp.Regexp
}

func (bs *BaseScraper) ClearData() {
	bs.data = []string{}
}

func (bs *BaseScraper) HasConnections() bool {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	return bs.connections > 0
}

func (bs *BaseScraper) AddConnection() {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	bs.connections++
}

func (bs *BaseScraper) RemoveConnection() {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	bs.connections--
}

func (bs *BaseScraper) ParseData(data string, dataRowCnt int) error {
	var parsedData string
	matches := bs.re.FindStringSubmatch(data)

	if len(matches) != dataRowCnt+1 {
		return ErrParseInitialData
	}
	for k, v := range matches {
		if k == 0 {
			continue
		}
		parsedData += v
		if k < len(matches)-1 {
			parsedData += " "
		}
	}
	bs.data = append(bs.data, parsedData)

	return nil
}

func NewCollection(system, cpu bool) map[MetricType]Scraper {
	m := make(map[MetricType]Scraper)

	if system {
		m[System] = NewScraperSystem()
		log.Println("system load scraper enabled")
	}
	if cpu {
		m[CPU] = NewScraperCPU()
		log.Println("cpu scraper enabled")
	}
	return m
}

func getLastN(data []string, n int64) []string {
	length := int64(len(data))
	if n > length {
		n = length
	}
	return data[length-n:]
}

func calculateAverage(numbers []float64) float64 {
	sum := 0.0
	if len(numbers) == 0 {
		return sum
	}

	for _, num := range numbers {
		sum += num
	}
	return sum / float64(len(numbers))
}

func fetchData(cmd *exec.Cmd) (string, error) {
	if cmd == nil {
		return "", ErrEmptyParseCommand
	}
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error executing command %v: %s", cmd.Args, err.Error())
	}

	return strings.TrimSpace(string(output)), err
}

func prepareData(data []string, seconds int64, dataRowCnt int) ([][]float64, error) {
	result := make([][]float64, dataRowCnt)
	for i := range result {
		result[i] = make([]float64, 0)
	}
	if seconds < 1 {
		return result, ErrSecondsValue
	}

	items := getLastN(data, seconds)
	if len(items) == 0 {
		return result, ErrEmptyScraperData
	}

	for _, i := range items {
		fields := strings.Fields(i)
		if len(fields) != dataRowCnt {
			return result, ErrParseValues
		}
		for k, f := range fields {
			float, err := strconv.ParseFloat(f, 64)
			if err != nil {
				return result, ErrParseValues
			}
			result[k] = append(result[k], float)
		}
	}

	return result, nil
}

func resultString(headers string, format string, values ...[]float64) ([]string, error) {
	var result []string
	var sb strings.Builder

	if len(values) == 0 {
		return result, ErrStringOutputEmptyData
	}

	sb.WriteString(headers)

	for k, v := range values {
		sb.WriteString(fmt.Sprintf(format, calculateAverage(v)))
		if k < len(values)-1 {
			sb.WriteString(" ")
		}
	}

	result = append(result, sb.String())

	return result, nil
}

func snapshot(data []string, seconds int64, elemsCnt int, headers string, format string) ([]string, error) {
	var result []string
	preparedData, err := prepareData(data, seconds, elemsCnt)
	if err != nil {
		return result, err
	}
	result, err = resultString(headers, format, preparedData...)
	if err != nil {
		return result, err
	}

	return result, nil
}
