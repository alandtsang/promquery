package promquery

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	promapi "github.com/prometheus/client_golang/api"
	promv1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

type queryType int

const (
	// query range
	cpuUsageQueryRange queryType = iota
	memUsageQueryRange
)

type QueryParams struct {
	node       string
	namespace  string
	pod        string
	container  string
	queryRange string
}

type PromQuery struct {
	client promapi.Client
	api    promv1.API
	addr   string
}

func New(addr string) (*PromQuery, error) {
	client, err := newPrometheusClient(addr)
	if err != nil {
		return nil, err
	}

	return &PromQuery{
		client: client,
		api:    promv1.NewAPI(client),
		addr:   addr,
	}, nil
}

func newPrometheusClient(addr string) (promapi.Client, error) {
	client, err := promapi.NewClient(promapi.Config{
		Address: addr,
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return nil, err
	}
	return client, nil
}

func (pq *PromQuery) CPUUsageQueryRange(params *QueryParams) (QueryResult, error) {
	expr, err := buildPromQueryRangeExpr(cpuUsageQueryRange, params)
	if err != nil {
		fmt.Println("CPUUsageQueryRange error", err)
		return QueryResult{}, err
	}
	return pq.queryRange(expr)
}

func (pq *PromQuery) MemUsageQueryRange(params *QueryParams) (QueryResult, error) {
	expr, err := buildPromQueryRangeExpr(memUsageQueryRange, params)
	if err != nil {
		fmt.Println("MemUsageQueryRange error", err)
		return QueryResult{}, err
	}
	return pq.queryRange(expr)
}

func buildPromQueryRangeExpr(t queryType, params *QueryParams) (string, error) {
	var queryBase string

	switch t {
	case cpuUsageQueryRange:
		queryBase = "sum(rate(container_cpu_usage_seconds_total{node=\"%s\", namespace=\"%s\", pod=\"%s\", container=~\"%s\"}[%s])) by (pod)"
	case memUsageQueryRange:
		queryBase = "rate(container_memory_usage_bytes{node=\"%s\", namespace=\"%s\", pod=\"%s\", container=\"%s\"}[%s])"
	default:
		return "", errors.New("invalid query type")
	}

	return fmt.Sprintf(queryBase, params.node, params.namespace, params.pod, params.container, params.queryRange), nil
}

func (pq *PromQuery) query(expr string) (QueryResult, error) {
	var parseResult QueryResult

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, warnings, err := pq.api.Query(ctx, expr, time.Now())
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		return parseResult, err
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}

	fmt.Printf("result=%+v", result)

	parseResult, err = parseResponse(result)
	if err != nil {
		fmt.Println("query error,", err)
		return parseResult, err
	}

	return parseResult, nil
}

func (pq *PromQuery) queryRange(expr string) (QueryResult, error) {
	var parseResult QueryResult

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	r := promv1.Range{
		Start: time.Now().Add(-5 * time.Minute),
		End:   time.Now(),
		Step:  time.Minute,
	}

	log.Print("expr=", expr)
	result, warnings, err := pq.api.QueryRange(ctx, expr, r)
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		return parseResult, err
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}

	log.Printf("result=%+v", result)

	parseResult, err = parseResponse(result)
	if err != nil {
		return parseResult, err
	}

	return parseResult, nil
}

// BuildQueryParams builds query PromQuery need parameters.
func BuildQueryParams(node, namespace, pod, container, queryRange string) *QueryParams {
	return &QueryParams{
		node:       node,
		namespace:  namespace,
		pod:        pod,
		container:  container,
		queryRange: queryRange,
	}
}
