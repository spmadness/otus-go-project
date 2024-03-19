package scraper

import (
	"testing"
)

var casesCommon = []struct {
	data  string
	error bool
}{
	{
		data:  "1.83 2.00 1.00",
		error: false,
	},
	{
		data:  "1.83 2.00",
		error: true,
	},
	{
		data:  "",
		error: true,
	},
	{
		data:  "wrong value",
		error: true,
	},
}

func TestCpuScraper(t *testing.T) {
	dataCPU := []string{
		"5.0 0.0 0.0",
		"10.0 1.0 0.0",
		"15.0 2.0 10.0",
	}

	casesCPU := []struct {
		data     []string
		expected []float64
		seconds  int64
	}{
		{
			data:     dataCPU,
			expected: []float64{10.0, 1.0, 3.3},
			seconds:  3,
		},
		{
			data:     dataCPU,
			expected: []float64{12.5, 1.5, 5.0},
			seconds:  2,
		},
	}

	t.Run("parse data success", func(t *testing.T) {
		for _, tc := range casesCommon {
			s := &CPUScraper{}
			s.Init()
			parseTest(t, s, tc.data, tc.error)
		}
	})

	t.Run("snapshot average count success", func(t *testing.T) {
		for _, tc := range casesCPU {
			s := &CPUScraper{}
			s.data = append(s.data, tc.data...)

			snapshotTest(t, s, tc.seconds, tc.expected, s.GetSnapshotHeaders(), "%.1f")
		}
	})
}
