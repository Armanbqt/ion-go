package ion

import (
	"math/big"
	"reflect"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type ionEqual interface {
	eq(other interface{}) bool
}

type ionFloat struct{ float64 }
type ionInt struct {
	i32  int
	i64  int64
	ui64 uint64
	bi   *big.Int
}
type ionBool struct{ bool }
type ionString struct{ string }
type ionTimestamp struct{ time.Time }
type ionDecimal struct{ Decimal }

func (thisFloat ionFloat) eq(other interface{}) bool {
	return cmp.Equal(thisFloat.float64, other, cmpopts.EquateNaNs())
}

func (thisInt ionInt) eq(other interface{}) bool {
	switch val := other.(type) {
	case int:
		return cmp.Equal(thisInt.i32, val)
	case int64:
		return cmp.Equal(thisInt.i64, val)
	case uint64:
		return cmp.Equal(thisInt.ui64, val)
	case *big.Int:
		return thisInt.bi.Cmp(val) == 0
	default:
		return false
	}
}

func (thisBool ionBool) eq(other interface{}) bool {
	return cmp.Equal(thisBool.bool, other)
}

func (thisStr ionString) eq(other interface{}) bool {
	return cmp.Equal(thisStr.string, other)
}

func (thisTimestamp ionTimestamp) eq(other interface{}) bool {
	return thisTimestamp.Time.Equal(other.(time.Time))
}

func (thisDecimal ionDecimal) eq(other interface{}) bool {
	return thisDecimal.Decimal.Equal(other.(*Decimal))
}

func cmpAnnotations(thisAnnotations, otherAnnotations []string) bool {
	if len(thisAnnotations) != len(otherAnnotations) {
		return false
	}

	for idx, this := range thisAnnotations {
		other := otherAnnotations[idx]
		thisText, otherText := this, other

		if !cmp.Equal(thisText, otherText) {
			return false
		}
	}

	return true
}

func cmpValueSlices(thisValues, otherValues []interface{}) bool {
	if len(thisValues) == 0 && len(otherValues) == 0 {
		return true
	}

	if len(thisValues) != len(otherValues) {
		return false
	}

	res := false
	for idx, this := range thisValues {
		other := otherValues[idx]
		thisType, otherType := reflect.TypeOf(this), reflect.TypeOf(other)

		if thisType != otherType {
			return false
		}

		switch this.(type) {
		case string: // null.Sexp, null.List, null.Struct
			res = strNullTypeCmp(this, other)
		default:
			thisItem := this.(ionItem)
			otherItem := other.(ionItem)
			res = thisItem.equal(otherItem)
		}
		if !res {
			return false
		}
	}
	return res
}

func cmpFloats(thisValue, otherValue interface{}) bool {
	if !haveSameTypes(thisValue, otherValue) {
		return false
	}

	switch val := thisValue.(type) {
	case string: // null.float
		return strNullTypeCmp(val, otherValue)
	case float64:
		thisFloat := ionFloat{val}
		return thisFloat.eq(otherValue.(float64))
	default:
		return false
	}
}

func cmpInts(thisValue, otherValue interface{}) bool {
	if !haveSameTypes(thisValue, otherValue) {
		return false
	}

	switch val := thisValue.(type) {
	case string: // null.float
		return strNullTypeCmp(val, otherValue)
	case int:
		thisInt := ionInt{i32: val}
		return thisInt.eq(otherValue.(int))
	case int64:
		thisInt := ionInt{i64: val}
		return thisInt.eq(otherValue.(int64))
	case uint64:
		thisInt := ionInt{ui64: val}
		return thisInt.eq(otherValue.(uint64))
	case *big.Int:
		thisInt := ionInt{bi: val}
		return thisInt.eq(otherValue.(*big.Int))
	default:
		return false
	}
}

func cmpBools(thisValue, otherValue interface{}) bool {
	if !haveSameTypes(thisValue, otherValue) {
		return false
	}

	switch val := thisValue.(type) {
	case string: // null.bool
		return strNullTypeCmp(val, otherValue)
	case bool:
		thisBool := ionBool{val}
		return thisBool.eq(otherValue.(bool))
	default:
		return false
	}
}

func cmpStrings(thisValue, otherValue interface{}) bool {
	if !haveSameTypes(thisValue, otherValue) {
		return false
	}

	switch val := thisValue.(type) {
	case string:
		thisStr := ionString{val}
		return thisStr.eq(otherValue.(string))
	default:
		return false
	}
}

func cmpTimestamps(thisValue, otherValue interface{}) bool {
	if !haveSameTypes(thisValue, otherValue) {
		return false
	}

	switch val := thisValue.(type) {
	case string: // null.Timestamp
		return strNullTypeCmp(val, otherValue)
	case time.Time:
		thisTimestamp := ionTimestamp{val}
		return thisTimestamp.eq(otherValue)
	default:
		return false
	}
}

func cmpDecimals(thisValue, otherValue interface{}) bool {
	if !haveSameTypes(thisValue, otherValue) {
		return false
	}

	switch val := thisValue.(type) {
	case string: // null.bool
		return strNullTypeCmp(val, otherValue)
	case Decimal:
		thisDecimal := ionDecimal{val}
		return thisDecimal.eq(otherValue)
	default:
		return false
	}
}

func strNullTypeCmp(this, other interface{}) bool {
	thisStr := this.(string)
	otherStr := other.(string)
	return cmp.Equal(thisStr, otherStr)
}

func haveSameTypes(this, other interface{}) bool {
	return reflect.TypeOf(this) == reflect.TypeOf(other)
}
