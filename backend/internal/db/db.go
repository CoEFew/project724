package db

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool
var ErrNoQuiz = errors.New("no quiz")

func Init(ctx context.Context) error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return fmt.Errorf("DATABASE_URL is empty")
	}
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return err
	}
	cfg.MaxConns = 5
	cfg.MinConns = 0
	cfg.MaxConnLifetime = time.Hour
	cfg.HealthCheckPeriod = 30 * time.Second

	p, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return err
	}
	if err := p.Ping(ctx); err != nil {
		p.Close()
		return err
	}
	pool = p
	return nil
}

func Close() {
	if pool != nil {
		pool.Close()
	}
}

type ScoreRow struct {
	Name      string
	Score     int
	GameName  string
	CreatedAt time.Time
}

func InsertScore(ctx context.Context, name string, score int, gamename string) error {
	_, err := pool.Exec(ctx,
		`INSERT INTO public.scores(name, score, gamename) VALUES ($1, $2, $3)`,
		name, score, gamename,
	)
	return err
}

func GetTopScores(ctx context.Context, gamename string, limit int) ([]ScoreRow, error) {
	rows, err := pool.Query(ctx, `
		SELECT name, score, gamename, created_at
		FROM public.scores
		WHERE ($1 = '' OR gamename = $1)
		ORDER BY score DESC, created_at ASC
		LIMIT $2
	`, gamename, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []ScoreRow
	for rows.Next() {
		var s ScoreRow
		if err := rows.Scan(&s.Name, &s.Score, &s.GameName, &s.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

/* ===== Quizzes (คงเดิม) ===== */

type QuizRow struct {
	Answer string
	Hint1  string
	Hint2  string
	Tier   int
}

func GetAllAnswers(ctx context.Context) ([]string, error) {
	rows, err := pool.Query(ctx, `SELECT answer FROM public.quizzes WHERE active`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []string
	for rows.Next() {
		var a string
		if err := rows.Scan(&a); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func GetRandomQuizByTier(ctx context.Context, tier int) (QuizRow, error) {
	var q QuizRow
	err := pool.QueryRow(ctx, `
		SELECT answer, hint1, hint2, tier
		FROM public.quizzes
		WHERE active AND tier = $1
		ORDER BY random()
		LIMIT 1
	`, tier).Scan(&q.Answer, &q.Hint1, &q.Hint2, &q.Tier)
	if errors.Is(err, pgx.ErrNoRows) {
		return q, ErrNoQuiz
	}
	return q, err
}

type QuizLite struct {
	Answer string
	Hint1  string
	Hint2  string
}

func GetAllQuizLite(ctx context.Context) ([]QuizLite, error) {
	rows, err := pool.Query(ctx, `
        SELECT answer, hint1, hint2
        FROM public.quizzes
        WHERE active = TRUE
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []QuizLite{}
	for rows.Next() {
		var q QuizLite
		if err := rows.Scan(&q.Answer, &q.Hint1, &q.Hint2); err != nil {
			return nil, err
		}
		out = append(out, q)
	}
	return out, rows.Err()
}
