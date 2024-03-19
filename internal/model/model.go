package model

import (
	"time"
)

type Category struct {
	ID          uint
	Text        string
	CategoryKey string // 4, 5, 6, ...
	Questions   []Question
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Question struct {
	ID             uint
	ImagePath      string
	Text           *string
	QuestionNumber string // 1, 2, 3, ...
	QuestionKey    string // 1732, 1733, 1734, ...
	IsFetched      bool
	Answers        []Answer
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Answer struct {
	ID         uint
	QuestionID uint
	Text       string
	IsCorrect  bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type UpdateQuestionRequest struct {
	ImagePath string
	IsFetched bool
}

func CreateCategory(text, categoryKey string) *Category {
	modelCategory := Category{
		Text:        text,
		CategoryKey: categoryKey,
	}

	return &modelCategory
}

func CreateQuestion(questionNumber, questionKey string) *Question {
	modelQuestion := Question{
		QuestionNumber: questionNumber,
		QuestionKey:    questionKey,
	}
	return &modelQuestion
}

func UpdateQuestion(question *Question, updateQuestionRequest *UpdateQuestionRequest) *Question {
	if updateQuestionRequest.ImagePath != "" {
		question.ImagePath = updateQuestionRequest.ImagePath
	}
	if updateQuestionRequest.IsFetched {
		question.IsFetched = updateQuestionRequest.IsFetched
	}

	return question
}

func CreateAnswer(text string, isCorrect bool, question *Question) *Answer {
	modelAnswer := Answer{
		Text:       text,
		IsCorrect:  isCorrect,
		QuestionID: question.ID,
	}
	return &modelAnswer
}
