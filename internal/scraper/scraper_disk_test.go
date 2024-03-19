package scraper

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiskScraper(t *testing.T) {
	t.Run("disk load snapshot average count success", func(t *testing.T) {
		data := []string{
			"nvme0n1p9 0.01 0.05 0.00\nnvme0n1p4 1.00 0.03 100.00\n",
			"nvme0n1p9 0.01 0.35 0.00\nnvme0n1p4 0.00 0.03 0.00\n",
			"nvme0n1p9 0.11 0.70 0.00\nnvme0n1p4 10.00 0.03 0.00\n",
		}

		cases := []struct {
			data     []string
			expected []float64
			seconds  int64
		}{
			{
				data:     data,
				expected: []float64{0.04, 0.37, 0.00, 3.67, 0.03, 33.33},
				seconds:  3,
			},
			{
				data:     data,
				expected: []float64{0.06, 0.52, 0.00, 5.00, 0.03, 0.00},
				seconds:  2,
			},
		}

		for _, tc := range cases {
			s := &DiskScraper{}
			s.dataLoad = append(s.dataLoad, tc.data...)

			snapshot := s.GetSnapshotLoad(tc.seconds)

			expected := fmt.Sprintf(
				"Disk Load Average\nDevice tps kB_read/s kB_wrtn/s\nnvme0n1p9 %.2f %.2f %.2f\nnvme0n1p4 %.2f %.2f %.2f\n",
				tc.expected[0], tc.expected[1], tc.expected[2], tc.expected[3], tc.expected[4], tc.expected[5])
			require.Equalf(t, expected, snapshot, "expected: %s, actual: %s", expected, snapshot)
		}
	})

	t.Run("disk usage snapshot average count success", func(t *testing.T) {
		data := []string{
			"/dev/nvme0n1p9 1000 50 100 20\n/dev/nvme0n1p4 1000 10 10 2\n",
			"/dev/nvme0n1p9 500 25 200 40\n/dev/nvme0n1p4 5000 50 0 0\n",
			"/dev/nvme0n1p9 0 0 300 50\n/dev/nvme0n1p4 10000 100 100 20\n",
		}

		cases := []struct {
			data     []string
			expected []int64
			seconds  int64
		}{
			{
				data:     data,
				expected: []int64{500, 25, 200, 36, 5333, 53, 36, 7},
				seconds:  3,
			},
			{
				data:     data,
				expected: []int64{250, 12, 250, 45, 7500, 75, 50, 10},
				seconds:  2,
			},
		}

		for _, tc := range cases {
			s := &DiskScraper{}
			s.dataUse = append(s.dataUse, tc.data...)

			snapshot := s.GetSnapshotUse(tc.seconds)

			expected := fmt.Sprintf(
				"Disk Usage Average\nFilesystem MB_Used Use%% IUsed IUse%%\n"+
					"/dev/nvme0n1p9 %d %d%% %d %d%%\n/dev/nvme0n1p4 %d %d%% %d %d%%\n",
				tc.expected[0], tc.expected[1], tc.expected[2], tc.expected[3],
				tc.expected[4], tc.expected[5], tc.expected[6], tc.expected[7])
			require.Equalf(t, expected, snapshot, "expected: %s, actual: %s", expected, snapshot)
		}
	})
}
