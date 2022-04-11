package util

import (
	"strconv"
	"testing"
)

func TestParseArgs(t *testing.T) {
	wantedAnnotation := "test-annotation"
	wantedIntervalSeconds := 10
	wantedNamespace := "test-namespace"
	wantedGracePeriodSeconds := 11

	t.Setenv("PODCHAOSMONKEY_ANNOTATIONS", wantedAnnotation)
	t.Setenv("PODCHAOSMONKEY_INTERVAL_SECONDS", strconv.Itoa(wantedIntervalSeconds))
	t.Setenv("PODCHAOSMONKEY_NAMESPACE", wantedNamespace)
	t.Setenv("PODCHAOSMONKEY_GRACE_PERIOD_SECONDS", strconv.Itoa(wantedGracePeriodSeconds))

	annotations, namespace, intervalSeconds, gracePeriodSeconds, err := ParseArgs()

	if err != nil {
		t.Error(err.Error())
	}

	if annotations.String() != wantedAnnotation {
		t.Errorf("got %s, wanted %s", annotations.String(), wantedAnnotation)
	}
	if intervalSeconds != wantedIntervalSeconds {
		t.Errorf("got %d, wanted %d", intervalSeconds, wantedIntervalSeconds)
	}
	if namespace != wantedNamespace {
		t.Errorf("got %s, wanted %s", namespace, wantedNamespace)
	}
	if gracePeriodSeconds != wantedGracePeriodSeconds {
		t.Errorf("got %d, wanted %d", gracePeriodSeconds, wantedGracePeriodSeconds)
	}
}
