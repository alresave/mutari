package matching

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

type Candidate struct {
	UserID int64   `json:"user_id"`
	Score  float64 `json:"score"`
}

func (r *Repository) GetCandidates(ctx context.Context, userID int64) ([]Candidate, error) {
	fmt.Println("***************************NEW BUILD444444444444444***********************************")
	query := `
	SELECT 
    u2.user_id,
    (
        0.65 * (1 - (u1.identity_vector <=> u2.preference_vector)) +
        0.35 * (1 - (u2.identity_vector <=> u1.preference_vector))
    ) AS score
FROM user_vectors u1
CROSS JOIN user_vectors u2
WHERE u1.user_id = $1
  AND u2.user_id <> u1.user_id
ORDER BY score DESC
LIMIT 10;
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Candidate

	for rows.Next() {
		var c Candidate
		if err := rows.Scan(&c.UserID, &c.Score); err != nil {
			return nil, err
		}
		results = append(results, c)
	}

	return results, nil
}
