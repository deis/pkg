package aboutme

import (
	"net"
	"os"
	"testing"

	"k8s.io/kubernetes/pkg/client/unversioned"
)

func TestFromEnv(t *testing.T) {
	if _, err := unversioned.InClusterConfig(); err != nil {
		t.Skip("This can only be run inside Kubernetes. Skipping.")
	}

	me, err := FromEnv()
	if err != nil {
		t.Errorf("Could not get an environment: %s", err)
	}
	if len(me.Name) == 0 {
		t.Error("Could not get a pod name.")
	}
}

func TestShuntEnv(t *testing.T) {
	e := &Me{
		Annotations: map[string]string{"a": "a"},
		Labels:      map[string]string{"b": "b"},
		Name:        "c",
	}

	e.ShuntEnv()

	if "a" != os.Getenv("MY_ANNOTATION_A") {
		t.Errorf("Expected 'a', got '%s'", os.Getenv("MY_ANNOTATION_A"))
	}
	if "b" != os.Getenv("MY_LABEL_B") {
		t.Errorf("Expected 'b', got '%s'", os.Getenv("MY_LABEL_B"))
	}

	if "c" != os.Getenv("MY_NAME") {
		t.Errorf("Expected 'c', got '%s'", os.Getenv("MY_NAME"))
	}
}

func TestMyIP(t *testing.T) {
	if _, err := net.InterfaceByName("eth0"); err != nil {
		t.Skip("Host operating system does not have an eth0 device to test.")
	}

	ip, err := MyIP()
	if err != nil {
		t.Errorf("Could not get IP address: %s", err)
	}

	if len(ip) == 0 {
		t.Errorf("Expected a valid IP address. Got nuthin.")
	}
}
