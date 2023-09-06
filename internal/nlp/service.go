package nlp

import (
	"context"
	"reflect"
)

type BoolResult struct {
	r       bool
	forSure bool
}

type Processor[T any] interface {
	getSample() T
	process(text string) (T, bool)
}

type NlpService struct {
	processors map[string][]Processor[any]
}

func NewNlpService() *NlpService {
	processors := map[string][]Processor[any]{}
	return &NlpService{processors}
}

func (s *NlpService) process(ctx context.Context, text string, resultTypeName string) (any, bool) {
	if processors, ok := s.processors[resultTypeName]; ok {
		var (
			result any = nil
			ok         = false
		)
		for _, processor := range processors {
			result, ok = processor.process(text)
			if ok {
				return result, ok
			}
		}
	}
	return nil, false
}

func (s *NlpService) Use(processor Processor[any]) {
	resultType := reflect.TypeOf(processor.getSample())
	if _, ok := s.processors[resultType.Name()]; !ok {
		s.processors[resultType.Name()] = []Processor[any]{processor}
	} else {
		s.processors[resultType.Name()] = append(s.processors[resultType.Name()], processor)
	}
}

var BoolType = reflect.TypeOf(true).Name()

func (s *NlpService) GetBoolResult(ctx context.Context, text string) (bool, bool) {
	result, ok := s.process(ctx, text, BoolType)
	if ok {
		return result.(bool), ok
	} else {
		return false, ok
	}
}

var IntType = reflect.TypeOf(1).Name()

func (s *NlpService) GetIntResult(ctx context.Context, text string) (int, bool) {
	result, ok := s.process(ctx, text, IntType)
	if ok {
		return result.(int), ok
	} else {
		return 0, ok
	}
}
