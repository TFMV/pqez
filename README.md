# pqez

This repository is based on [frostdb](https://github.com/polarsignals/frostdb), with modifications to handle Arrow version compatibility issues.

## Overview

pqez is a high-performance query engine that combines the power of Apache Arrow and Apache Parquet. It provides efficient data processing capabilities while handling compatibility between different versions of Arrow.

## Arrow Version Compatibility

The project uses Arrow v18 (`github.com/apache/arrow-go/v18`) as its primary Arrow dependency, while maintaining compatibility with dependencies that use Arrow v17. This is achieved through a compatibility layer that handles type conversions between the two versions.

### Compatibility Layer

The compatibility layer is implemented in `pqarrow/convert/compat.go` and provides:

1. `ConvertArrowType`: Converts Arrow types between v17 and v18 formats
2. `TypesEqual`: Compares Arrow types across versions

Supported types include:

- Binary
- String
- Int64/Int32
- Float64
- Boolean
- List
- Struct
- Map
- Dictionary

## Usage

When working with Arrow types in the codebase:

```go
import (
    "github.com/TFMV/pqez/pqarrow/convert"
    "github.com/apache/arrow-go/v18/arrow"
)

// Compare types using the compatibility layer
if convert.TypesEqual(type1, type2) {
    // Types are equal regardless of Arrow version
}

// Convert types to v18 format
v18Type := convert.ConvertArrowType(someType)
```

## Dependencies

- `github.com/apache/arrow-go/v18` v18.0.0 (primary Arrow dependency)
- `github.com/polarsignals/frostdb` (uses Arrow v17 internally)

## Contributing

When adding new features or modifying existing code:

1. Use Arrow v18 types in new code
2. Use the compatibility layer when comparing or converting types
3. Add tests for any new type conversions
4. Ensure all tests pass with both Arrow v17 and v18 dependencies

## License

This project maintains the same license as the original frostdb repository.
