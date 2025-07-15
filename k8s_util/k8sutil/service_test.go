// service_test.go
package k8sutil

import (
	"context"
	"os"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	// "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestGetServiceInfo(t *testing.T) {
	tests := []struct {
		name        string
		namespace   string
		serviceName string
		setup       func(clientset *fake.Clientset)
		expectedErr bool
	}{
		{
			name:        "existing service",
			namespace:   "default",
			serviceName: "test-svc",
			setup: func(cs *fake.Clientset) {
				cs.CoreV1().Services("default").Create(context.TODO(), &v1.Service{
					ObjectMeta: metav1.ObjectMeta{Name: "test-svc"},
					Spec: v1.ServiceSpec{
						Type:      v1.ServiceTypeNodePort,
						ClusterIP: "10.96.0.1",
						Ports:     []v1.ServicePort{{Port: 80}},
					},
				}, metav1.CreateOptions{})
			},
			expectedErr: false,
		},
		{
			name:        "non-existing service",
			namespace:   "default",
			serviceName: "non-existent",
			setup:       nil,
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeClient := fake.NewSimpleClientset()
			if tt.setup != nil {
				tt.setup(fakeClient)
			}

			oldNewForConfig := kubernetes.NewForConfig
			defer func() { kubernetes.NewForConfig = oldNewForConfig }()
			kubernetes.NewForConfig = func(cfg *rest.Config) (*kubernetes.Clientset, error) {
				return fakeClient, nil
			}

			os.Setenv("NAMESPACE", tt.namespace)
			defer os.Unsetenv("NAMESPACE")

			err := GetServiceInfo(tt.serviceName, tt.namespace)
			if (err != nil) != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}
