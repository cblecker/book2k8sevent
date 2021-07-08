// Package book2k8sevent allows you to convert a book into a set of kubernetes events.
package main

import (
	"bufio"
	"bytes"
	"flag"
	"path/filepath"
	"time"

	"github.com/golang/glog"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	r := eventRecorder(clientset)
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "book",
			Namespace: "default",
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
			r.Event(pod, "Normal", "Paragraph", line)
			time.Sleep(2 * time.Second)
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
