package http

import (
	"encoding/json"
	"go-task-queue-system/domain"
	"go-task-queue-system/usecase"
	"log"
	"net/http"
	"strings"
)

type Handler struct {
	submitTaskUC *usecase.SubmitTaskUseCase
	getTaskUC    *usecase.GetTaskUseCase
	listTasksUC  *usecase.ListTasksUseCase
	cancelTaskUC *usecase.CancelTaskUseCase
	getStatsUC   *usecase.GetStatsUseCase
	workerPool   WorkerPool
}

type WorkerPool interface {
	GetStatus() map[string]interface{}
}

func NewHandler(
	submitTaskUC *usecase.SubmitTaskUseCase,
	getTaskUC *usecase.GetTaskUseCase,
	listTasksUC *usecase.ListTasksUseCase,
	cancelTaskUC *usecase.CancelTaskUseCase,
	getStatsUC *usecase.GetStatsUseCase,
	workerPool WorkerPool,
) *Handler {
	return &Handler{
		submitTaskUC: submitTaskUC,
		getTaskUC:    getTaskUC,
		listTasksUC:  listTasksUC,
		cancelTaskUC: cancelTaskUC,
		getStatsUC:   getStatsUC,
		workerPool:   workerPool,
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:  "healthy",
		Version: "1.0.0",
	}
	respondJSON(w, http.StatusOK, response)
}

func (h *Handler) SubmitTask(w http.ResponseWriter, r *http.Request) {
	var req SubmitTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	taskType := domain.TaskType(req.Type)
	if !taskType.IsValid() {
		respondError(w, http.StatusBadRequest, "Invalid task type", "")
		return
	}

	priority := domain.TaskPriorityMedium
	if req.Priority != "" {
		priority = domain.TaskPriority(req.Priority)
		if !priority.IsValid() {
			respondError(w, http.StatusBadRequest, "Invalid priority", "")
			return
		}
	}

	task, err := h.submitTaskUC.Execute(taskType, priority, req.Payload)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to submit task", err.Error())
		return
	}

	log.Printf("✅ Task submitted: %s (type: %s)", task.ID, task.Type)
	respondJSON(w, http.StatusCreated, ToTaskResponse(task))
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	// Extract task ID from URL path
	taskID := strings.TrimPrefix(r.URL.Path, "/tasks/")

	if taskID == "" {
		respondError(w, http.StatusBadRequest, "Task ID is required", "")
		return
	}

	task, err := h.getTaskUC.Execute(taskID)
	if err != nil {
		if err == domain.ErrTaskNotFound {
			respondError(w, http.StatusNotFound, "Task not found", "")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to retrieve task", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, ToTaskResponse(task))
}

func (h *Handler) ListTasks(w http.ResponseWriter, r *http.Request) {
	statusParam := r.URL.Query().Get("status")

	var tasks []*domain.Task
	var err error

	if statusParam != "" {
		status := domain.TaskStatus(statusParam)
		if !status.IsValid() {
			respondError(w, http.StatusBadRequest, "Invalid status", "")
			return
		}
		tasks, err = h.listTasksUC.ExecuteByStatus(status)
	} else {
		tasks, err = h.listTasksUC.Execute()
	}

	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to retrieve tasks", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, ToTaskListResponse(tasks))
}

func (h *Handler) CancelTask(w http.ResponseWriter, r *http.Request) {
	taskID := strings.TrimPrefix(r.URL.Path, "/tasks/")
	taskID = strings.TrimSuffix(taskID, "/cancel")

	if taskID == "" {
		respondError(w, http.StatusBadRequest, "Task ID is required", "")
		return
	}

	err := h.cancelTaskUC.Execute(taskID)
	if err != nil {
		if err == domain.ErrTaskNotFound {
			respondError(w, http.StatusNotFound, "Task not found", "")
			return
		}
		respondError(w, http.StatusBadRequest, "Failed to cancel task", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, SuccessResponse{Message: "Task cancelled successfully"})
}

func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.getStatsUC.Execute()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to retrieve stats", err.Error())
		return
	}

	response := StatsResponse{
		TotalTasks:      stats.TotalTasks,
		PendingTasks:    stats.PendingTasks,
		ProcessingTasks: stats.ProcessingTasks,
		CompletedTasks:  stats.CompletedTasks,
		FailedTasks:     stats.FailedTasks,
		CancelledTasks:  stats.CancelledTasks,
		QueueSize:       stats.QueueSize,
	}

	respondJSON(w, http.StatusOK, response)
}

func (h *Handler) GetWorkerStatus(w http.ResponseWriter, r *http.Request) {
	status := h.workerPool.GetStatus()
	respondJSON(w, http.StatusOK, status)
}

func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, statusCode int, message string, details string) {
	response := ErrorResponse{
		Error:   message,
		Message: details,
	}
	respondJSON(w, statusCode, response)
}
