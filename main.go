// Package book2k8sevent allows you to convert a book into a set of kubernetes events.
package main

import (
	"bufio"
	"bytes"
	"time"

	"github.com/golang/glog"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
)

func main() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	r := eventRecorder(clientset)
	d := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "rabbit-hole",
			Namespace: "wonderland",
		},
	}

	// Open the file.
	data := MustAsset("data/alice_in_wonderland.txt")
	// Create a new scanner for the file.
	scanner := bufio.NewScanner(bytes.NewReader(data))
	// Loop over all lines in the file and print them.
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			r.Event(d, "Normal", "Paragraph", line)
			time.Sleep(7 * time.Second)
		}
	}
}

func eventRecorder(
	kubeClient *kubernetes.Clientset) record.EventRecorder {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(
		&typedcorev1.EventSinkImpl{
			Interface: kubeClient.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(
		scheme.Scheme,
		v1.EventSource{Component: "book2k8sevent"})
	return recorder
}
