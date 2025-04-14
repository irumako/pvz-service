package usecase

import (
	"context"
	"github.com/avito-tech/go-transaction-manager/trm/v2"
)

type FakeManager struct{}

func (fm *FakeManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func (fm *FakeManager) DoWithSettings(ctx context.Context, _ trm.Settings, fn func(ctx context.Context) error) error {
	return fn(ctx)
}
