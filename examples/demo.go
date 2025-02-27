package main

import (
	"context"
	"fmt"
	"log"

	"github.com/apache/arrow-go/v18/arrow/memory"

	"github.com/TFMV/pqez/dynparquet"
	"github.com/TFMV/pqez/logicalplan"
)

// This demo demonstrates the core capabilities of the pqez query engine:
// 1. Using the sample schema from dynparquet
// 2. Creating a logical plan for querying
// 3. Showing how to build different types of queries

func main() {
	// Use the sample schema from dynparquet
	schema := dynparquet.NewSampleSchema()
	fmt.Println("Using sample schema:")
	fmt.Println(schema)

	// Create a mock table provider that uses the sample schema
	tableProvider := &MockTableProvider{schema: schema}

	// Demo different types of queries
	fmt.Println("\n--- Simple Query Example ---")
	demonstrateSimpleQuery(tableProvider)

	fmt.Println("\n--- Aggregation Query Example ---")
	demonstrateAggregationQuery(tableProvider)

	fmt.Println("\n--- Filtered Query Example ---")
	demonstrateFilteredQuery(tableProvider)

	fmt.Println("\n--- Dynamic Columns Example ---")
	demonstrateDynamicColumnsQuery(tableProvider)
}

// MockTableProvider implements the logicalplan.TableProvider interface
type MockTableProvider struct {
	schema *dynparquet.Schema
}

// GetTable returns a mock table reader
func (p *MockTableProvider) GetTable(name string) (logicalplan.TableReader, error) {
	return &MockTableReader{schema: p.schema}, nil
}

// MockTableReader implements the logicalplan.TableReader interface
type MockTableReader struct {
	schema *dynparquet.Schema
}

// View implements the TableReader interface
func (r *MockTableReader) View(ctx context.Context, fn func(ctx context.Context, tx uint64) error) error {
	return fn(ctx, 0)
}

// Iterator implements the TableReader interface
func (r *MockTableReader) Iterator(
	ctx context.Context,
	tx uint64,
	pool memory.Allocator,
	callbacks []logicalplan.Callback,
	options ...logicalplan.Option,
) error {
	// This is a mock implementation that doesn't actually return data
	// In a real implementation, this would read data and call the callbacks
	fmt.Println("Iterator called with options:", formatOptions(options))
	return nil
}

// SchemaIterator implements the TableReader interface
func (r *MockTableReader) SchemaIterator(
	ctx context.Context,
	tx uint64,
	pool memory.Allocator,
	callbacks []logicalplan.Callback,
	options ...logicalplan.Option,
) error {
	// This is a mock implementation that doesn't actually return data
	fmt.Println("SchemaIterator called with options:", formatOptions(options))
	return nil
}

// Schema returns the schema of the table
func (r *MockTableReader) Schema() *dynparquet.Schema {
	return r.schema
}

// Helper function to format options for display
func formatOptions(options []logicalplan.Option) string {
	opts := logicalplan.IterOptions{}
	for _, opt := range options {
		opt(&opts)
	}

	result := ""
	if opts.Filter != nil {
		result += fmt.Sprintf("Filter: %v, ", opts.Filter)
	}
	if len(opts.Projection) > 0 {
		result += fmt.Sprintf("Projection: %v, ", opts.Projection)
	}
	if len(opts.PhysicalProjection) > 0 {
		result += fmt.Sprintf("PhysicalProjection: %v, ", opts.PhysicalProjection)
	}
	if len(opts.DistinctColumns) > 0 {
		result += fmt.Sprintf("DistinctColumns: %v, ", opts.DistinctColumns)
	}
	if result == "" {
		return "none"
	}
	return result
}

// demonstrateSimpleQuery shows how to build a simple query
func demonstrateSimpleQuery(tableProvider *MockTableProvider) {
	// Build the logical plan
	builder := logicalplan.Builder{}
	plan, err := builder.
		Scan(tableProvider, "sample_table").
		Project(
			logicalplan.Col("timestamp"),
			logicalplan.Col("value"),
			logicalplan.Col("labels"),
		).
		Build()

	if err != nil {
		log.Fatalf("Failed to build logical plan: %v", err)
	}

	// Print the logical plan
	fmt.Println("Logical Plan:")
	fmt.Println(plan)
}

// demonstrateAggregationQuery shows how to build an aggregation query
func demonstrateAggregationQuery(tableProvider *MockTableProvider) {
	// Build the logical plan
	builder := logicalplan.Builder{}
	plan, err := builder.
		Scan(tableProvider, "sample_table").
		Aggregate(
			[]*logicalplan.AggregationFunction{
				logicalplan.Sum(logicalplan.Col("value")),
				logicalplan.Count(logicalplan.Col("value")),
				logicalplan.Avg(logicalplan.Col("value")),
			},
			[]logicalplan.Expr{
				logicalplan.Col("labels"),
			},
		).
		Build()

	if err != nil {
		log.Fatalf("Failed to build logical plan: %v", err)
	}

	// Print the logical plan
	fmt.Println("Logical Plan:")
	fmt.Println(plan)
}

// demonstrateFilteredQuery shows how to build a query with a filter
func demonstrateFilteredQuery(tableProvider *MockTableProvider) {
	// Build the logical plan
	builder := logicalplan.Builder{}
	plan, err := builder.
		Scan(tableProvider, "sample_table").
		Filter(
			logicalplan.Col("labels").Eq(logicalplan.Literal("cpu")),
		).
		Project(
			logicalplan.Col("timestamp"),
			logicalplan.Col("value"),
			logicalplan.Col("labels"),
		).
		Build()

	if err != nil {
		log.Fatalf("Failed to build logical plan: %v", err)
	}

	// Print the logical plan
	fmt.Println("Logical Plan:")
	fmt.Println(plan)
}

// demonstrateDynamicColumnsQuery shows how to build a query with dynamic columns
func demonstrateDynamicColumnsQuery(tableProvider *MockTableProvider) {
	// Build the logical plan
	builder := logicalplan.Builder{}
	plan, err := builder.
		Scan(tableProvider, "sample_table").
		Project(
			logicalplan.Col("value"),
			logicalplan.Col("labels"),
			// Access dynamic columns with dot notation
			logicalplan.Col("labels.environment"),
			logicalplan.Col("labels.region"),
			logicalplan.Col("labels.service"),
		).
		Build()

	if err != nil {
		log.Fatalf("Failed to build logical plan: %v", err)
	}

	// Print the logical plan
	fmt.Println("Logical Plan:")
	fmt.Println(plan)
}
