package store

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/shahin-bayat/go-scraper/model"
	"github.com/shahin-bayat/go-scraper/util"
)

type Store struct {
	db *sql.DB
}

func NewPostgresStore() (*Store, error) {

	host := util.GetEnvVariable("DB_HOST")
	user := util.GetEnvVariable("DB_USER")
	password := util.GetEnvVariable("DB_PASSWORD")
	dbname := util.GetEnvVariable("DB_NAME")
	port := util.GetEnvVariable("DB_PORT")
	sslmode := util.GetEnvVariable("DB_SSLMODE")
	timezone := util.GetEnvVariable("DB_TIMEZONE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", host, user, password, dbname, port, sslmode, timezone)

	// if _ "github.com/lib/pq" is not imported, the following line will throw an error
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, errors.New("failed to connect database: " + err.Error())
	}

	return &Store{db: db}, nil
}

func (s *Store) createCategoryTable() error {
	query := `CREATE TABLE IF NOT EXISTS categories (
		id SERIAL PRIMARY KEY,
		text VARCHAR(50) NOT NULL,
		category_key VARCHAR(50) UNIQUE NOT NULL,
		created_at TIMEStAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMEStAMP DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := s.db.Exec(query)
	return err
}

func (s *Store) createQuestionTable() error {
	query := `CREATE TABLE IF NOT EXISTS questions (
		id SERIAL PRIMARY KEY,
		category_id INTEGER NOT NULL REFERENCES categories(id),
		image_path VARCHAR(255),
		text VARCHAR(255),
		question_number VARCHAR(50) NOT NULL,
		question_key VARCHAR(50) UNIQUE NOT NULL,
		is_fetched BOOLEAN DEFAULT FALSE,
		created_at TIMEStAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMEStAMP DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := s.db.Exec(query)
	return err
}

func (s *Store) createAnswerTable() error {
	query := `CREATE TABLE IF NOT EXISTS answers (
		id SERIAL PRIMARY KEY,
		question_id INTEGER NOT NULL REFERENCES questions(id),
		text VARCHAR(255) NOT NULL,
		is_correct BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := s.db.Exec(query)
	return err
}

func (s *Store) Init() error {
	if err := s.createCategoryTable(); err != nil {
		return err
	}
	if err := s.createQuestionTable(); err != nil {
		return err
	}
	if err := s.createAnswerTable(); err != nil {
		return err
	}
	return nil
}

func (s *Store) CreateCategory(category *model.Category) error {
	query := `INSERT INTO categories (text, category_key) VALUES ($1, $2) ON CONFLICT (category_key) DO NOTHING`
	_, err := s.db.Exec(query, category.Text, category.CategoryKey)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetCategoryByCategoryKey(categoryKey string) (*model.Category, error) {
	var category model.Category
	query := `SELECT * FROM categories WHERE category_key = $1`
	err := s.db.QueryRow(query, categoryKey).Scan(&category.ID, &category.Text, &category.CategoryKey, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (s *Store) GetCategories() ([]model.Category, error) {
	var categories []model.Category
	query := `SELECT * FROM categories`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var category model.Category
		err := rows.Scan(&category.ID, &category.Text, &category.CategoryKey, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

// TODO: for debugging purposes, remove this function later
func (s *Store) GetCategoryByText(text string) (*model.Category, error) {
	var category model.Category
	query := `SELECT * FROM categories WHERE text = $1`
	err := s.db.QueryRow(query, text).Scan(&category.ID, &category.Text, &category.CategoryKey, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (s *Store) CreateQuestion(question *model.Question) error {
	query := `INSERT INTO questions (category_id, image_path, text, question_number, question_key) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (question_key) DO NOTHING`
	_, err := s.db.Exec(query, question.CategoryID, question.ImagePath, question.Text, question.QuestionNumber, question.QuestionKey)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetQuestionByQuestionKey(questionKey string) (*model.Question, error) {
	var question model.Question
	query := `SELECT * FROM questions WHERE question_key = $1`
	err := s.db.QueryRow(query, questionKey).Scan(&question.ID, &question.CategoryID, &question.ImagePath, &question.Text, &question.QuestionNumber, &question.QuestionKey, &question.IsFetched, &question.CreatedAt, &question.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &question, nil
}

func (s *Store) GetQuestionsByCategoryId(categoryId uint) ([]model.Question, error) {
	var questions []model.Question
	query := `SELECT * FROM questions WHERE category_id = $1`
	rows, err := s.db.Query(query, categoryId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var question model.Question
		err := rows.Scan(&question.ID, &question.CategoryID, &question.ImagePath, &question.Text, &question.QuestionNumber, &question.QuestionKey, &question.IsFetched, &question.CreatedAt, &question.UpdatedAt)
		if err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}
	return questions, nil
}

func (s *Store) UpdateQuestion(question *model.Question) error {
	query := `UPDATE questions SET image_path = $1, is_fetched = $2 WHERE id = $3`
	_, err := s.db.Exec(query, question.ImagePath, question.IsFetched, question.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) CreateAnswer(answer *model.Answer) error {
	query := `INSERT INTO answers (question_id, text, is_correct) VALUES ($1, $2, $3)`
	_, err := s.db.Exec(query, answer.QuestionID, answer.Text, answer.IsCorrect)
	if err != nil {
		return err
	}
	return nil
}
