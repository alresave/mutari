package matching

import "context"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetCandidates(ctx context.Context, userID int64) ([]Candidate, error) {
	return s.repo.GetCandidates(ctx, userID)
}
