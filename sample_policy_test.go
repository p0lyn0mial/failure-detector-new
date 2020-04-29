package failure_detector

import (
	"fmt"
	"testing"
)

func TestSimpleWeightedEndpointStatusEvaluator(t *testing.T) {
	scenarios := []struct {
		name               string
		endpoint           *WeightedEndpointStatus
		expectedStatus     string
		expectedWeight     float32
		expectedHasChanged bool
	}{
		{
			name:               "an endpoint received 10 errors - status set to error, weight set to 0",
			endpoint:           createWeightedEndpointStatus(errToSampleFunc(genErrors(10)...)),
			expectedStatus:     EndpointStatusReasonTooManyErrors,
			expectedWeight:     0.0,
			expectedHasChanged: true,
		},
		{
			name:               "an endpoint received 4 errors - status NOT set, weight set to 0.6",
			endpoint:           createWeightedEndpointStatus(errToSampleFunc(genErrors(4)...)),
			expectedWeight:     0.6,
			expectedHasChanged: true,
		},
		{
			name:               "an endpoint received 0 errors - status NOT set, weight set to 1",
			endpoint:           createWeightedEndpointStatus(errToSampleFunc(genNilErrors(1)...)),
			expectedWeight:     1,
			expectedHasChanged: false,
		},
		{
			name: "an endpoint received a non-error sample after receiving 10 errors - status reset, weight set to 0.1",
			endpoint: func() *WeightedEndpointStatus {
				errors := genErrors(10)
				errors = append(errors, genNilErrors(1)...)
				ep := createWeightedEndpointStatus(errToSampleFunc(errors...))
				return ep
			}(),
			expectedWeight:     0.1,
			expectedHasChanged: true,
		},
		{
			name: "an endpoint received a non-error sample after receiving 3 non-errors samples - status reset, weight set to 0.1",
			endpoint: func() *WeightedEndpointStatus {
				errors := genNilErrors(3)
				errors = append(errors, genNilErrors(1)...)
				ep := createWeightedEndpointStatus(errToSampleFunc(errors...))
				ep.weight = 1
				return ep
			}(),
			expectedWeight:     1,
			expectedHasChanged: false,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			// test data
			hasChanged := SimpleWeightedEndpointStatusEvaluator(scenario.endpoint)
			if hasChanged != scenario.expectedHasChanged {
				t.Errorf("expected the method to return %v value bot got %v value", scenario.expectedHasChanged, hasChanged)
			}
			if scenario.endpoint.status != scenario.expectedStatus {
				t.Errorf("expected to get %s status but got %s", scenario.expectedStatus, scenario.endpoint.status)
			}
			if weightToErrorCount(scenario.endpoint.weight) != weightToErrorCount(scenario.expectedWeight) {
				t.Errorf("expected to get %d errors based on %f weight but got %d errors based on %f", weightToErrorCount(scenario.expectedWeight), scenario.expectedWeight, weightToErrorCount(scenario.endpoint.weight), scenario.endpoint.weight)
			}
		})
	}
}

func createWeightedEndpointStatus(samples []*Sample) *WeightedEndpointStatus {
	target := newWeightedEndpoint(10, nil)
	for _, sample := range samples {
		target.Add(sample)
	}
	return target
}

func errToSampleFunc(err ...error) []*Sample {
	ret := []*Sample{}
	for _, e := range err {
		ret = append(ret, &Sample{err: e})
	}
	return ret
}

func genErrors(number int) []error {
	ret := make([]error, number)
	for i := 0; i < number; i++ {
		ret[i] = fmt.Errorf("error %d", i)
	}
	return ret
}

func genNilErrors(number int) []error {
	ret := make([]error, number)
	for i := 0; i < number; i++ {
		ret[i] = nil
	}
	return ret
}
