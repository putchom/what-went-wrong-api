package services

import (
	"fmt"
)

type AIService interface {
	GenerateExcuse(tone string, context string) ([]string, error)
}

type MockAIService struct{}

func NewMockAIService() *MockAIService {
	return &MockAIService{}
}

func (s *MockAIService) GenerateExcuse(tone string, context string) ([]string, error) {
	// Mock implementation returning dummy excuses based on tone
	var candidates []string
	switch tone {
	case "surreal":
		candidates = []string{
			fmt.Sprintf("重力が強すぎて、%s ができませんでした。", context),
			fmt.Sprintf("時空の歪みにより、%s という概念が消滅していました。", context),
		}
	case "philosophical":
		candidates = []string{
			fmt.Sprintf("%s をすることは、宇宙の真理に反すると感じました。", context),
			fmt.Sprintf("存在そのものが %s を拒否していました。", context),
		}
	default:
		candidates = []string{
			fmt.Sprintf("なんとなく %s ができませんでした。", context),
			fmt.Sprintf("今日は %s の日ではありませんでした。", context),
		}
	}
	return candidates, nil
}
