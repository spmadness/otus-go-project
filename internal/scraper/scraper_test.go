package scraper

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var casesCommon = []struct {
	name          string
	data          string
	error         bool
	expectedError error
}{
	{
		name:          "correct values amount",
		data:          "1.83 2.00 1.00",
		error:         false,
		expectedError: nil,
	},
	{
		name:          "correct values amount with prefix dot",
		data:          "0.0 .0.0 .100.0",
		error:         false,
		expectedError: nil,
	},
	{
		name:          "wrong values amount",
		data:          "1.83 2.00",
		error:         true,
		expectedError: ErrParseInitialData,
	},
	{
		name:          "empty values",
		data:          "",
		error:         true,
		expectedError: ErrParseInitialData,
	},
	{
		name:          "wrong values type",
		data:          "wrong value",
		error:         true,
		expectedError: ErrParseInitialData,
	},
}

func TestScraperSystemParseData(t *testing.T) {
	s := NewScraperSystem()
	for _, tc := range casesCommon {
		t.Run(tc.name, func(t *testing.T) {
			err := s.ParseData(tc.data, s.GetSnapshotDataRowElementsCnt())
			if tc.error {
				require.Error(t, err)
				require.Equalf(t, tc.expectedError, err, "Unexpected error: expected %v, actual %v", tc.expectedError, err)
				return
			}
			require.NoError(t, err)
			s.ClearData()
		})
	}
}

func TestScraperCPUParseData(t *testing.T) {
	s := NewScraperCPU()
	for _, tc := range casesCommon {
		t.Run(tc.name, func(t *testing.T) {
			err := s.ParseData(tc.data, s.GetSnapshotDataRowElementsCnt())
			if tc.error {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			s.ClearData()
		})
	}
}

func TestSystemScraperGetSnapshot(t *testing.T) {
	dataSystem := []string{
		"15.00 2.00 10.00",
		"10.00 1.00 0.00",
		"5.00 0.00 0.00",
	}

	casesSystem := []struct {
		name          string
		data          []string
		expected      []float64
		error         bool
		errorExpected error
		seconds       int64
	}{
		{
			name:          "system scraper get average from 3 last seconds",
			data:          dataSystem,
			expected:      []float64{10.00, 1.0, 3.33},
			error:         false,
			errorExpected: nil,
			seconds:       3,
		},
		{
			name:          "system scraper get average from 2 last seconds",
			data:          dataSystem,
			expected:      []float64{7.50, 0.50, 0.00},
			error:         false,
			errorExpected: nil,
			seconds:       2,
		},
		{
			name:          "system scraper non-positive seconds value",
			data:          dataSystem,
			expected:      []float64{},
			error:         true,
			errorExpected: ErrSecondsValue,
			seconds:       0,
		},
	}

	for _, tc := range casesSystem {
		t.Run(tc.name, func(t *testing.T) {
			s := NewScraperSystem()
			s.data = append(s.data, tc.data...)

			snapshotTest(t, s, tc.seconds, tc.expected, tc.error, tc.errorExpected)
		})
	}
}

func TestScraperCPUGetSnapshot(t *testing.T) {
	dataCPU := []string{
		"5.0 0.0 0.0",
		"10.0 1.0 0.0",
		"15.0 2.0 10.0",
	}

	casesCPU := []struct {
		name          string
		data          []string
		expected      []float64
		error         bool
		errorExpected error
		seconds       int64
	}{
		{
			name:          "cpu scraper non-positive seconds value",
			data:          dataCPU,
			expected:      []float64{},
			error:         true,
			errorExpected: ErrSecondsValue,
			seconds:       0,
		},
		{
			name:          "cpu scraper get average from 2 last seconds",
			data:          dataCPU,
			expected:      []float64{12.5, 1.5, 5.0},
			error:         false,
			errorExpected: nil,
			seconds:       2,
		},
		{
			name:          "cpu scraper get average from 3 last seconds",
			data:          dataCPU,
			expected:      []float64{10.0, 1.0, 3.3},
			error:         false,
			errorExpected: nil,
			seconds:       3,
		},
	}
	for _, tc := range casesCPU {
		t.Run(tc.name, func(t *testing.T) {
			s := NewScraperCPU()
			s.data = append(s.data, tc.data...)

			snapshotTest(t, s, tc.seconds, tc.expected, tc.error, tc.errorExpected)
		})
	}
}

func snapshotTest(t *testing.T, s Scraper, seconds int64, caseValues []float64, hasError bool, errorExpected error) {
	t.Helper()
	snapshot, err := s.GetSnapshot(seconds)
	if hasError {
		require.Equalf(t, errorExpected, err, "Unexpected error: expected %v, actual %v", errorExpected, err)
		return
	}

	formattedFloats := make([]string, len(caseValues))
	for i, val := range caseValues {
		formattedFloats[i] = fmt.Sprintf(s.GetSnapshotFormat(), val)
	}

	expected := fmt.Sprintf("%s%s", s.GetSnapshotHeaders(), strings.Join(formattedFloats, " "))
	require.Equalf(t, expected, snapshot[0], "expected: %s, actual: %s", expected, snapshot)
}

func TestGetLastN(t *testing.T) {
	cases := []struct {
		name     string
		data     []string
		n        int64
		expected []string
	}{
		{
			name:     "get last 3 from slice with 5 elements",
			data:     []string{"1.1", "2.2", "3.3", "4.4", "5.5"},
			n:        3,
			expected: []string{"3.3", "4.4", "5.5"},
		},
		{
			name:     "get last 5 from slice with 3 elements",
			data:     []string{"6.1", "7.2", "8.3"},
			n:        5,
			expected: []string{"6.1", "7.2", "8.3"},
		},
		{
			name:     "get last 2 from empty slice",
			data:     []string{},
			n:        2,
			expected: []string{},
		},
		{
			name:     "get last 0 from slice",
			data:     []string{"9.1", "10.2", "11.3"},
			n:        0,
			expected: []string{},
		},
		{
			name:     "get last 5 from slice with 5 elements",
			data:     []string{"12.1", "13.2", "14.3", "15.4", "16.5"},
			n:        5,
			expected: []string{"12.1", "13.2", "14.3", "15.4", "16.5"},
		},
		{
			name:     "get last 7 from slice with 5 elements",
			data:     []string{"17.1", "18.2", "19.3", "20.4", "21.5"},
			n:        7,
			expected: []string{"17.1", "18.2", "19.3", "20.4", "21.5"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := getLastN(tc.data, tc.n)
			require.Equalf(t, tc.expected, result, "expected %v, got %v", tc.expected, result)
		})
	}
}

func TestCalculateAverage(t *testing.T) {
	cases := []struct {
		name     string
		numbers  []float64
		expected float64
	}{
		{
			name:     "calculate average of positive numbers",
			numbers:  []float64{1.5, 2.5, 3.5, 4.5, 5.5},
			expected: 3.5,
		},
		{
			name:     "calculate average with negative numbers",
			numbers:  []float64{-10.5, -5.5, 0.5, 5.5, 10.5},
			expected: 0.1,
		},
		{
			name:     "calculate average of single number",
			numbers:  []float64{7.0},
			expected: 7.0,
		},
		{
			name:     "calculate average of empty slice",
			numbers:  []float64{},
			expected: 0.0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := calculateAverage(tc.numbers)
			require.Equalf(t, tc.expected, result, "expected %v, got %v", tc.expected, result)
		})
	}
}

