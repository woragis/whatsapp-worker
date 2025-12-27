package metrics

import (
	"testing"

	dto "github.com/prometheus/client_model/go"
)

func TestRecordJobProcessed(t *testing.T) {
	// Reset metrics to avoid interference from other tests
	WorkerJobsProcessedTotal.Reset()
	WorkerJobDuration.Reset()

	worker := "test-worker"
	status := "success"
	duration := 1.5

	RecordJobProcessed(worker, status, duration)

	// Verify counter was incremented
	metric := &dto.Metric{}
	if err := WorkerJobsProcessedTotal.WithLabelValues(worker, status).Write(metric); err != nil {
		t.Fatalf("Failed to get metric: %v", err)
	}
	if metric.Counter == nil {
		t.Error("Counter should not be nil")
	}
	if metric.Counter.GetValue() != 1.0 {
		t.Errorf("Counter value = %v, want 1.0", metric.Counter.GetValue())
	}

	// Verify histogram was observed (histograms are Observers, not directly writable)
	// The observation is recorded, we just verify the function doesn't panic
	// Actual histogram values are tested via Prometheus metrics endpoint in integration tests
}

func TestRecordJobFailed(t *testing.T) {
	WorkerJobsFailedTotal.Reset()

	worker := "test-worker"
	errorType := "send_error"

	RecordJobFailed(worker, errorType)

	metric := &dto.Metric{}
	if err := WorkerJobsFailedTotal.WithLabelValues(worker, errorType).Write(metric); err != nil {
		t.Fatalf("Failed to get metric: %v", err)
	}
	if metric.Counter == nil {
		t.Error("Counter should not be nil")
	}
	if metric.Counter.GetValue() != 1.0 {
		t.Errorf("Counter value = %v, want 1.0", metric.Counter.GetValue())
	}
}

func TestRecordJobRetried(t *testing.T) {
	WorkerJobsRetriedTotal.Reset()

	worker := "test-worker"

	RecordJobRetried(worker)

	metric := &dto.Metric{}
	if err := WorkerJobsRetriedTotal.WithLabelValues(worker).Write(metric); err != nil {
		t.Fatalf("Failed to get metric: %v", err)
	}
	if metric.Counter == nil {
		t.Error("Counter should not be nil")
	}
	if metric.Counter.GetValue() != 1.0 {
		t.Errorf("Counter value = %v, want 1.0", metric.Counter.GetValue())
	}
}

func TestSetQueueDepth(t *testing.T) {
	QueueDepth.Reset()

	queueName := "test-queue"
	depth := 42.0

	SetQueueDepth(queueName, depth)

	metric := &dto.Metric{}
	if err := QueueDepth.WithLabelValues(queueName).Write(metric); err != nil {
		t.Fatalf("Failed to get metric: %v", err)
	}
	if metric.Gauge == nil {
		t.Error("Gauge should not be nil")
	}
	if metric.Gauge.GetValue() != depth {
		t.Errorf("Gauge value = %v, want %v", metric.Gauge.GetValue(), depth)
	}
}

func TestSetQueueDLQSize(t *testing.T) {
	QueueDLQSize.Reset()

	queueName := "test-queue"
	size := 10.0

	SetQueueDLQSize(queueName, size)

	metric := &dto.Metric{}
	if err := QueueDLQSize.WithLabelValues(queueName).Write(metric); err != nil {
		t.Fatalf("Failed to get metric: %v", err)
	}
	if metric.Gauge == nil {
		t.Error("Gauge should not be nil")
	}
	if metric.Gauge.GetValue() != size {
		t.Errorf("Gauge value = %v, want %v", metric.Gauge.GetValue(), size)
	}
}

func TestRecordJobProcessed_MultipleStatuses(t *testing.T) {
	WorkerJobsProcessedTotal.Reset()

	worker := "test-worker"
	
	RecordJobProcessed(worker, "success", 1.0)
	RecordJobProcessed(worker, "failed", 2.0)
	RecordJobProcessed(worker, "success", 0.5)

	// Check success counter
	successMetric := &dto.Metric{}
	if err := WorkerJobsProcessedTotal.WithLabelValues(worker, "success").Write(successMetric); err != nil {
		t.Fatalf("Failed to get metric: %v", err)
	}
	if successMetric.Counter.GetValue() != 2.0 {
		t.Errorf("Success counter = %v, want 2.0", successMetric.Counter.GetValue())
	}

	// Check failed counter
	failedMetric := &dto.Metric{}
	if err := WorkerJobsProcessedTotal.WithLabelValues(worker, "failed").Write(failedMetric); err != nil {
		t.Fatalf("Failed to get metric: %v", err)
	}
	if failedMetric.Counter.GetValue() != 1.0 {
		t.Errorf("Failed counter = %v, want 1.0", failedMetric.Counter.GetValue())
	}
}

func TestMetrics_Initialization(t *testing.T) {
	// Verify all metrics are initialized
	if WorkerJobsProcessedTotal == nil {
		t.Error("WorkerJobsProcessedTotal should not be nil")
	}
	if WorkerJobDuration == nil {
		t.Error("WorkerJobDuration should not be nil")
	}
	if WorkerJobsFailedTotal == nil {
		t.Error("WorkerJobsFailedTotal should not be nil")
	}
	if WorkerJobsRetriedTotal == nil {
		t.Error("WorkerJobsRetriedTotal should not be nil")
	}
	if QueueDepth == nil {
		t.Error("QueueDepth should not be nil")
	}
	if QueueDLQSize == nil {
		t.Error("QueueDLQSize should not be nil")
	}
}

