package auth

import (
	"context"

	"cloud.google.com/go/auth/credentials/idtoken"
	domainauth "github.com/chocoko/character-support-walk-backend/internal/domain/auth"
)

type GoogleIDTokenVerifier struct {
	audience string
}

func NewGoogleIDTokenVerifier(audience string) *GoogleIDTokenVerifier {
	return &GoogleIDTokenVerifier{
		audience: audience,
	}
}

func (v *GoogleIDTokenVerifier) Validate(
	ctx context.Context,
	token string,
) (*domainauth.IDTokenPayload, error) {
	payload, err := idtoken.Validate(ctx, token, v.audience)
	if err != nil {
		return nil, err
	}

	return &domainauth.IDTokenPayload{
		Subject: payload.Subject,
	}, nil
}
