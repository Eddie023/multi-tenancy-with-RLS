package store

import (
	"context"
	"errors"
)

type ctxKey int

const tenantKey ctxKey = 1

func SetTenantID(ctx context.Context, tenantID int) context.Context {
	return context.WithValue(ctx, tenantKey, tenantID)
}

func GetTenantID(ctx context.Context) (int, error) {
	v, ok := ctx.Value(tenantKey).(int)
	if !ok {
		return 0, errors.New("tenant ID not set in context")
	}

	return v, nil
}
