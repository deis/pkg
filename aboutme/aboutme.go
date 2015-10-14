// Package aboutme provides information to a pod about itself.
//
// Typical usage is to let the Pod auto-detect information about itself:
//
//	my, err := aboutme.FromEnv()
//  if err != nil {
// 		// Error connecting to tke k8s API server
// 	}
//
// 	fmt.Printf("My Pod Name is %s", my.Name)
package aboutme

import (
	"errors"
	"net"
	"os"
	"strings"

	"github.com/deis/deis/pkg/k8s"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/labels"
)

type Me struct {
	ApiServer, Name                      string
	IP, NodeIP, Namespace, SelfLink, UID string
	Labels                               map[string]string
	Annotations                          map[string]string

	c *unversioned.Client
}

// FromEnv uses the environment to create a new Me.
//
// To use this, a client MUST be running inside of a Pod environment. It uses
// a combination of environment variables and file paths to determine
// information about the cluster.
func FromEnv() (*Me, error) {

	host := os.Getenv("KUBERNETES_SERVICE_HOST")
	port := os.Getenv("KUBERNETES_SERVICE_PORT")
	name := os.Getenv("HOSTNAME")

	// FIXME: Better way? Probably scanning secrets for
	// an SSL cert would help?
	proto := "https"

	url := proto + "://" + host + ":" + port

	me := &Me{
		ApiServer: url,
		Name:      name,

		// FIXME: This is a chicken-and-egg problem. We need the namespace
		// to get pod info, and we can only get info from the pod.
		Namespace: "default",
	}

	client, err := k8s.PodClient()
	if err != nil {
		return me, err
	}
	me.c = client

	me.init()

	return me, nil
}

// Client returns an initialized Kubernetes API client.
func (me *Me) Client() *unversioned.Client {
	return me.c
}

// ShuntEnv puts the Me object into the environment.
//
// The properties of Me are placed into the environment according to the
// following rules:
//
// 	- In general, all variables are prefaced with MY_ (MY_IP, MY_NAMESPACE)
// 	- Labels become MY_LABEL_[NAME]=[value]
// 	- Annotations become MY_ANNOTATION_[NAME] = [value]
func (me *Me) ShuntEnv() {
	env := map[string]string{
		"MY_APISERVER": me.ApiServer,
		"MY_NAME":      me.Name,
		"MY_IP":        me.IP,
		"MY_NODEIP":    me.NodeIP,
		"MY_NAMESPACE": me.Namespace,
		"MY_SELFLINK":  me.SelfLink,
		"MY_UID":       me.UID,
	}
	for k, v := range env {
		os.Setenv(k, v)
	}
	var name string
	for k, v := range me.Labels {
		name = "MY_LABEL_" + strings.ToUpper(k)
		os.Setenv(name, v)
	}
	for k, v := range me.Annotations {
		name = "MY_ANNOTATION_" + strings.ToUpper(k)
		os.Setenv(name, v)
	}
}

func (me *Me) init() error {
	p, n, err := me.findPodInNamespaces()
	if err != nil {
		return err
	}

	me.Namespace = n
	me.IP = p.Status.PodIP
	me.NodeIP = p.Status.HostIP
	me.SelfLink = p.SelfLink
	me.UID = string(p.UID)
	me.Labels = p.Labels
	me.Annotations = me.Annotations

	return nil
}

// findPodInNamespaces searches relevant namespaces for this pod.
//
// It returns a PodInterface for working with the pod, a namespace name as a
// string, and an error if something goes wrong.
//
// The search pattern is to look for namespaces that have the "deis" name in
// their labels, and then to fall back to default. We don't look at all
// namespaces.
func (me *Me) findPodInNamespaces() (*api.Pod, string, error) {
	// Get the deis namespace. If it does not exist, get the default namespce.
	s, err := labels.Parse("name=deis")
	if err == nil {
		ns, err := me.c.Namespaces().List(s, nil)
		if err != nil {
			return nil, "default", err
		}
		for _, n := range ns.Items {
			p, err := me.c.Pods(n.Name).Get(me.Name)

			// If there is no error, we got a matching pod.
			if err == nil {
				return p, n.Name, nil
			}
		}
	}

	// If we get here, it's really the last ditch.
	p, err := me.c.Pods("default").Get(me.Name)
	return p, "default", err
}

// MyIP examines the local interfaces and guesses which is its IP.
//
// Containers tend to put the IP address in eth0, so this attempts to look up
// that interface and retrieve its IP. It is fairly naive. To get more
// thorough IP information, you may prefer to use the `net` package and
// look up the desired information.
//
// Because this queries the interfaces, not the Kube API server, this could,
// in theory, return an IP address different from Me.IP.
func MyIP() (string, error) {
	iface, err := net.InterfaceByName("eth0")
	if err != nil {
		return "", err
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return "", err
	}
	var ip string
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
			}
		}
	}
	if len(ip) == 0 {
		return ip, errors.New("Found no IPv4 addresses.")
	}
	return ip, nil
}
