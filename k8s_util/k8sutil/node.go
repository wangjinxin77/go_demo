// node.go
package k8sutil

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GetNodeInfo 获取节点信息
func GetNodeInfo(nodeName string) error {
	config, err := getConfig()
	if err != nil {
		return fmt.Errorf("failed to get config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create clientset: %v", err)
	}

	node, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get node: %v", err)
	}

	fmt.Printf("Node: %s\n", node.Name)
	fmt.Printf("Status: %s\n", node.Status.Phase)
	fmt.Printf("Roles: %s\n", node.ObjectMeta.Labels["node-role.kubernetes.io/"])
	fmt.Printf("OS: %s\n", node.Status.NodeInfo.OSImage)
	fmt.Printf("Kernel: %s\n", node.Status.NodeInfo.KernelVersion)
	fmt.Printf("Age: %s\n", node.CreationTimestamp.Time.Format("2006-01-02 15:04:05"))

	return nil
}
