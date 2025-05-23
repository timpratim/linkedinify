// internal/service/linkedin_service.go
package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/you/linkedinify/internal/ai"
	"github.com/you/linkedinify/internal/model"
	"github.com/you/linkedinify/internal/repository"
)

type LinkedInService struct {
	ai    ai.Client
	posts repository.PostRepository
}

func NewLinkedIn(ai ai.Client, pr repository.PostRepository) *LinkedInService {
	return &LinkedInService{ai: ai, posts: pr}
}

func (l *LinkedInService) Transform(ctx context.Context, userID uuid.UUID, text string) (string, error) {
	out, err := l.ai.Transform(ctx, text)
	if err != nil {
		return "", err
	}
	post := &model.LinkedInPost{
		ID:         uuid.New(),
		UserID:     userID,
		InputText:  text,
		OutputText: out,
	}
	if err = l.posts.Save(ctx, post); err != nil {
		return "", err
	}
	return out, nil
}

func (l *LinkedInService) History(ctx context.Context, userID uuid.UUID) ([]model.LinkedInPost, error) {
	return l.posts.ListByUser(ctx, userID, 20)
}
