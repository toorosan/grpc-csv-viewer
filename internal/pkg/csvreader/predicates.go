package csvreader

import (
	"fmt"
)

// FloatEqualityFunc describes interface for predicates to check float64 values equality.
type FloatEqualityFunc func(one, other float64) bool

// BasicFloatEquals checks if float numbers are completely equal.
func BasicFloatEquals(one, other float64) bool {
	return one == other
}

// CheckRoundedFloatEquals checks if float numbers rounded to the precision 2 are equal.
func CheckRoundedFloatEquals(one, other float64) bool {
	return fmt.Sprintf("%.2f", one) == fmt.Sprintf("%.2f", other)
}
