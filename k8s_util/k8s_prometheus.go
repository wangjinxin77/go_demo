package comp_attr

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// ServiceInfo 定义返回的服务信息结构
type ServiceInfo struct {
	Name      string           `json:"name"`
	Namespace string           `json:"namespace"`
	ClusterIP string           `json:"clusterIP"`
	Ports     []v1.ServicePort `json:"ports"`
}

// GetServiceInfo 获取指定服务信息
func GetServiceInfo(serviceName, namespace string) (*ServiceInfo, error) {
	config, err := getConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10)
	defer cancel()

	service, err := clientset.CoreV1().Services(namespace).Get(ctx, serviceName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %v", err)
	}

	return &ServiceInfo{
		Name:      service.Name,
		Namespace: service.Namespace,
		ClusterIP: service.Spec.ClusterIP,
		Ports:     service.Spec.Ports,
	}, nil
}

// getConfig 根据环境自动选择配置加载方式
func getConfig() (*rest.Config, error) {
	// 优先尝试集群内配置
	config, err := rest.InClusterConfig()
	if err == nil {
		return config, nil
	}

	// 本地环境尝试kubeconfig文件
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	if _, err := os.Stat(kubeconfig); os.IsNotExist(err) {
		kubeconfig = "" // 允许空值，测试时处理
	}

	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	return nil, fmt.Errorf("no valid k8s config found")
}

// 修正后（返回值和错误）
func GetServiceInfoForTest(service *v1.Service, err error) (*ServiceInfo, error) {
	if err != nil {
		return nil, err
	}
	return &ServiceInfo{
		Name:      service.Name,
		Namespace: service.Namespace,
		ClusterIP: service.Spec.ClusterIP,
		Ports:     service.Spec.Ports,
	}, nil
}
