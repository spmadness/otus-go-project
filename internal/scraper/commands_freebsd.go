//go:build freebsd

package scraper

import "os/exec"

func (s *SystemScraper) command() *exec.Cmd {
	c := "top -bn 1 | sed -n 's/^.*load averages: //p' | sed 's/, / /g' | sed 's/,/./g' | awk '{print $1 FS $2 FS $3}'"

	return exec.Command("sh", "-c", c)
}

func (s *CPUScraper) command() *exec.Cmd {
	c := "top -bn 1 | awk '/CPU/ {gsub(\"%\",\"\"); print $2, $6, $10}' | head -1"

	return exec.Command("sh", "-c", c)
}
