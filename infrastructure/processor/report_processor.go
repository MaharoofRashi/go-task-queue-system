package processor

import (
	"context"
	"fmt"
	"go-task-queue-system/domain"
	"log"
	"math/rand"
	"time"
)

type ReportProcessor struct{}

func NewReportProcessor() *ReportProcessor {
	return &ReportProcessor{}
}

// Process simulates generating a report (PDF, CSV, etc.)
func (p *ReportProcessor) Process(ctx context.Context, task *domain.Task) (map[string]interface{}, error) {
	payload := task.Payload

	reportType, _ := payload["report_type"].(string)
	startDate, _ := payload["start_date"].(string)
	endDate, _ := payload["end_date"].(string)
	format, _ := payload["format"].(string)

	log.Printf("📊 [Report Processor] Processing task %s", task.ID)
	log.Printf("   Report Type: %s", reportType)
	log.Printf("   Period: %s to %s", startDate, endDate)
	log.Printf("   Format: %s", format)

	// Simulate data fetching (2 seconds)
	log.Printf("   📥 Fetching data...")
	select {
	case <-time.After(2 * time.Second):
	case <-ctx.Done():
		return nil, fmt.Errorf("task cancelled during data fetch: %v", ctx.Err())
	}

	// Simulate report generation (4-6 seconds)
	processingTime := 4 + rand.Intn(3) // 4-6 seconds
	log.Printf("   📝 Generating report...")

	select {
	case <-time.After(time.Duration(processingTime) * time.Second):
		// Processing completed
	case <-ctx.Done():
		return nil, fmt.Errorf("task cancelled during generation: %v", ctx.Err())
	}

	// Simulate 5% failure rate for testing retry mechanism
	if rand.Float32() < 0.05 {
		log.Printf("❌ [Report Processor] Failed to generate report")
		return nil, fmt.Errorf("report generation failed: insufficient data")
	}

	log.Printf("✅ [Report Processor] Report generated successfully")

	reportURL := fmt.Sprintf("https://reports.example.com/%s-%d.%s", reportType, time.Now().Unix(), format)
	result := map[string]interface{}{
		"report_url":    reportURL,
		"report_type":   reportType,
		"format":        format,
		"total_records": rand.Intn(10000) + 1000, // 1000-11000 records
		"file_size":     rand.Intn(5000) + 500,   // 500-5500 KB
		"generated_at":  time.Now().Format(time.RFC3339),
	}

	return result, nil
}

func (p *ReportProcessor) CanProcess(taskType domain.TaskType) bool {
	return taskType == domain.TaskTypeReportGeneration
}
