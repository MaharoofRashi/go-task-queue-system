package main

import (
	"go-task-queue-system/domain"
	"go-task-queue-system/infrastructure/processor"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpDelivery "go-task-queue-system/delivery/http"
	"go-task-queue-system/infrastructure/queue"
	"go-task-queue-system/infrastructure/repository"
	"go-task-queue-system/infrastructure/worker"
	"go-task-queue-system/usecase"
)

const (
	serverPort    = ":8080"
	queueCapacity = 100
	workerCount   = 5
	workerTimeout = 30 * time.Second
)

func main() {
	log.Println("🚀 Starting Task Queue System...")

	// 1. Initialize Infrastructure Layer

	// Repository (in-memory storage)
	taskRepository := repository.NewMemoryRepository()
	log.Println("✅ Repository initialized")

	// Queue (channel-based)
	taskQueue := queue.NewChannelQueue(queueCapacity)
	log.Printf("✅ Queue initialized (capacity: %d)", queueCapacity)

	// Processor Registry
	processorRegistry := processor.NewProcessorRegistry()

	// Register task processors
	processorRegistry.Register(domain.TaskTypeEmail, processor.NewEmailProcessor())
	processorRegistry.Register(domain.TaskTypeImageProcessing, processor.NewImageProcessor())
	processorRegistry.Register(domain.TaskTypeReportGeneration, processor.NewReportProcessor())
	log.Println("✅ Task processors registered")

	// Worker Pool
	workerPool := worker.NewWorkerPool(
		workerCount,
		taskQueue.GetChannel(),
		taskRepository,
		processorRegistry,
		workerTimeout,
	)
	workerPool.Start()
	log.Printf("✅ Worker pool started (%d workers)", workerCount)

	// 2. Initialize Use Cases Layer

	submitTaskUC := usecase.NewSubmitTaskUseCase(taskRepository, taskQueue)
	getTaskUC := usecase.NewGetTaskUseCase(taskRepository)
	listTasksUC := usecase.NewListTasksUseCase(taskRepository)
	cancelTaskUC := usecase.NewCancelTaskUseCase(taskRepository)
	getStatsUC := usecase.NewGetStatsUseCase(taskRepository, taskQueue)
	log.Println("✅ Use cases initialized")

	// 3. Initialize HTTP Delivery Layer

	handler := httpDelivery.NewHandler(
		submitTaskUC,
		getTaskUC,
		listTasksUC,
		cancelTaskUC,
		getStatsUC,
		workerPool,
	)

	router := httpDelivery.SetupRoutes(handler)
	log.Println("✅ HTTP routes configured")

	// 4. Start HTTP Server

	server := &http.Server{
		Addr:         serverPort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🌐 Server starting on http://localhost%s", serverPort)
		log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		log.Println("📋 Available Endpoints:")
		log.Println("   GET  /health              - Health check")
		log.Println("   POST /tasks               - Submit a task")
		log.Println("   GET  /tasks               - List all tasks")
		log.Println("   GET  /tasks?status=pending - Filter by status")
		log.Println("   GET  /tasks/{id}          - Get task by ID")
		log.Println("   POST /tasks/{id}/cancel   - Cancel a task")
		log.Println("   GET  /stats               - System statistics")
		log.Println("   GET  /workers/status      - Worker pool status")
		log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		log.Println("✨ Ready to accept requests!")
		log.Println("")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server failed to start: %v", err)
		}
	}()

	// 5. Graceful Shutdown

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("")
	log.Println("🛑 Shutting down gracefully...")

	// Stop accepting new tasks
	taskQueue.Close()
	log.Println("✅ Queue closed")

	// Stop worker pool
	workerPool.Stop()
	log.Println("✅ Workers stopped")

	// Shutdown HTTP server
	log.Println("✅ HTTP server stopped")

	log.Println("👋 Goodbye!")
}
