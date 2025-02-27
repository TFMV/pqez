# pqez

PQEZ is a high-performance query engine for Parquet files with dynamic schemas, built on top of Apache Arrow and Parquet. It provides a flexible and efficient way to query columnar data with support for dynamic columns.

## Features

- **Dynamic Columns**: Support for schemas with dynamic columns that can vary between datasets
- **Logical Query Planning**: Build complex queries using a logical plan builder
- **Arrow Integration**: Seamless integration with Apache Arrow for efficient in-memory processing
- **Compatibility Layer**: Handles compatibility between Arrow v17 and v18

## Architecture

PQEZ consists of several key components:

- **dynparquet**: Extends Parquet with support for dynamic columns
- **logicalplan**: Provides a query planning and optimization layer
- **pqarrow**: Handles conversion between Parquet and Arrow formats

## Getting Started

### Prerequisites

- Go 1.18 or later
- Apache Arrow v18

### Installation

```bash
go get github.com/TFMV/pqez
```

### Basic Usage

Here's a simple example of how to use PQEZ:

```go
package main

import (
 "context"
 "fmt"
 "log"

 "github.com/apache/arrow-go/v18/arrow/memory"
 "github.com/TFMV/pqez/dynparquet"
 "github.com/TFMV/pqez/logicalplan"
)

func main() {
 // Use the sample schema from dynparquet
 schema := dynparquet.NewSampleSchema()
 
 // Create a table provider
 tableProvider := &YourTableProvider{schema: schema}
 
 // Build a query
 builder := logicalplan.Builder{}
 plan, err := builder.
  Scan(tableProvider, "your_table").
  Project(
   logicalplan.Col("timestamp"),
   logicalplan.Col("value"),
   logicalplan.Col("labels"),
  ).
  Build()
  
 if err != nil {
  log.Fatalf("Failed to build logical plan: %v", err)
 }
}
```

## Query Capabilities

PQEZ supports a variety of query operations:

### Projection

Select specific columns from your data:

```go
plan, err := builder.
 Scan(tableProvider, "table").
 Project(
  logicalplan.Col("timestamp"),
  logicalplan.Col("value"),
 ).
 Build()
```

### Filtering

Filter data based on conditions:

```go
plan, err := builder.
 Scan(tableProvider, "table").
 Filter(
  logicalplan.Col("value").Gt(logicalplan.Literal(100)),
 ).
 Build()
```

### Aggregation

Perform aggregation operations:

```go
plan, err := builder.
 Scan(tableProvider, "table").
 Aggregate(
  []*logicalplan.AggregationFunction{
   logicalplan.Sum(logicalplan.Col("value")),
   logicalplan.Count(logicalplan.Col("value")),
  },
  []logicalplan.Expr{
   logicalplan.Col("category"),
  },
 ).
 Build()
```

### Dynamic Columns

Access dynamic columns using dot notation:

```go
plan, err := builder.
 Scan(tableProvider, "table").
 Project(
  logicalplan.Col("value"),
  logicalplan.Col("labels.environment"),
  logicalplan.Col("labels.region"),
 ).
 Build()
```

## Arrow Compatibility

PQEZ includes a compatibility layer to handle differences between Arrow v17 and v18. This allows it to work with dependencies that might be using different versions of Arrow.

The compatibility layer is implemented in `pqarrow/convert/compat.go` and provides functions for converting between Arrow types from different versions.

## Examples

Check out the `examples` directory for more detailed examples of how to use PQEZ:

- `demo.go`: Demonstrates basic query capabilities
- More examples coming soon!

## License

This project is licensed under the [MIT LICENSE](LICENSE).

## Acknowledgements

This repository is based on [frostdb](https://github.com/polarsignals/frostdb).
