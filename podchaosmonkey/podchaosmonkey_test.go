package podchaosmonkey

import (
	"context"
	"fmt"
	"podchaosmonkey/util"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

var (
	testNamespace                  = "test-namespace"
	testAnnotationKey              = "test-annotation"
	testPeriodSeconds              = 10
	testGracePeriodSeconds         = 11
	matchingTestAnnotationValue    = "true"
	matchingAnnotationString       = fmt.Sprintf("%s=%s", testAnnotationKey, matchingTestAnnotationValue)
	matchingAnnotationSelector, _  = labels.Parse(matchingAnnotationString)
	nonMatchingTestAnnotationValue = "false"
)

func TestNewPodChaosMonkey(t *testing.T) {
	wantedClient := &kubernetes.Clientset{}
	wantedAnnotations := labels.Everything()

	podChaosMonkey := NewPodChaosMonkey(wantedClient, wantedAnnotations, testNamespace, testPeriodSeconds, testGracePeriodSeconds)
	if podChaosMonkey.Client != wantedClient {
		t.Errorf("error setting client in NewPodChaosMonkey")
	}
	if podChaosMonkey.Annotations.String() != wantedAnnotations.String() {
		t.Errorf("error setting client in NewPodChaosMonkey")
	}
	if podChaosMonkey.Namespace != testNamespace {
		t.Errorf("error setting namespace in NewPodChaosMonkey: got %s, wanted %s", podChaosMonkey.Namespace, testNamespace)
	}
	if podChaosMonkey.PeriodSeconds != testPeriodSeconds {
		t.Errorf("got %d, wanted %d", podChaosMonkey.PeriodSeconds, testPeriodSeconds)
	}
	if podChaosMonkey.GracePeriodSeconds != int64(testGracePeriodSeconds) {
		t.Errorf("got %d, wanted %d", podChaosMonkey.GracePeriodSeconds, int64(testGracePeriodSeconds))
	}
}

func TestFilterPodsByAnnotations(t *testing.T) {
	matchingAnnotations := make(map[string]string)
	matchingAnnotations[testAnnotationKey] = matchingTestAnnotationValue
	matchingPod := v1.Pod{}
	matchingPod.Annotations = matchingAnnotations

	nonMatchingAnnotations := make(map[string]string)
	nonMatchingAnnotations[testAnnotationKey] = nonMatchingTestAnnotationValue
	nonMatchingPod := v1.Pod{}
	nonMatchingPod.Annotations = nonMatchingAnnotations

	podList := []v1.Pod{matchingPod, nonMatchingPod}

	filteredPods := filterPodsByAnnotations(podList, matchingAnnotationSelector)

	if len(filteredPods) != 1 {
		t.Errorf("Error filtering pods, expected len = %d, got %d", 1, len(filteredPods))
	}

	if filteredPods[0].Annotations[testAnnotationKey] != matchingTestAnnotationValue {
		t.Errorf("Expected filtered pod annotation '%s', got %s", matchingTestAnnotationValue, filteredPods[0].Annotations[testAnnotationKey])
	}
}

func TestFilterAndDeletePods(t *testing.T) {
	fakeClient := fake.NewSimpleClientset()
	podChaosMonkey := NewPodChaosMonkey(fakeClient, matchingAnnotationSelector, testNamespace, testPeriodSeconds, testGracePeriodSeconds)
	matchingPod := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "matching-pod",
			Annotations: map[string]string{
				testAnnotationKey: matchingTestAnnotationValue,
			},
		},
	}

	nonMatchingPodName := "non-matching-pod"
	nonMatchingPod := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: nonMatchingPodName,
			Annotations: map[string]string{
				testAnnotationKey: nonMatchingTestAnnotationValue,
			},
		},
	}

	fakeClient.CoreV1().Pods(testNamespace).Create(context.TODO(), &matchingPod, metav1.CreateOptions{})
	fakeClient.CoreV1().Pods(testNamespace).Create(context.TODO(), &nonMatchingPod, metav1.CreateOptions{})

	podChaosMonkey.filterAndDeletePods(util.NewFakeRandomGenerator(0))

	pods := podChaosMonkey.getPods()

	if len(pods) != 1 {
		t.Errorf("Error filtering and deleting pods, expected len = %d, got %d", 1, len(pods))
	}

	if pods[0].Name != nonMatchingPodName {
		t.Errorf("Expected '%s', got %s", &nonMatchingPod, pods[0].Name)
	}

}
