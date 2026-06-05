package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/chocoko/character-support-walk-backend/internal/domain/entity"
	"github.com/chocoko/character-support-walk-backend/internal/domain/repository"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"golang.org/x/sync/errgroup"
)

type HomeOutput struct {
	Walker    HomeWalkerOutput     `json:"walker"`
	Character HomeCharacterOutput  `json:"character"`
	TodayWalk *HomeTodayWalkOutput `json:"today_walk"`
}

type HomeWalkerOutput struct {
	Nickname string `json:"nickname"`
}

type HomeCharacterOutput struct {
	SupporterType string `json:"supporter_type"`
}

type HomeTodayWalkOutput struct {
	WalkID uuid.UUID         `json:"walk_id"`
	Status entity.WalkStatus `json:"status"`
}

type HomeUsecase struct {
	userRepo      repository.UserRepository
	characterRepo repository.CharacterRepository
	walkRepo      repository.WalkRepository
}

func NewHomeUsecase(
	userRepo repository.UserRepository,
	characterRepo repository.CharacterRepository,
	walkRepo repository.WalkRepository,
) *HomeUsecase {
	return &HomeUsecase{userRepo: userRepo, characterRepo: characterRepo, walkRepo: walkRepo}
}

func (uc *HomeUsecase) GetHome(ctx context.Context, userID string) (HomeOutput, error) {
	ctx, span := otel.Tracer("character-support-walk-api").Start(
		ctx,
		"HomeUsecase.GetHome",
	)
	defer span.End()

	// JSTの"今日"に対応するUTC範囲で検索をしたいため
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return HomeOutput{}, fmt.Errorf("failed to load location: %w", err)
	}
	nowInJST := time.Now().In(loc)
	todayStartJST := time.Date(nowInJST.Year(), nowInJST.Month(), nowInJST.Day(), 0, 0, 0, 0, loc)
	tomorrowStartJST := todayStartJST.AddDate(0, 0, 1)

	searchFromUTC := todayStartJST.UTC()
	searchToUTC := tomorrowStartJST.UTC()

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return HomeOutput{}, fmt.Errorf("failed to parse uuid :%v", err)
	}

	var (
		user      *entity.User
		character *entity.Character
		result    *entity.Walk
	)

	// ユーザー数増加の場合、コネクションプーリングの要因になるので普通に戻す。
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		var err error
		user, err = uc.userRepo.FindByID(ctx, parsedUserID)
		if err != nil {
			return fmt.Errorf("failed to find user: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		var err error
		character, err = uc.characterRepo.FindByUserID(ctx, parsedUserID)
		if err != nil {
			return fmt.Errorf("failed to find character: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		var err error
		result, err = uc.walkRepo.FindTodayWalkByUserID(ctx, parsedUserID, searchFromUTC, searchToUTC)
		if err != nil {
			return fmt.Errorf("failed to find today walk: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return HomeOutput{}, err
	}
	if user == nil || user.DeletedAt != nil {
		return HomeOutput{}, fmt.Errorf("inactive user: %w", ErrUnauthorized)
	}
	if character == nil {
		return HomeOutput{}, fmt.Errorf("inactive user: %w", ErrInternalServer)
	}

	var todayWalk *HomeTodayWalkOutput

	if result != nil {
		todayWalk = &HomeTodayWalkOutput{
			WalkID: result.ID,
			Status: result.Status,
		}
	}

	userOutput := HomeWalkerOutput{
		Nickname: user.Nickname,
	}

	characterOutput := HomeCharacterOutput{
		SupporterType: string(character.SupporterType),
	}

	return HomeOutput{
		Walker:    userOutput,
		Character: characterOutput,
		TodayWalk: todayWalk,
	}, nil
}
