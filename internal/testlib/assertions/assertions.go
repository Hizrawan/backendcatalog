package assertions

import (
	"fmt"
	"net/http"
	"time"

	"github.com/onsi/gomega/format"

	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/onsi/gomega/types"
	"github.com/xinchuantw/hoki-tabloid-backend/internal/testlib"
	"gopkg.in/guregu/null.v4"
)

func ReturnsValidationError(fields []string) types.GomegaMatcher {
	return And(
		HaveHTTPStatus(http.StatusUnprocessableEntity),
		WithTransform(
			testlib.MapFromResponseBody,
			And(
				MatchKeys(IgnoreExtras, Keys{
					"error_code": Equal("validation_error"),
				}),
				WithTransform(func(body map[string]any) []string {
					fieldErrs := body["fields"].(map[string]any)
					fields := getChildKeys(fieldErrs, "")
					return fields
				}, ContainElements(fields)),
			),
		),
	)
}

func getChildKeys(child map[string]any, prefix string) []string {
	var fields []string
	for k, v := range child {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}
		fields = append(fields, key)
		switch parsed := v.(type) {
		case map[string]any:
			childFields := getChildKeys(parsed, key)
			if len(childFields) > 0 {
				fields = append(fields, childFields...)
			}
		}
	}

	return fields
}

func MatchTime(t null.Time, exact bool, accuracy time.Duration) types.GomegaMatcher {
	return &TimeMatcher{
		Expected:   t,
		Exact:      exact,
		TruncateBy: accuracy,
	}
}

// BeFormattedTimeOf checks whether the string provided matches the expected time, up
// until the specified point of accuracy.
//
// This functions consider formatted time with different timezone as equal as long as
// both the expected and provided time refer to the same moment in time. For example,
// the expected time could be in UTC+08:00 while the formatted form is in UTC.
//
// An accuracy value of 0 means checking up until nanosecond precision. If given a
// time.Second accuracy, it will only check that the time is equivalent up to the
// second, ignoring the sub-second components. This is done by using time.Truncate
// before doing the comparison.
func BeFormattedTimeOf(t null.Time, accuracy time.Duration) types.GomegaMatcher {
	return &TimeMatcher{
		Expected:   t,
		Exact:      false,
		TruncateBy: accuracy,
	}
}

func BeFormattedTimeAndTimezoneOf(t null.Time, accuracy time.Duration) types.GomegaMatcher {
	return &TimeMatcher{
		Expected:   t,
		Exact:      true,
		TruncateBy: accuracy,
	}
}

type TimeMatcher struct {
	Expected   null.Time
	Exact      bool
	TruncateBy time.Duration
}

func (matcher *TimeMatcher) Match(actual interface{}) (success bool, err error) {
	if actual == nil && !matcher.Expected.Valid {
		return true, nil
	}

	var actualAsTime *time.Time
	formats := []string{
		time.RFC850,
		time.RFC3339,
		time.RFC3339Nano,
		time.ANSIC,
	}
	for _, format := range formats {
		t, err := time.Parse(format, actual.(string))
		if err == nil {
			actualAsTime = &t
			break
		}
	}

	if actualAsTime == nil {
		return false, fmt.Errorf("%v cannot be deserialized into time", actual)
	}

	actualAsNullTime := null.TimeFrom(actualAsTime.Truncate(matcher.TruncateBy))
	compareWith := null.TimeFrom(matcher.Expected.Time.Truncate(matcher.TruncateBy))

	if matcher.Exact {
		return compareWith.ExactEqual(actualAsNullTime), nil
	}
	return compareWith.Equal(actualAsNullTime), nil
}

func (matcher *TimeMatcher) FailureMessage(actual interface{}) (message string) {
	json, err := matcher.Expected.MarshalJSON()
	if err != nil {
		panic(err)
	}

	return format.Message(actual, "to equal", string(json))
}

func (matcher *TimeMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	json, err := matcher.Expected.MarshalJSON()
	if err != nil {
		panic(err)
	}

	return format.Message(actual, "not to equal", string(json))
}
