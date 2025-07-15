package comp_attr

import (
	"context"
	"reflect"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	intstr "k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGetServiceInfo(t *testing.T) {
	// Scheme 初始化
	scheme := runtime.NewScheme()
	v1.AddToScheme(scheme)

	// 测试用例
	tests := []struct {
		name        string
		namespace   string
		serviceName string
		setup       func(cs *fake.Clientset)
		expected    *ServiceInfo
		expectErr   bool
	}{
		{
			name:        "正常服务获取",
			namespace:   "monitoring",
			serviceName: "kube-prometheus-prometheus",
			setup: func(cs *fake.Clientset) {
				_, _ = cs.CoreV1().Services("monitoring").Create(
					context.TODO(),
					&v1.Service{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "kube-prometheus-prometheus",
							Namespace: "monitoring",
						},
						Spec: v1.ServiceSpec{
							ClusterIP: "10.109.24.227",
							Ports: []v1.ServicePort{
								{
									Name:       "http",
									Port:       9090,
									TargetPort: intstr.FromInt(9090),
								},
							},
						},
					},
					metav1.CreateOptions{},
				)
			},
			expected: &ServiceInfo{
				Name:      "kube-prometheus-prometheus",
				Namespace: "monitoring",
				ClusterIP: "10.109.24.227",
				Ports: []v1.ServicePort{
					{
						Name:       "http",
						Port:       9090,
						TargetPort: intstr.FromInt(9090),
					},
				},
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 初始化 fake client
			cs := fake.NewSimpleClientset()
			tt.setup(cs)

			// 调用被测函数
			service, err := cs.CoreV1().Services(tt.namespace).Get(
				context.TODO(),
				tt.serviceName,
				metav1.GetOptions{},
			)

			info, testErr := GetServiceInfoForTest(service, err)

			// 错误校验
			if (testErr != nil) != tt.expectErr {
				t.Fatalf("预期错误状态: %v, 实际: %v", tt.expectErr, testErr)
			}

			// 结果校验
			if !reflect.DeepEqual(info, tt.expected) {
				t.Errorf("预期结果: %+v, 实际结果: %+v", tt.expected, info)
			}
		})
	}
}
