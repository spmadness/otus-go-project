//go:build linux

package scraper

import "os/exec"

func (s *SystemScraper) command() *exec.Cmd {
	c := "top -bn 1 | sed -n 's/^.*load average: //p' | sed 's/, / /g' | sed 's/,/./g'"

	return exec.Command("bash", "-c", c)
}

func (s *CPUScraper) command() *exec.Cmd {
	c := "top -bn 1 | awk '/Cpu/' | grep -oE \"[0-9,\\.]+ (us|sy|id)\" | awk '{printf $1\" \"}' | sed 's/,/./g'"

	return exec.Command("bash", "-c", c)
}
