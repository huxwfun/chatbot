package nlp

import (
	"strconv"

	"github.com/cdipaolo/sentiment"
)

type StaticBoolProcessor struct {
	Processor[bool]
}

func (s StaticBoolProcessor) getSample() any {
	r := true
	return r
}

func (s StaticBoolProcessor) process(text string) (any, bool) {
	result, ok := false, false
	if text == "ok" || text == "y" || text == "yes" || text == "sure" {
		result = true
		ok = true
	} else if text == "no" || text == "n" {
		result = false
		ok = true
	}
	return result, ok
}

type SentimentBoolProcessor struct {
	Processor[bool]
	model sentiment.Models
}

func NewSentimentBoolProcessor() SentimentBoolProcessor {
	model, err := sentiment.Restore()
	if err != nil {
		panic(err)
	}
	return SentimentBoolProcessor{
		model: model,
	}
}
func (s SentimentBoolProcessor) getSample() any {
	return true
}

func (s SentimentBoolProcessor) process(text string) (any, bool) {
	analysis := s.model.SentimentAnalysis(text, sentiment.English)
	score := analysis.Score
	return score > 0, true
}

type StaticIntProcessor struct {
	Processor[int]
}

func (s StaticIntProcessor) getSample() any {
	return 1
}

func (s StaticIntProcessor) process(text string) (any, bool) {
	result, err := strconv.Atoi(text)
	if err != nil {
		return 0, false
	}
	if result > 5 {
		result = 5
	} else if result < 1 {
		result = 1
	}
	return result, true
}
