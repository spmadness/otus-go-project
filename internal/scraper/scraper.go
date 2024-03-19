package scraper

import (
	"regexp"
	"strings"
	"sync"
)

type MetricType int32

const (
	System MetricType = iota
	CPU
	Disk
)

var scraperRegex = regexp.MustCompile(`(\d+\.\d+) (\d+\.\d+) (\d+\.\d+)`)

type Scraper interface {
	Init()
	GetData()
	ClearData()
	ParseData(data string) error
	GetSnapshot(seconds int64) ([]string, error)
	GetSnapshotHeaders() string
	HasConnections() bool
	AddConnection()
	RemoveConnection()
}

type DataParse interface {
	ParseData(data string) error
}

type BaseScraper struct {
	data []string

	re          *regexp.Regexp
	connections uint
	mu          sync.Mutex
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

func NewCollection(system, cpu, disk bool) map[MetricType]Scraper {
	m := make(map[MetricType]Scraper)

	if system {
		m[System] = &SystemScraper{}
	}
	if cpu {
		m[CPU] = &CPUScraper{}
	}
	if disk {
		m[Disk] = &DiskScraper{}
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
	for _, num := range numbers {
		sum += num
	}
	return sum / float64(len(numbers))
}

func parseLines[T float64 | int64](
	m map[string][][]T,
	data []string,
	valuesNum int,
	f func(data string) (T, error),
) (map[string][][]T, []string) {
	var deviceOrder []string

	for k, s := range data {
		lines := strings.Split(s, "\n")
		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) < valuesNum {
				continue
			}
			device := fields[0]
			if k == 0 {
				deviceOrder = append(deviceOrder, device)
			}
			var chunk []T
			for i := 1; i < len(fields); i++ {
				value, err := f(fields[i])
				if err != nil {
					continue
				}

				chunk = append(chunk, value)
			}
			m[device] = append(m[device], chunk)
		}
	}
	return m, deviceOrder
}

func calculateAvgMap[T float64 | int64](data map[string][][]T, fieldsCnt int) map[string][]T {
	m := make(map[string][]T)
	for device, chunks := range data {
		if len(chunks) == 0 {
			continue
		}
		averages := make([]T, fieldsCnt)
		for _, chunk := range chunks {
			for k, value := range chunk {
				averages[k] += value
			}
		}
		var result []T
		for _, sum := range averages {
			result = append(result, sum/T(len(chunks)))
		}
		m[device] = result
	}

	return m
}
