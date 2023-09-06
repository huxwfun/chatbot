package storage

import (
	"context"
	"huxwfun/chatbot/internal/models"
	"sort"
)

type WorkflowStorage struct {
	InMemStorage[models.Workflow]
}

func NewWorkflowStorage() *WorkflowStorage {
	inMem := NewInMemStorage[models.Workflow]()
	return &WorkflowStorage{
		inMem,
	}
}

// start WorkflowExecutionStorage
type WorkflowExecutionStorage struct {
	InMemStorage[models.WorkflowExecution]
}

func NewWorkflowExecutionStorage() *WorkflowExecutionStorage {
	inMem := NewInMemStorage[models.WorkflowExecution]()
	return &WorkflowExecutionStorage{
		inMem,
	}
}

func (w *WorkflowExecutionStorage) FindByCustomer(ctx context.Context, customerId string) []models.WorkflowExecution {
	result := make([]models.WorkflowExecution, 0, len(w.storage))
	for _, v := range w.storage {
		if v.CustomerId == customerId {
			result = append(result, v)
		}
	}
	return result
}

func (w *WorkflowExecutionStorage) FindActiveByCustomer(ctx context.Context, customerId string) []models.WorkflowExecution {
	result := make([]models.WorkflowExecution, 0, len(w.storage))
	for _, v := range w.storage {
		if v.CustomerId == customerId && v.CurrentState != models.WorkflowStateDone {
			result = append(result, v)
		}
	}
	return result
}
func (w *WorkflowExecutionStorage) FindByWorkflowAndCustomer(ctx context.Context, workflowId, customerId string) []models.WorkflowExecution {
	result := make([]models.WorkflowExecution, 0, len(w.storage))
	for _, v := range w.storage {
		if v.WorkflowId == workflowId && v.CustomerId == customerId {
			result = append(result, v)
		}
	}
	return result
}

// end WorkflowExecutionStorage

// start WorkflowExecutionLogStorage
type WorkflowLogStorage struct {
	InMemStorage[models.WorkflowLog]
}

func NewWorkflowLogStorage() *WorkflowLogStorage {
	inMem := NewInMemStorage[models.WorkflowLog]()
	return &WorkflowLogStorage{
		inMem,
	}
}

func (w *WorkflowLogStorage) FindByExecution(ctx context.Context, executionId string) []models.WorkflowLog {
	result := make([]models.WorkflowLog, 0, len(w.storage))
	for _, v := range w.storage {
		if v.ExecutionId == executionId {
			result = append(result, v)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		if !result[i].TimeCreated.Equal(result[j].TimeCreated) {
			return result[i].TimeCreated.Before(result[j].TimeCreated)
		}
		return result[i].Id < result[j].Id
	})
	return result
}

// end WorkflowExecutionLogStorage
