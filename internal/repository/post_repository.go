// internal/repository/post_repository.go
package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/you/linkedinify/internal/model"
)

type PostRepository interface {
	Save(ctx context.Context, p *model.LinkedInPost) error
	ListByUser(ctx context.Context, userID uuid.UUID, limit int) ([]model.LinkedInPost, error)
}

type postRepo struct{ db *bun.DB }

func NewPostRepo(db *bun.DB) PostRepository { return &postRepo{db} }

func (p *postRepo) Save(ctx context.Context, post *model.LinkedInPost) error {
	_, err := p.db.NewInsert().Model(post).Exec(ctx)
	return err
}

func (p *postRepo) ListByUser(ctx context.Context, userID uuid.UUID, limit int) ([]model.LinkedInPost, error) {
	var posts []model.LinkedInPost
	err := p.db.NewSelect().
		Model(&posts).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Scan(ctx)
	return posts, err
}
