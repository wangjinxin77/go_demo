package periodic

import (
	"fmt"
	"reflect"
	"sync/atomic"
	"testing"
	"time"
)

// PeriodicScheduler 定义了一个周期性任务调度器
type PeriodicScheduler struct {
	interval   time.Duration               // 任务执行的时间间隔
	task       func() (interface{}, error) // 任务函数，返回结果或错误
	resultChan chan<- interface{}          // 用于传递任务结果的通道
	stopChan   chan struct{}               // 用于停止调度器的通道
}

// NewPeriodicScheduler 创建一个新的 PeriodicScheduler 实例
func NewPeriodicScheduler(
	interval time.Duration,
	task func() (interface{}, error),
	resultChan chan<- interface{},
) *PeriodicScheduler {
	return &PeriodicScheduler{
		interval:   interval,
		task:       task,
		resultChan: resultChan,
		stopChan:   make(chan struct{}),
	}
}

// Start 启动调度器，定期执行任务并将结果发送到 resultChan
func (s *PeriodicScheduler) Start() {
	ticker := time.NewTicker(s.interval) // 定时器
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C: // 定时触发任务
			result, err := s.task()
			if err != nil {
				// 将错误通过通道传递
				s.resultChan <- fmt.Errorf("task error: %w", err)
			} else {
				// 将正常结果通过通道传递
				s.resultChan <- result
			}
		case <-s.stopChan: // 停止信号
			return
		}
	}
}

// Stop 停止调度器
func (s *PeriodicScheduler) Stop() {
	close(s.stopChan)
}

// TestPeriodicScheduler 测试 PeriodicScheduler 的功能
func TestPeriodicScheduler(t *testing.T) {
	// 创建带缓冲的通道以防止阻塞
	taskResult := make(chan interface{}, 10)

	// 使用原子计数器模拟任务
	var counter int64
	scheduler := NewPeriodicScheduler(
		time.Millisecond*100, // 每100ms执行一次任务
		func() (interface{}, error) {
			// 使用原子操作保证并发安全
			return atomic.AddInt64(&counter, 1), nil
		},
		taskResult,
	)

	go scheduler.Start()
	defer scheduler.Stop()

	// 等待足够的时间以确保任务执行
	time.Sleep(time.Millisecond * 350)

	// 期望的结果
	expected := []int64{1, 2, 3}
	actual := []int64{}

	// 从通道中读取任务结果
	for len(actual) < len(expected) {
		select {
		case res := <-taskResult:
			// 检查结果类型并追加到实际结果中
			if val, ok := res.(int64); ok {
				actual = append(actual, val)
			} else {
				t.Errorf("Unexpected result type: %T", res)
			}
		case <-time.After(time.Millisecond * 50): // 超时处理
			t.Fatal("Timeout waiting for scheduled task result")
		}
	}

	// 验证实际结果是否与期望结果一致
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func periodicTest() {
	// 示例用法：每秒打印当前时间戳
	timestampChan := make(chan interface{}, 5) // 创建一个缓冲通道
	scheduler := NewPeriodicScheduler(
		time.Second, // 每秒执行一次任务
		func() (interface{}, error) {
			// 返回当前时间戳
			return time.Now().Format(time.RFC3339), nil
		},
		timestampChan,
	)

	go scheduler.Start() // 启动调度器
	defer scheduler.Stop()

	// 接收并打印5次时间戳
	for i := 0; i < 5; i++ {
		fmt.Println(<-timestampChan)
	}
}
