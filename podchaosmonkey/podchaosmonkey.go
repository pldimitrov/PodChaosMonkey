package podchaosmonkey

import (
	"context"
	"fmt"
	"podchaosmonkey/util"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// PodChaosMonkey randomly kills a pod at a specified interval (PeriodSeconds)
// matches pods by annotations and namespace
type PodChaosMonkey struct {
	Client             kubernetes.Interface
	Annotations        labels.Selector
	Namespace          string
	PeriodSeconds      int
	GracePeriodSeconds int64
}

// NewPodChaosMonkey constructor
func NewPodChaosMonkey(client kubernetes.Interface, annotations labels.Selector, namespace string, periodSeconds int, gracePeriodSeconds int) *PodChaosMonkey {
	return &PodChaosMonkey{
		Client:             client,
		Annotations:        annotations,
		Namespace:          namespace,
		PeriodSeconds:      periodSeconds,
		GracePeriodSeconds: int64(gracePeriodSeconds),
	}
}

// InitClientSet helper to initialize a k8s clientset
func InitClientSet() kubernetes.Interface {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return clientset
}

func filterPodsByAnnotations(pods []v1.Pod, annotations labels.Selector) []v1.Pod {
	results := []v1.Pod{}

	for _, pod := range pods {
		selector := labels.Set(pod.Annotations)
		if annotations.Matches(selector) {
			results = append(results, pod)
		}
	}

	return results
}

func (p *PodChaosMonkey) deletePod(pod v1.Pod) error {
	fmt.Printf("deleting pod: %s in namespace: %s\n", pod.Name, p.Namespace)
	err := p.Client.CoreV1().Pods(pod.Namespace).Delete(context.TODO(), pod.Name, metav1.DeleteOptions{
		GracePeriodSeconds: &p.GracePeriodSeconds,
	})

	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}

func (p *PodChaosMonkey) getPods() []v1.Pod {
	pods, err := p.Client.CoreV1().Pods(p.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		// a panic is appropriate here as being unable to list pods is probably due to an RBAC issue
		panic(err.Error())
	}

	return pods.Items
}

func (p *PodChaosMonkey) deletePods(pods []v1.Pod) {
	for _, pod := range pods {
		p.deletePod(pod)
	}
}

// filters pods by annotations and namespace, selects a random pod to kill
func (p *PodChaosMonkey) filterAndDeletePods(randomGenerator util.IRandomGenerator) {
	pods := filterPodsByAnnotations(p.getPods(), p.Annotations)
	if len(pods) > 0 {
		p.deletePod(pods[randomGenerator.Int()%len(pods)])
	}
}

// Run start PodChaosMonkey
func (p *PodChaosMonkey) Run(ctx context.Context) {
	for {
		p.filterAndDeletePods(util.NewRandomGenerator())
		time.Sleep(time.Duration(p.PeriodSeconds) * time.Second)
	}
}
