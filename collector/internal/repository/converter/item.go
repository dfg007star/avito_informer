package converter

import (
	"github.com/dfg007star/avito_informer/collector/internal/model"
	repoModel "github.com/dfg007star/avito_informer/collector/internal/repository/model"
)

func RepoItemToModel(item *repoModel.Item) *model.Item {
	return &model.Item{
		ID:          int64(item.ID),
		LinkId:      item.LinkId,
		Uid:         item.Uid,
		Title:       item.Title,
		Description: item.Description,
		Url:         item.Url,
		PreviewUrl:  item.PreviewUrl,
		Price:       item.Price,
		IsNotify:    item.IsNotify,
		CreatedAt:   item.CreatedAt,
	}
}
