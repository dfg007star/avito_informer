package item

import (
	"context"

	"github.com/dfg007star/avito_informer/notification/internal/model"
	repoModel "github.com/dfg007star/avito_informer/notification/internal/repository/model"
	"gorm.io/gorm"
)

func (r *repository) UpdateItem(ctx context.Context, item *model.Item) error {
	repoItem := repoModel.Item{
		Model:       gorm.Model{ID: uint(item.ID)},
		LinkId:      item.LinkId,
		Uid:         item.Uid,
		Title:       item.Title,
		Description: item.Description,
		Url:         item.Url,
		PreviewUrl:  item.PreviewUrl,
		Price:       item.Price,
		NeedNotify:  item.NeedNotify,
	}

	err := r.db.WithContext(ctx).Save(&repoItem).Error
	if err != nil {
		return err
	}

	return nil
}
