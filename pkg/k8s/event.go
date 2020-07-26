package k8s

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/record"
)

//EmitEvent emits the event
func EmitEvent(kubeclientset kubernetes.Interface, pod runtime.Object, eventType, reason, message string) {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events(v1.NamespaceAll)})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "kubectl-emit-event"})
	recorder.Event(pod, eventType, reason, message)
	// as the event is called asynchronously,
	// the process exit without sending the event
	// if the sleep is not provided
	time.Sleep(200 * time.Millisecond)
}
