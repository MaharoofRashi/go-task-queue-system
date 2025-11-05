package processor

import (
	"context"
	"fmt"
	"go-task-queue-system/domain"
	"log"
	"math/rand"
	"time"
)

type EmailProcessor struct{}

func NewEmailProcessor() *EmailProcessor {
	return &EmailProcessor{}
}

func (p *EmailProcessor) Process(ctx context.Context, task *domain.Task) (map[string]interface{}, error) {
	payload := task.Payload

	to, _ := payload["to"].(string)
	subject, _ := payload["subject"].(string)
	body, _ := payload["body"].(string)

	log.Printf("📧 [Email Processor] Processing task %s", task.ID)
	log.Printf("   To: %s", to)
	log.Printf("   Subject: %s", subject)
	log.Printf("   Body: %s", body)

	// Simulate email sending (2-3 seconds)
	processingTime := 2 + rand.Intn(2) // 2-3 seconds

	select {
	case <-time.After(time.Duration(processingTime) * time.Second):
		// Processing completed
	case <-ctx.Done():
		return nil, fmt.Errorf("task cancelled: %v", ctx.Err())
	}

	// Simulate 10% failure rate for testing retry mechanism
	if rand.Float32() < 0.10 {
		log.Printf("❌ [Email Processor] Failed to send email to %s", to)
		return nil, fmt.Errorf("failed to send email: SMTP connection timeout")
	}

	log.Printf("✅ [Email Processor] Email sent successfully to %s", to)

	result := map[string]interface{}{
		"message_id": fmt.Sprintf("msg-%d", time.Now().Unix()),
		"sent_at":    time.Now().Format(time.RFC3339),
		"recipient":  to,
	}

	return result, nil
}

func (p *EmailProcessor) CanProcess(taskType domain.TaskType) bool {
	return taskType == domain.TaskTypeEmail
}
