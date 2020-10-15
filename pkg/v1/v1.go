package v1

import (
	"fmt"

	"github.com/alandtsang/promquery/pkg/promquery"
)

// CPUQueryRange
func CPUQueryRange(query *promquery.PromQuery, params *promquery.QueryParams) {
	cpuUsageRange, err := query.CPUUsageQueryRange(params)
	if err != nil {
		return
	}
	fmt.Printf("=== cpuUsageRange=%+v\n", cpuUsageRange.Metrics)
}

// MemQueryRange
func MemQueryRange(query *promquery.PromQuery, params *promquery.QueryParams) {
	memUsageRange, err := query.MemUsageQueryRange(params)
	if err != nil {
		return
	}
	fmt.Printf("=== memUsageRange=%+v\n", memUsageRange)
}
