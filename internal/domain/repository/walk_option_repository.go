//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../test/mock/$GOPACKAGE/$GOFILE
package repository

import (
	"context"

	"github.com/chocoko/character-support-walk-backend/internal/domain/entity"
)

type WalkOptionRepository interface {
	SelectOption(ctx context.Context) ([]entity.WalkOption, error)
}
