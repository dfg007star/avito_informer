package converter

import (
	"github.com/dfg007star/avito_informer/http/internal/model"
	repoModel "github.com/dfg007star/avito_informer/http/internal/repository/model"
	"time"
)

func RepoItemToModel(item *repoModel.Item) *model.Item {
	createdAt := time.Now()
	if item.CreatedAt.Valid {
		createdAt = item.CreatedAt.Time
	}

	return &model.Item{
		ID:          item.ID,
		LinkId:      item.LinkId,
		Uid:         item.Uid,
		Title:       item.Title,
		Description: item.Description,
		Url:         item.Url,
		PreviewUrl:  item.PreviewUrl,
		Price:       item.Price,
		NeedNotify:  item.NeedNotify,
		CreatedAt:   createdAt,
	}
}
