package scraper

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

type DiskScraper struct {
	dataLoad []string
	dataUse  []string
	BaseScraper
}

func (s *DiskScraper) ParseData(_ string) error {
	return nil
}

func (s *DiskScraper) Init() {
}

func (s *DiskScraper) GetData() {
	var wg sync.WaitGroup
	wg.Add(2)

	go s.GetDataLoad(&wg)
	go s.GetDataUse(&wg)

	wg.Wait()
}

func (s *DiskScraper) GetSnapshot(seconds int64) ([]string, error) {
	var result []string

	result = append(result, s.GetSnapshotLoad(seconds))
	result = append(result, s.GetSnapshotUse(seconds))

	return result, nil
}

func (s *DiskScraper) GetSnapshotLoad(seconds int64) string {
	var sb strings.Builder
	sb.WriteString("Disk Load Average\n")
	sb.WriteString("Device tps kB_read/s kB_wrtn/s\n")
	data := getLastN(s.dataLoad, seconds)

	m := make(map[string][][]float64)

	values, order := parseLines(m, data, 4, func(data string) (float64, error) {
		return strconv.ParseFloat(data, 64)
	})
	averages := calculateAvgMap(values, 3)

	for _, device := range order {
		e := averages[device]
		str := fmt.Sprintf("%s %.2f %.2f %.2f\n", device, e[0], e[1], e[2])
		sb.WriteString(str)
	}
	return sb.String()
}

func (s *DiskScraper) GetSnapshotUse(seconds int64) string {
	var sb strings.Builder
	sb.WriteString("Disk Usage Average\n")
	sb.WriteString("Filesystem MB_Used Use% IUsed IUse%\n")
	data := getLastN(s.dataUse, seconds)

	m := make(map[string][][]int64)

	values, order := parseLines(m, data, 5, func(data string) (int64, error) {
		return strconv.ParseInt(data, 10, 64)
	})
	averages := calculateAvgMap(values, 4)

	for _, device := range order {
		e := averages[device]
		str := fmt.Sprintf("%s %d %d%% %d %d%%\n", device, e[0], e[1], e[2], e[3])
		sb.WriteString(str)
	}
	return sb.String()
}

func (s *DiskScraper) GetDataLoad(wg *sync.WaitGroup) {
	defer wg.Done()

	s.mu.Lock()
	defer s.mu.Unlock()

	c := "iostat -d -k | grep 'Device' -A1000 | grep -v '^$' | " +
		"awk '{ print $1 FS $2 FS $3 FS $4 }' | sed 's/,/./g' | tail -n +2"
	cmd := exec.Command("bash", "-c", c)

	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}

	s.dataLoad = append(s.dataLoad, string(out))
}

func (s *DiskScraper) GetDataUse(wg *sync.WaitGroup) {
	defer wg.Done()

	s.mu.Lock()
	defer s.mu.Unlock()

	c := "paste <(df -m | awk '{print $1 FS $3 FS $5}') <(df -i | awk '{print $3 FS $5}') | tail -n +2 | sed 's/%//g'"
	cmd := exec.Command("bash", "-c", c)

	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}

	s.dataUse = append(s.dataUse, string(out))
}

func (s *DiskScraper) ClearData() {
	s.dataUse = []string{}
	s.dataLoad = []string{}
}

func (s *DiskScraper) GetSnapshotHeaders() string {
	return ""
}