func TestPrepareData(t *testing.T) {
	testCases := []struct {
		name          string
		data          []string
		seconds       int64
		dataRowCnt    int
		expected      [][]float64
		expectedError error
	}{
		{
			name:          "ValidData",
			data:          []string{"1 2 3", "4 5 6", "7 8 9"},
			seconds:       5,
			dataRowCnt:    3,
			expected:      [][]float64{{1, 4, 7}, {2, 5, 8}, {3, 6, 9}},
			expectedError: nil,
		},
		{
			name:          "EmptyData",
			data:          []string{},
			seconds:       5,
			dataRowCnt:    3,
			expected:      [][]float64{},
			expectedError: ErrEmptyScraperData,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := prepareData(tc.data, tc.seconds, tc.dataRowCnt)
			if err != nil {
				require.Equalf(t, tc.expectedError, err, "Unexpected error: expected %v, actual %v", tc.expectedError, err)
				return
			}
			require.Equalf(t, tc.expected, result, "actual %v, expected %v", result, tc.expected)
		})
	}
}

func TestResultString(t *testing.T) {
	testCases := []struct {
		name          string
		headers       string
		format        string
		values        [][]float64
		expected      []string
		expectedError error
	}{
		{
			name:          "valid data %.2f",
			headers:       "Averages: ",
			format:        "%.2f",
			values:        [][]float64{{1.10, 2.20, 3.30}, {4.40, 5.50, 6.60}, {7.70, 8.80, 9.90}},
			expected:      []string{"Averages: 2.20 5.50 8.80"},
			expectedError: nil,
		},
		{
			name:          "valid data %.1f",
			headers:       "Averages: ",
			format:        "%.1f",
			values:        [][]float64{{1.10, 2.20, 3.30}, {4.40, 5.50, 6.60}, {7.70, 8.80, 9.90}},
			expected:      []string{"Averages: 2.2 5.5 8.8"},
			expectedError: nil,
		},
		{
			name:          "empty data",
			headers:       "Averages: ",
			format:        "%.2f",
			values:        [][]float64{},
			expected:      []string{},
			expectedError: ErrStringOutputEmptyData,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := resultString(tc.headers, tc.format, tc.values...)
			if err != nil {
				require.Equalf(t, tc.expectedError, err, "Unexpected error: expected %v, actual %v", tc.expectedError, err)
				return
			}
			require.Equalf(t, tc.expected, result, "actual %v, expected %v", result, tc.expected)
		})
	}
}
