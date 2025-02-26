package convert

import (
	"github.com/apache/arrow-go/v18/arrow"
)

// ConvertArrowType converts an Arrow type from v17 to v18 format.
// This is needed because we have dependencies using v17 while we use v18.
func ConvertArrowType(dt arrow.DataType) arrow.DataType {
	// Both v17 and v18 implement the arrow.DataType interface
	// with the same methods, so we can use type assertions to
	// convert between them.
	switch dt := dt.(type) {
	case *arrow.BinaryType:
		return arrow.BinaryTypes.Binary
	case *arrow.StringType:
		return arrow.BinaryTypes.String
	case *arrow.Int64Type:
		return arrow.PrimitiveTypes.Int64
	case *arrow.Int32Type:
		return arrow.PrimitiveTypes.Int32
	case *arrow.Float64Type:
		return arrow.PrimitiveTypes.Float64
	case *arrow.BooleanType:
		return arrow.FixedWidthTypes.Boolean
	case *arrow.ListType:
		return arrow.ListOf(ConvertArrowType(dt.Elem()))
	case *arrow.StructType:
		fields := make([]arrow.Field, len(dt.Fields()))
		for i, f := range dt.Fields() {
			fields[i] = arrow.Field{
				Name:     f.Name,
				Type:     ConvertArrowType(f.Type),
				Nullable: f.Nullable,
			}
		}
		return arrow.StructOf(fields...)
	case *arrow.MapType:
		return arrow.MapOf(ConvertArrowType(dt.KeyType()), ConvertArrowType(dt.ItemType()))
	case *arrow.DictionaryType:
		return &arrow.DictionaryType{
			IndexType: ConvertArrowType(dt.IndexType).(arrow.FixedWidthDataType),
			ValueType: ConvertArrowType(dt.ValueType),
		}
	default:
		// For any other types, return as is since they should be compatible
		return dt
	}
}

// TypesEqual checks if two Arrow types are equal, handling v17 and v18 compatibility
func TypesEqual(a, b arrow.DataType) bool {
	// Convert both types to v18 format for comparison
	a = ConvertArrowType(a)
	b = ConvertArrowType(b)
	return arrow.TypeEqual(a, b)
}
