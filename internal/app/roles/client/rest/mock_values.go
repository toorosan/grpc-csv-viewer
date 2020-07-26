package rest

import (
	"time"

	"grpc-csv-viewer/internal/app/roles/client"
)

func mockTimeSeries() client.TimeSeries {
	return client.TimeSeries{
		FileName:  "mocked-values.csv",
		StartDate: time.Now().Add(-time.Hour * 3),
		StopDate:  time.Now(),
		Values: []client.SeriesItem{
			{
				Date:  time.Now().Add(-time.Hour * 3),
				Value: 1,
			},
			{
				Date:  time.Now().Add(-time.Hour * 2),
				Value: 2,
			},
			{
				Date:  time.Now().Add(-time.Hour * 1),
				Value: 4,
			},
			{
				Date:  time.Now(),
				Value: 1,
			},
		},
	}
}
