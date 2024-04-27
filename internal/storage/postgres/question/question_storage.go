package question

import (
	"context"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/internal/entities"
	"github.com/go-park-mail-ru/2024_1_ResCogitans/utils/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QuestionStorage struct {
	db *pgxpool.Pool
}

func NewQuestionStorage(db *pgxpool.Pool) *QuestionStorage {
	return &QuestionStorage{
		db: db,
	}
}

func (qs *QuestionStorage) AddReview(review entities.Review) error {
	currentTime := time.Now()
	moscowLocation, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return err
	}
	// Преобразуем текущее время в часовой пояс Москвы
	currentTimeInMoscow := currentTime.In(moscowLocation)
	timePlusThreeHours := currentTimeInMoscow.Add(3 * time.Hour)

	query := `
        INSERT INTO quiz (user_id, rating, question_id, created_at)
        VALUES ($1, $2, $3, $4)
    `
	_, err = qs.db.Exec(context.Background(), query, review.UserID, review.Rating, review.QuestionID, timePlusThreeHours)
	return err
}

func (qs *QuestionStorage) GetQuestions() ([]entities.QuestionResponse, error) {
	var questions []*entities.QuestionResponse
	ctx := context.Background()
	err := pgxscan.Select(ctx, qs.db, &questions, `SELECT * FROM question`)
	if err != nil {
		logger.Logger().Error(err.Error())
		return []entities.QuestionResponse{}, err
	}

	var questionList []entities.QuestionResponse
	for _, q := range questions {
		questionList = append(questionList, *q)
	}

	return questionList, nil
}

func (qs *QuestionStorage) GetReview(userID int) ([]entities.Review, error) {
	var review []*entities.Review
	ctx := context.Background()
	err := pgxscan.Select(ctx, qs.db, &review, `SELECT * FROM quiz WHERE user_id = $1`, userID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return []entities.Review{}, err
	}

	var reviewList []entities.Review
	for _, r := range review {
		reviewList = append(reviewList, *r)
	}

	return reviewList, nil
}

func (qs *QuestionStorage) SetStat(userID int) ([]entities.Statistic, error) {
	var statistic []*entities.Statistic
	ctx := context.Background()
	err := pgxscan.Select(ctx, qs.db, &statistic, `SELECT q.text, r.rating, AVG(r.rating) 
	FROM quiz r 
	INNER JOIN question q ON r.question_id = q.id 
	WHERE user_id = $1 
	GROUP BY r.question_id`, userID)
	if err != nil {
		logger.Logger().Error(err.Error())
		return []entities.Statistic{}, err
	}

	var statisticList []entities.Statistic
	for _, s := range statistic {
		statisticList = append(statisticList, *s)
	}

	return statisticList, nil
}