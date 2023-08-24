package gconfig

import (
	"fmt"
	"time"
)

// Int64 returns the int64 value of a given key path or 0 if the path
// does not exist or if the value is not a valid int64.
func (c *configMgr) Int64(path string) int64 {
	return c.ko.Int64(path)
}

// MustInt64 returns the int64 value of a given key path or panics
// if the value is not set or set to default value of 0.
func (c *configMgr) MustInt64(path string) int64 {
	return c.ko.MustInt64(path)
}

// Int64s returns the []int64 slice value of a given key path or an
// empty []int64 slice if the path does not exist or if the value
// is not a valid int slice.
func (c *configMgr) Int64s(path string) []int64 {
	return c.ko.Int64s(path)
}

// MustInt64s returns the []int64 slice value of a given key path or panics
// if the value is not set or its default value.
func (c *configMgr) MustInt64s(path string) []int64 {
	return c.ko.MustInt64s(path)
}

// Int64Map returns the map[string]int64 value of a given key path
// or an empty map[string]int64 if the path does not exist or if the
// value is not a valid int64 map.
func (c *configMgr) Int64Map(path string) map[string]int64 {
	return c.ko.Int64Map(path)
}

// MustInt64Map returns the map[string]int64 value of a given key path
// or panics if its not set or set to default value.
func (c *configMgr) MustInt64Map(path string) map[string]int64 {
	return c.ko.MustInt64Map(path)
}

// Int returns the int value of a given key path or 0 if the path
// does not exist or if the value is not a valid int.
func (c *configMgr) Int(path string) int {
	return int(c.ko.Int64(path))
}

// MustInt returns the int value of a given key path or panics
// or panics if its not set or set to default value of 0.
func (c *configMgr) MustInt(path string) int {
	return c.ko.MustInt(path)
}

// Ints returns the []int slice value of a given key path or an
// empty []int slice if the path does not exist or if the value
// is not a valid int slice.
func (c *configMgr) Ints(path string) []int {
	return c.ko.Ints(path)
}

// MustInts returns the []int slice value of a given key path or panics
// if the value is not set or set to default value.
func (c *configMgr) MustInts(path string) []int {
	return c.ko.MustInts(path)
}

// IntMap returns the map[string]int value of a given key path
// or an empty map[string]int if the path does not exist or if the
// value is not a valid int map.
func (c *configMgr) IntMap(path string) map[string]int {
	return c.ko.IntMap(path)
}

// MustIntMap returns the map[string]int value of a given key path or panics
// if the value is not set or set to default value.
func (c *configMgr) MustIntMap(path string) map[string]int {
	return c.ko.MustIntMap(path)

}

// Float64 returns the float64 value of a given key path or 0 if the path
// does not exist or if the value is not a valid float64.
func (c *configMgr) Float64(path string) float64 {
	return c.ko.Float64(path)
}

// MustFloat64 returns the float64 value of a given key path or panics
// or panics if its not set or set to default value 0.
func (c *configMgr) MustFloat64(path string) float64 {
	return c.ko.MustFloat64(path)
}

// Float64s returns the []float64 slice value of a given key path or an
// empty []float64 slice if the path does not exist or if the value
// is not a valid float64 slice.
func (c *configMgr) Float64s(path string) []float64 {
	return c.ko.Float64s(path)
}

// MustFloat64s returns the []Float64 slice value of a given key path or panics
// if the value is not set or set to default value.
func (c *configMgr) MustFloat64s(path string) []float64 {
	return c.ko.MustFloat64s(path)
}

// Float64Map returns the map[string]float64 value of a given key path
// or an empty map[string]float64 if the path does not exist or if the
// value is not a valid float64 map.
func (c *configMgr) Float64Map(path string) map[string]float64 {
	return c.ko.Float64Map(path)
}

// MustFloat64Map returns the map[string]float64 value of a given key path or panics
// if the value is not set or set to default value.
func (c *configMgr) MustFloat64Map(path string) map[string]float64 {
	return c.ko.MustFloat64Map(path)
}

// Duration returns the time.Duration value of a given key path assuming
// that the key contains a valid numeric value.
func (c *configMgr) Duration(path string) time.Duration {
	return c.ko.Duration(path)
}

