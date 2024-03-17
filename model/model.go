package model

import (
	"log"
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID          uint `gprm:"primaryKey"`
	Text        string
	CategoryKey string     // 4, 5, 6, ...
	Questions   []Question `gorm:"foreignKey:CategoryID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Question struct {
	ID             uint `gprm:"primaryKey"`
	CategoryID     uint
	ImagePath      string
	Text           *string
	QuestionNumber string   // 1, 2, 3, ...
	QuestionKey    string   // 1732, 1733, 1734, ...
	IsFetched      bool     `gorm:"default:false"`
	Answers        []Answer `gorm:"foreignKey:QuestionID"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Answer struct {
	ID         uint `gprm:"primaryKey"`
	QuestionID uint
	Text       string
	IsCorrect  bool `gorm:"default:false"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func CreateCategory(text, categoryKey string, db *gorm.DB) error {
	modelCategory := Category{
		Text:        text,
		CategoryKey: categoryKey,
	}
	result := db.FirstOrCreate(&modelCategory, Category{CategoryKey: categoryKey})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func CreateQuestion(questionNumber, questionKey, categoryKey string, db *gorm.DB) error {
	var categoryId uint
	result := db.Model(&Category{}).Where("category_key = ?", categoryKey).Select("id").Find(&categoryId)
	if result.Error != nil {
		log.Fatalf("Error fetching category id:%s", result.Error)
		return result.Error
	}

	modelQuestion := Question{
		QuestionNumber: questionNumber,
		QuestionKey:    questionKey,
		CategoryID:     categoryId,
	}
	result = db.FirstOrCreate(&modelQuestion, Question{QuestionKey: questionKey})
	if result.Error != nil {
		log.Fatalf("Error creating question:%s", result.Error)
		return result.Error
	}

	return nil
}

func UpdateQuestion(questionKey, imagePath string, db *gorm.DB) error {
	result := db.Model(&Question{}).Where("question_key = ?", questionKey).Update("image_path", imagePath)
	if result.Error != nil {
		log.Fatalf("Error updating question image path:%s", result.Error)
		return result.Error
	}
	return nil
}

func CreateAnswer(questionKey, text string, isCorrect bool, db *gorm.DB) error {
	var questionId uint
	result := db.Model(&Question{}).Where("question_key = ?", questionKey).Select("id").Find(&questionId)
	if result.Error != nil {
		log.Fatalf("Error fetching question id:%s", result.Error)
		return result.Error
	}

	modelAnswer := Answer{
		QuestionID: questionId,
		Text:       text,
		IsCorrect:  isCorrect,
	}
	result = db.Create(&modelAnswer)
	if result.Error != nil {
		log.Fatalf("Error creating answer:%s", result.Error)
		return result.Error
	}

	// update questions set is_fetched = true where question_key = questionKey
	result = db.Model(&Question{}).Where("question_key = ?", questionKey).Update("is_fetched", true)
	if result.Error != nil {
		log.Fatalf("Error updating question is_fetched:%s", result.Error)
		return result.Error
	}

	return nil
}
