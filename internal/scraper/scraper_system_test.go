package scraper

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSystemScraper(t *testing.T) {
	dataSystem := []string{
		"15.00 2.00 10.00",
		"10.00 1.00 0.00",
		"5.00 0.00 0.00",
	}

	casesSystem := []struct {
		data     []string
		expected []float64
		seconds  int64
	}{
		{
			data:     dataSystem,
			expected: []float64{10.00, 1.0, 3.33},
			seconds:  3,
		},
		{
			data:     dataSystem,
			expected: []float64{7.50, 0.50, 0.00},
			seconds:  2,
		},
	}

	t.Run("parse data success", func(t *testing.T) {
		scrapers := []Scraper{&SystemScraper{}, &CPUScraper{}}

		for _, s := range scrapers {
			s.Init()

			for _, tc := range casesCommon {
				err := s.ParseData(tc.data)
				if tc.error {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)

				s.ClearData()
			}
		}
	})

	t.Run("snapshot average count success", func(t *testing.T) {
		for _, tc := range casesSystem {
			s := &SystemScraper{}
			s.data = append(s.data, tc.data...)

			snapshotTest(t, s, tc.seconds, tc.expected, s.GetSnapshotHeaders(), "%.2f")
		}
	})
}

func parseTest(t *testing.T, dp DataParse, data string, hasError bool) {
	t.Helper()
	err := dp.ParseData(data)
	if hasError {
		require.Error(t, err)
		return
	}
	require.NoError(t, err)
}

func snapshotTest(t *testing.T, s Scraper, seconds int64, caseValues []float64, headers string, formatValue string) {
	t.Helper()
	snapshot, err := s.GetSnapshot(seconds)
	if err != nil {
		t.Error("get snapshot error")
	}

	formattedFloats := make([]string, len(caseValues))
	for i, val := range caseValues {
		formattedFloats[i] = fmt.Sprintf(formatValue, val)
	}

	expected := fmt.Sprintf("%s%s", headers, strings.Join(formattedFloats, " "))
	require.Equalf(t, expected, snapshot[0], "expected: %s, actual: %s", expected, snapshot)
}
