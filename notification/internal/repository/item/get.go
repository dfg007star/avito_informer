package item

import (
	"context"

	"github.com/dfg007star/avito_informer/notification/internal/model"
	"github.com/dfg007star/avito_informer/notification/internal/repository/converter"
	repoModel "github.com/dfg007star/avito_informer/notification/internal/repository/model"
)

func (r *repository) GetItems(ctx context.Context) ([]*model.Item, error) {
	var items []*repoModel.Item

	err := r.db.WithContext(ctx).Find(&items).Error
	if err != nil {
		return nil, err
	}

	return converter.RepoItemToModelList(items), nil
}

func (r *repository) GetNotNotifiedItems(ctx context.Context) ([]*model.Item, error) {
	var items []*repoModel.Item

	err := r.db.WithContext(ctx).Where("need_notify = ?", false).Find(&items).Error
	if err != nil {
		return nil, err
	}

	return converter.RepoItemToModelList(items), nil
}
