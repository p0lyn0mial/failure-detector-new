package failure_detector

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestEndpointType(t *testing.T) {
	errToSampleFunc := func(err ...error) []*Sample {
		ret := []*Sample{}
		for _, e := range err {
			ret = append(ret, &Sample{err: e})
		}
		return ret
	}
	scenarios := []struct {
		name       string
		size       int
		itemsToAdd int
		output     []*Sample
	}{

		{
			name:       "scenario 1: adds exactly \"size\" items to the endpoint type - no wrapping",
			size:       5,
			itemsToAdd: 5,
			output:     errToSampleFunc(errors.New("0"), errors.New("1"), errors.New("2"), errors.New("3"), errors.New("4")),
		},
		{
			name:       "scenario 2: adds more than \"size\" items to the endpoint type - wrapping",
			size:       5,
			itemsToAdd: 6,
			output:     errToSampleFunc(errors.New("1"), errors.New("2"), errors.New("3"), errors.New("4"), errors.New("5")),
		},
		{
			name:       "scenario 3: adds twice \"size\" items to the endpoint type - wrapping",
			size:       5,
			itemsToAdd: 10,
			output:     errToSampleFunc(errors.New("5"), errors.New("6"), errors.New("7"), errors.New("8"), errors.New("9")),
		},
		{
			name:       "scenario 4: adds twice + 3 \"size\" items to the endpoint type - wrapping",
			size:       5,
			itemsToAdd: 13,
			output:     errToSampleFunc(errors.New("8"), errors.New("9"), errors.New("10"), errors.New("11"), errors.New("12")),
		},
		{
			name:       "scenario 5: adds less than \"size\" items to the endpoint type - no wrapping",
			size:       5,
			itemsToAdd: 3,
			output:     errToSampleFunc(errors.New("0"), errors.New("1"), errors.New("2")),
		},
		{
			name:       "scenario 6: adds 0 items to the endpoint type - no wrapping",
			size:       5,
			itemsToAdd: 0,
			output:     errToSampleFunc(),
		},
		{
			name:       "scenario 7: adds 1984 items to the endpoint type - wrapping",
			size:       5,
			itemsToAdd: 1984,
			output:     errToSampleFunc(errors.New("1979"), errors.New("1980"), errors.New("1981"), errors.New("1982"), errors.New("1983")),
		},
		{
			name:       "scenario 8: adds 99 items to the endpoint type - wrapping - size 1",
			size:       1,
			itemsToAdd: 99,
			output:     errToSampleFunc(errors.New("98")),
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {

			target := newWeightedEndpoint(scenario.size, nil)
			for i := 0; i < scenario.itemsToAdd; i++ {
				target.Add(errToSampleFunc(fmt.Errorf("%d", i))[0])
			}

			actualOutput := target.Get()
			if !reflect.DeepEqual(actualOutput, scenario.output) {
				t.Fatalf("expected %v got %v", scenario.output, actualOutput)
			}
		})
	}
}
