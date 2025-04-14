package postgres

import (
	"context"
	"github.com/avito-tech/go-transaction-manager/trm/v2"
	trmSettings "github.com/avito-tech/go-transaction-manager/trm/v2/settings"
)

func getForUpdateClause(ctx context.Context) string {
	if _, ok := ctx.Value(trmSettings.DefaultCtxKey).(trm.Transaction); ok {
		return "FOR UPDATE"
	}
	return ""
}
