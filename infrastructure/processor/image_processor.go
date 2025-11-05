package processor

import (
	"context"
	"fmt"
	"go-task-queue-system/domain"
	"log"
	"math/rand"
	"time"
)

type ImageProcessor struct{}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{}
}

func (p *ImageProcessor) Process(ctx context.Context, task *domain.Task) (map[string]interface{}, error) {
	payload := task.Payload

	imageURL, _ := payload["image_url"].(string)
	width, _ := payload["width"].(float64)
	height, _ := payload["height"].(float64)
	format, _ := payload["format"].(string)

	log.Printf("🖼️  [Image Processor] Processing task %s", task.ID)
	log.Printf("   Image URL: %s", imageURL)
	log.Printf("   Target Size: %.0fx%.0f", width, height)
	log.Printf("   Format: %s", format)

	// Simulate downloading image (1 second)
	log.Printf("   ⬇️  Downloading image...")
	select {
	case <-time.After(1 * time.Second):
	case <-ctx.Done():
		return nil, fmt.Errorf("task cancelled during download: %v", ctx.Err())
	}

	// Simulate image processing (3-5 seconds)
	processingTime := 3 + rand.Intn(3) // 3-5 seconds
	log.Printf("   ⚙️  Processing image...")

	select {
	case <-time.After(time.Duration(processingTime) * time.Second):
		// Processing completed
	case <-ctx.Done():
		return nil, fmt.Errorf("task cancelled during processing: %v", ctx.Err())
	}

	// Simulate 15% failure rate for testing retry mechanism
	if rand.Float32() < 0.15 {
		log.Printf("❌ [Image Processor] Failed to process image")
		return nil, fmt.Errorf("image processing failed: invalid image format")
	}

	log.Printf("✅ [Image Processor] Image processed successfully")

	processedURL := fmt.Sprintf("https://cdn.example.com/processed/%d.%s", time.Now().Unix(), format)
	result := map[string]interface{}{
		"processed_url": processedURL,
		"width":         width,
		"height":        height,
		"format":        format,
		"file_size":     rand.Intn(500) + 100, // 100-600 KB
		"processed_at":  time.Now().Format(time.RFC3339),
	}

	return result, nil
}

func (p *ImageProcessor) CanProcess(taskType domain.TaskType) bool {
	return taskType == domain.TaskTypeImageProcessing
}
