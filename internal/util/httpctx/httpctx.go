package httpctx

import (
	"context"
	"fmt"

	"github.com/escoutdoor/study-platform/internal/util/token"
)

type ctxKey string

const (
	UserIDContextKey ctxKey = "id"

	RolesContextKey ctxKey = "roles"
)

func GetID(ctx context.Context) (int, error) {
	val := ctx.Value(UserIDContextKey)
	if val == nil {
		return 0, fmt.Errorf("userId missing in context")
	}

	id, ok := val.(int)
	if !ok {
		return 0, fmt.Errorf("userId cast failed")
	}

	if id == 0 {
		return 0, fmt.Errorf("userId is empty")
	}

	return id, nil
}

func GetRoles(ctx context.Context) ([]token.Role, error) {
	roles, ok := ctx.Value(RolesContextKey).([]token.Role)
	if !ok {
		return nil, fmt.Errorf("user roles not found in context")
	}

	return roles, nil
}
