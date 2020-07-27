package csvreader

import (
	"encoding/csv"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"grpc-csv-viewer/internal/pkg/logger"

	"github.com/pkg/errors"
)

// ParsedCSVFile describes structure of parsed data gathered from CSV file.
type ParsedCSVFile struct {
	ColumnNames struct {
		Date  string
		Value string
	}
	TimeSeries []TimeSeriesItem
}

// EqualTo allows to check whether passed parsed CSV file is equal to ours or not.
// Returns error with basic diff details or nil if objects are identical.
func (pcf *ParsedCSVFile) EqualTo(other *ParsedCSVFile, valueCheckPredicate FloatEqualityFunc) error {
	if pcf.ColumnNames.Value != other.ColumnNames.Value || pcf.ColumnNames.Date != other.ColumnNames.Date {
		return errors.New("column names are not equal")
	}
	if len(pcf.TimeSeries) != len(other.TimeSeries) {
		return errors.New("sizes of time series are not equal")
	}
	for i := range pcf.TimeSeries {
		if pcf.TimeSeries[i].Date.Unix() != other.TimeSeries[i].Date.Unix() {
			return errors.Errorf("time series are not equal: record %d contains different date: %q vs %q ", i, pcf.TimeSeries[i].Date, other.TimeSeries[i].Date)
		}
		if valueCheckPredicate == nil {
			valueCheckPredicate = BasicFloatEquals
		}
		if !valueCheckPredicate(pcf.TimeSeries[i].Value, other.TimeSeries[i].Value) {
			return errors.Errorf("time series are not equal: record %d contains different values: '%.2f' vs '%.2f' ", i, pcf.TimeSeries[i].Value, other.TimeSeries[i].Value)
		}
	}

	return nil
}

// TimeSeriesItem describes TimeSeries item gathered from CSV file.
type TimeSeriesItem struct {
	Date  time.Time
	Value float64
}

// ReadTimeSeriesFromCSV gathers time series from CSV file.
// Assuming first column is date, second column is float value.
func ReadTimeSeriesFromCSV(fileName string) (ParsedCSVFile, error) {
	fail := func(e error) (ParsedCSVFile, error) {
		return ParsedCSVFile{}, e
	}
	csvFile, err := os.Open(filepath.Clean(fileName))
	if err != nil {
		return fail(errors.Wrapf(err, "failed to open csv file %q", fileName))
	}
	defer func() {
		err = csvFile.Close()
		if err != nil {
			logger.Error(errors.Wrapf(err, "failed to close file %q", fileName).Error())
		}
	}()

	// Parse the file
	r := csv.NewReader(csvFile)

	result := ParsedCSVFile{}
	// Iterate through the records
	for i := 0; ; i++ {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fail(errors.Wrapf(err, "failed to read data from csv file %q", fileName))
		}
		if i == 0 {
			// Store column names from first row.
			result.ColumnNames.Date = record[0]
			result.ColumnNames.Value = record[1]

			continue
		}
		floatValue, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return fail(errors.Wrapf(err, "failed to convert value %q to float", record[1]))
		}
		timeValue, err := time.Parse("2006-01-02 15:04:05", record[0])
		if err != nil {
			return fail(errors.Wrapf(err, "failed to convert value %q to timestamp", record[0]))
		}
		result.TimeSeries = append(result.TimeSeries, TimeSeriesItem{
			Date:  timeValue,
			Value: floatValue,
		})
	}

	return result, nil
}
