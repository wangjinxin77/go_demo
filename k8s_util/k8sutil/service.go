// service.go
package k8sutil

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// GetServiceInfo 获取指定namespace下的Service信息
func GetServiceInfo(serviceName, namespace string) error {
	config, err := getConfig()
	if err != nil {
		return fmt.Errorf("failed to get config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create clientset: %v", err)
	}

	service, err := clientset.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get service: %v", err)
	}

	fmt.Printf("Service: %s\n", service.Name)
	fmt.Printf("Type: %s\n", service.Spec.Type)
	fmt.Printf("ClusterIP: %s\n", service.Spec.ClusterIP)
	fmt.Printf("Ports: %v\n", service.Spec.Ports)
	fmt.Printf("Age: %s\n", service.CreationTimestamp.Time.Format("2006-01-02 15:04:05"))

	return nil
}

// getConfig 根据环境获取Kubernetes配置
func getConfig() (*rest.Config, error) {
	if _, err := os.Stat("/var/run/secrets/kubernetes.io/serviceaccount/token"); err == nil {
		// Pod环境使用InClusterConfig
		return rest.InClusterConfig()
	}
	// 本地环境使用kubeconfig
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	}
	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}
