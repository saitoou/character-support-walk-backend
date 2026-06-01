//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../test/mock/$GOPACKAGE/$GOFILE

package auth

import (
	"context"
)

type IDTokenPayload struct {
	Subject string
}

type IDTokenValidator interface {
	Validate(
		ctx context.Context,
		token string,
	) (*IDTokenPayload, error)
}
