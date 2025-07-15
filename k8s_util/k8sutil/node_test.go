// node_test.go (修正后)
package k8sutil

import (
	"context"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
)

func TestGetNodeInfo(t *testing.T) {
	tests := []struct {
		name           string
		nodeName       string
		setup          func(clientset *fake.Clientset) error
		expectedErr    bool
		expectedPhase  v1.NodePhase
		expectedOS     string
		expectedKernel string
	}{
		{
			name:     "existing node",
			nodeName: "edgecloud",
			setup: func(cs *fake.Clientset) error {
				_, err := cs.CoreV1().Nodes().Create(context.TODO(), &v1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name:   "edgecloud",
						Labels: map[string]string{"node-role.kubernetes.io/": "master"},
					},
					Status: v1.NodeStatus{
						Phase:    v1.NodeRunning,
						NodeInfo: v1.NodeSystemInfo{OSImage: "openEuler 24.03", KernelVersion: "6.6.0"},
					},
				}, metav1.CreateOptions{})
				return err
			},
			expectedErr:    false,
			expectedPhase:  v1.NodeRunning,
			expectedOS:     "openEuler 24.03",
			expectedKernel: "6.6.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeClient := fake.NewSimpleClientset()
			if tt.setup != nil {
				if err := tt.setup(fakeClient); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			}

			// 注入 mock client
			oldNewForConfig := kubernetes.NewForConfig
			defer func() { kubernetes.NewForConfig = oldNewForConfig }()
			kubernetes.NewForConfig = func(cfg *rest.Config) (*kubernetes.Clientset, error) {
				return fakeClient, nil
			}

			err := GetNodeInfo(tt.nodeName)
			if (err != nil) != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}
