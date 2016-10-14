package conditions

import (
	"math"

	"github.com/grafana/grafana/pkg/tsdb"
	"gopkg.in/guregu/null.v3"
)

type QueryReducer interface {
	Reduce(timeSeries *tsdb.TimeSeries) null.Float
}

type SimpleReducer struct {
	Type string
}

func (s *SimpleReducer) Reduce(series *tsdb.TimeSeries) null.Float {
	if len(series.Points) == 0 {
		return null.FloatFromPtr(nil)
	}

	value := float64(0)
	allNull := true

	switch s.Type {
	case "avg":
		for _, point := range series.Points {
			if point[0].Valid {
				value += point[0].Float64
				allNull = false
			}
		}
		value = value / float64(len(series.Points))
	case "sum":
		for _, point := range series.Points {
			if point[0].Valid {
				value += point[0].Float64
				allNull = false
			}
		}
	case "min":
		value = math.MaxFloat64
		for _, point := range series.Points {
			if point[0].Valid {
				allNull = false
				if value > point[0].Float64 {
					value = point[0].Float64
				}
			}
		}
	case "max":
		value = -math.MaxFloat64
		for _, point := range series.Points {
			if point[0].Valid {
				allNull = false
				if value < point[0].Float64 {
					value = point[0].Float64
				}
			}
		}
	case "count":
		value = float64(len(series.Points))
		allNull = false
	}

	if allNull {
		return null.FloatFromPtr(nil)
	}

	return null.FloatFrom(value)
}

func NewSimpleReducer(typ string) *SimpleReducer {
	return &SimpleReducer{Type: typ}
}
