package promquery

import (
	"fmt"
	"github.com/prometheus/common/model"
)

type QueryResult struct {
	Metrics map[string]string `json:"metrics"`
	Data    []DataPair        `json:"data"`
}

type DataPair struct {
	Value     string `json:"value"`
	Timestamp string `json:"timestamp"`
}

func parseResponse(value model.Value) (QueryResult, error) {
	var result QueryResult

	fmt.Printf("\nvalue model:%+v\n", value.Type())

	data, ok := value.(model.Matrix)
	if !ok {
		return result, fmt.Errorf("unsupported query result format")
	}

	result.Metrics = make(map[string]string, len(data))
	result.Data = make([]DataPair, 0, len(data))

	fmt.Println("data=", len(data))

	for _, v := range data {
		for k, v := range v.Metric {
			result.Metrics[string(k)] = string(v)
		}

		for _, k := range v.Values {
			result.Data = append(result.Data, DataPair{Value: k.Value.String(), Timestamp: k.Timestamp.String()})
		}
	}

	return result, nil
}
