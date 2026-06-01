//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=../../test/mock/$GOPACKAGE/$GOFILE
package repository

import "context"

type Transaction interface {
	RunInTx(ctx context.Context, fn func(ctx context.Context) error) error
}
