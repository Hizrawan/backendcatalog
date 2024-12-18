package factories

import (
	"errors"
	"gopkg.in/guregu/null.v4"
	"time"
)

// mapContains checks whether the map contains the provided key, regardless
// whether the value is a zero-value or nil.
func mapContains(m map[string]any, key string) bool {
	if m == nil {
		return false
	}
	_, ok := m[key]
	return ok
}

// getOrDefault returns the value in the map for a given key if exists.
// Otherwise, it will return the provided default value.
func getOrDefault(m map[string]any, key string, def any) any {
	if mapContains(m, key) {
		return m[key]
	}
	return def
}

// getOrDefaultFunction returns the value in the map for a given key if exists.
// Otherwise, it will call the default function and returns it.
func getOrDefaultFunction(m map[string]any, key string, def func(map[string]any, string) any) any {
	if mapContains(m, key) {
		return m[key]
	}
	return def(m, key)
}

// getOrDefaultWithTransform returns the value in the map for a given key if it
// exists. Before returning, it will call the transform function, converting the
// original value. Otherwise, it will return the provided default value.
//func getOrDefaultWithTransform(m map[string]any, key string, transform func(any) any, def any) any {
//	if mapContains(m, key) {
//		return transform(m[key])
//	}
//	return def
//}
//
// getOrDefaultFunctionWithTransform returns the value in the map for a given key
// if it exists. Before returning, it will call the transform function,
// converting the original value. Otherwise, it will call the default function
// and returns it.
//func getOrDefaultFunctionWithTransform(m map[string]any, key string, transform func(any) any, def func(map[string]any, string) any) any {
//	if mapContains(m, key) {
//		return transform(m[key])
//	}
//	return def(m, key)
//}

// getOptionMap returns the value of the specified field from the map as
// map[string]any
func getOptionMap(m map[string]any, field string) map[string]any {
	var opt map[string]any
	if mapContains(m, field) {
		opt = m[field].(map[string]any)
	}
	return opt
}

// getAsTime returns the value of the specified key from the map as a null.Time
// instance. If the value in the map originally is a string, it will be parsed
// using time.RFC3339 (ISO8601) format before being returned. If the value is
// originally a time.Time instance, it will be wrapped as null.Time instance.
func getAsTime(m map[string]any, key string) null.Time {
	t := null.Time{}
	if mapContains(m, key) {
		t = toNullTime(m[key])
	}
	return t
}

// toNullTime converts the provided value into a null.Time instance if it's
// compatible. This function can only convert value from type string, time.Time,
// or null.Time. If the value is of other types, this function will panic.
func toNullTime(value any) null.Time {
	t := null.Time{}
	if value != nil {
		switch v := value.(type) {
		case string:
			parsed, err := time.Parse(time.RFC3339, v)
			if err != nil {
				panic(err)
			}
			t = null.TimeFrom(parsed)
		case time.Time:
			t = null.TimeFrom(v)
		case null.Time:
			t = v
		default:
			panic(errors.New("value must be a time.Time, a null.Time, or a RFC3339 (ISO8601) formatted string"))
		}
	}
	return t
}
