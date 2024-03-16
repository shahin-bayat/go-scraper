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

func SaveQuestion(questionNumber, questionKey string, categoryKey string, db *gorm.DB) error {
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
	result = db.Create(&modelQuestion)
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
