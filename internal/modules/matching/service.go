package matching

import "context"

type Service struct {
	repo *Repository
}

var blockNames = map[int]string{
	0: "Core Values",
	1: "Lifestyle Rhythm",
	2: "Emotional Stability",
	3: "Relational Architecture",
	4: "Ambition",
	5: "Growth",
	6: "Intimacy",
	7: "Cognitive Style",
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetMatches(ctx context.Context, userID int64) ([]Candidate, error) {
	return s.repo.GetMCPMatches(ctx, userID)
}

func (s *Service) SubmitMCPAnswers(ctx context.Context, userID int64, answers []MCPAnswer) error {

	err := s.repo.StoreMCPAnswers(ctx, userID, answers)
	if err != nil {
		return err
	}

	return s.repo.BuildIdentityVector(ctx, userID)
}

func (s *Service) GetMCPQuestions(ctx context.Context) ([]MCPQuestion, error) {
	return s.repo.GetMCPQuestions(ctx)
}
