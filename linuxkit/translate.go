package linuxkit

import "k8s.io/client-go/pkg/runtime"

//Template generates a linuxkit template from a k8 deployment object
func Template(obj runtime.Object) (string, error) {
	return "", nil
}
