package ion

import (
	"math/big"
	"reflect"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type ionEqual interface {
	eq(other ionEqual) bool
}

type ionFloat float64
type ionInt int
type ionInt64 int64
type ionUint64 uint64
type ionBigInt struct{ *big.Int }
type ionBool bool
type ionString string
type ionTimestamp struct{ time.Time }
type ionDecimal struct{ *Decimal }
type ionLob []byte

func (thisFloat ionFloat) eq(other ionEqual) bool {
	return cmp.Equal(thisFloat, other, cmpopts.EquateNaNs())
}

func (thisInt ionInt) eq(other ionEqual) bool {
	return cmp.Equal(thisInt, other)
}

func (thisInt64 ionInt64) eq(other ionEqual) bool {
	return cmp.Equal(thisInt64, other)
}

func (thisUint64 ionUint64) eq(other ionEqual) bool {
	return cmp.Equal(thisUint64, other)
}

func (thisBigInt ionBigInt) eq(other ionEqual) bool {
	if val, ok := other.(ionBigInt); ok {
		return thisBigInt.Int.Cmp(val.Int) == 0
	}
	return false
}

func (thisBool ionBool) eq(other ionEqual) bool {
	return cmp.Equal(thisBool, other)
}

func (thisStr ionString) eq(other ionEqual) bool {
	return cmp.Equal(thisStr, other)
}

func (thisTimestamp ionTimestamp) eq(other ionEqual) bool {
	if val, ok := other.(ionTimestamp); ok {
		return thisTimestamp.Time.Equal(val.Time)
	}
	return false
}

func (thisDecimal ionDecimal) eq(other ionEqual) bool {
	if val, ok := other.(ionDecimal); ok {
		return thisDecimal.Decimal.Equal(val.Decimal)
	}
	return false
}

func (thisLob ionLob) eq(other ionEqual) bool {
	return cmp.Equal(thisLob, other)
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
		thisFloat := ionFloat(val)
		return thisFloat.eq(ionFloat(otherValue.(float64)))
	default:
		return false
	}
}

func cmpInts(thisValue, otherValue interface{}) bool {
	if !haveSameTypes(thisValue, otherValue) {
		return false
	}

	switch val := thisValue.(type) {
	case string: // null.int
		return strNullTypeCmp(val, otherValue)
	case int:
		thisInt := ionInt(val)
		return thisInt.eq(ionInt(otherValue.(int)))
	case int64:
		thisInt := ionInt64(val)
		return thisInt.eq(ionInt64(otherValue.(int64)))
	case uint64:
		thisInt := ionUint64(val)
		return thisInt.eq(ionUint64(otherValue.(uint64)))
	case *big.Int:
		thisInt := ionBigInt{val}
		return thisInt.eq(ionBigInt{otherValue.(*big.Int)})
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
		thisBool := ionBool(val)
		return thisBool.eq(ionBool(otherValue.(bool)))
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
		thisStr := ionString(val)
		return thisStr.eq(ionString(otherValue.(string)))
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
		return thisTimestamp.eq(ionTimestamp{otherValue.(time.Time)})
	default:
		return false
	}
}

func cmpDecimals(thisValue, otherValue interface{}) bool {
	if !haveSameTypes(thisValue, otherValue) {
		return false
	}

	switch val := thisValue.(type) {
	case string: // null.decimal
		return strNullTypeCmp(val, otherValue)
	case *Decimal:
		thisDecimal := ionDecimal{val}
		return thisDecimal.eq(ionDecimal{otherValue.(*Decimal)})
	default:
		return false
	}
}

func cmpLobs(thisValue, otherValue interface{}) bool {
	if !haveSameTypes(thisValue, otherValue) {
		return false
	}

	switch val := thisValue.(type) {
	case string: // null.blob  null.clob
		return strNullTypeCmp(val, otherValue)
	case []byte:
		thisLob := ionLob(val)
		return thisLob.eq(ionLob(otherValue.([]byte)))
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
