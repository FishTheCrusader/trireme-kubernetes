package kubernetes

import (
	"fmt"

	"github.com/golang/glog"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/watch"
)

// PolicyWatcher iterates over the networkPolicyEvents. Each event generates a call to the parameter function.
func (c *Client) PolicyWatcher(namespace string, resultChan chan<- watch.Event, stopChan <-chan bool) error {
	for {
		watcher, err := c.kubeClient.Extensions().NetworkPolicies(namespace).Watch(api.ListOptions{})
		if err != nil {
			return fmt.Errorf("Couldn't open the Policy watch channel: %s", err)
		}
		for {
			select {
			case <-stopChan:
				return nil
			case req, open := <-watcher.ResultChan():
				if !open {
					glog.V(2).Infof("NetworkPolicy Watcher channel closed.")
					break
				}
				glog.V(4).Infof("Adding NetworkPolicyEvent")
				resultChan <- req
			}
		}
	}
}

// LocalPodWatcher iterates over the podEvents. Each event generates a call to the parameter function.
func (c *Client) LocalPodWatcher(namespace string, resultChan chan<- watch.Event, stopChan <-chan bool) error {
	option := c.localNodeOption()
	for {
		watcher, err := c.kubeClient.Pods(namespace).Watch(option)
		if err != nil {
			return fmt.Errorf("Couldn't open the Pod watch channel: %s", err)
		}
		for {
			select {
			case <-stopChan:
				return nil
			case req, open := <-watcher.ResultChan():
				if !open {
					glog.V(2).Infof("Namespace Watcher channel closed.")
					break
				}
				glog.V(4).Infof("Adding PodEvent")
				resultChan <- req
			}
		}
	}
}

// NamespaceWatcher iterates over the namespaceEvents. Each event generates a call to the parameter function.
func (c *Client) NamespaceWatcher(resultChan chan<- watch.Event, stopChan <-chan bool) error {
	for {
		watcher, err := c.kubeClient.Namespaces().Watch(api.ListOptions{})
		if err != nil {
			return fmt.Errorf("Couldn't open the Namespace watch channel: %s", err)
		}
		for {
			select {
			case <-stopChan:
				return nil
			case req, open := <-watcher.ResultChan():
				if !open {
					glog.V(2).Infof("Namespace Watcher channel closed.")
					break
				}
				glog.V(4).Infof("Adding NamespaceEvent")
				resultChan <- req
			}
		}
	}
}