// MustDuration returns the time.Duration value of a given key path or panics
// if its not set or set to default value 0.
func (c *configMgr) MustDuration(path string) time.Duration {
	return c.ko.MustDuration(path)
}

// Time attempts to parse the value of a given key path and return time.Time
// representation. If the value is numeric, it is treated as a UNIX timestamp
// and if it's string, a parse is attempted with the given layout.
func (c *configMgr) Time(path, layout string) time.Time {
	// Unix timestamp?
	return c.ko.Time(path, layout)
}

// MustTime attempts to parse the value of a given key path and return time.Time
// representation. If the value is numeric, it is treated as a UNIX timestamp
// and if it's string, a parse is attempted with the given layout. It panics if
// the parsed time is zero.
func (c *configMgr) MustTime(path, layout string) time.Time {
	return c.ko.MustTime(path, layout)
}

// String returns the string value of a given key path or "" if the path
// does not exist or if the value is not a valid string.
func (c *configMgr) String(path string) string {
	return c.ko.String(path)
}

// MustString returns the string value of a given key path
// or panics if its not set or set to default value "".
func (c *configMgr) MustString(path string) string {
	return c.ko.MustString(path)
}

// Strings returns the []string slice value of a given key path or an
// empty []string slice if the path does not exist or if the value
// is not a valid string slice.
func (c *configMgr) Strings(path string) []string {
	return c.ko.Strings(path)
}

// MustStrings returns the []string slice value of a given key path or panics
// if the value is not set or set to default value.
func (c *configMgr) MustStrings(path string) []string {
	return c.ko.MustStrings(path)
}

// StringMap returns the map[string]string value of a given key path
// or an empty map[string]string if the path does not exist or if the
// value is not a valid string map.
func (c *configMgr) StringMap(path string) map[string]string {
	return c.ko.StringMap(path)
}

// MustStringMap returns the map[string]string value of a given key path or panics
// if the value is not set or set to default value.
func (c *configMgr) MustStringMap(path string) map[string]string {
	return c.ko.MustStringMap(path)
}

// StringsMap returns the map[string][]string value of a given key path
// or an empty map[string][]string if the path does not exist or if the
// value is not a valid strings map.
func (c *configMgr) StringsMap(path string) map[string][]string {
	return c.ko.StringsMap(path)
}

// MustStringsMap returns the map[string][]string value of a given key path or panics
// if the value is not set or set to default value.
func (c *configMgr) MustStringsMap(path string) map[string][]string {
	return c.ko.MustStringsMap(path)
}

// Bytes returns the []byte value of a given key path or an empty
// []byte slice if the path does not exist or if the value is not a valid string.
func (c *configMgr) Bytes(path string) []byte {
	return c.ko.Bytes(path)
}

// MustBytes returns the []byte value of a given key path or panics
// if the value is not set or set to default value.
func (c *configMgr) MustBytes(path string) []byte {
	return c.ko.MustBytes(path)
}

// Bool returns the bool value of a given key path or false if the path
// does not exist or if the value is not a valid bool representation.
// Accepted string representations of bool are the ones supported by strconv.ParseBool.
func (c *configMgr) Bool(path string) bool {
	return c.ko.Bool(path)
}

// Bools returns the []bool slice value of a given key path or an
// empty []bool slice if the path does not exist or if the value
// is not a valid bool slice.
func (c *configMgr) Bools(path string) []bool {
	return c.ko.Bools(path)
}

// MustBools returns the []bool value of a given key path or panics
// if the value is not set or set to default value.
func (c *configMgr) MustBools(path string) []bool {
	return c.ko.MustBools(path)
}

// BoolMap returns the map[string]bool value of a given key path
// or an empty map[string]bool if the path does not exist or if the
// value is not a valid bool map.
func (c *configMgr) BoolMap(path string) map[string]bool {
	return c.ko.BoolMap(path)
}

// MustBoolMap returns the map[string]bool value of a given key path or panics
// if the value is not set or set to default value.
func (c *configMgr) MustBoolMap(path string) map[string]bool {
	val := c.ko.MustBoolMap(path)
	if len(val) == 0 {
		panic(fmt.Sprintf("invalid value: %s=%v", path, val))
	}
	return val
}
