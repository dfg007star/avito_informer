package converter

import (
	"github.com/dfg007star/avito_informer/http/internal/model"
	repoModel "github.com/dfg007star/avito_informer/http/internal/repository/model"
	"time"
)

func RepoLinkToModel(link *repoModel.Link) *model.Link {
	var parsedAt *time.Time
	if link.ParsedAt.Valid {
		t := link.ParsedAt.Time
		parsedAt = &t
	}

	items := make([]*model.Item, len(link.Items))
	for i, dbItem := range link.Items {
		items[i] = &model.Item{
			ID:          int64(dbItem.ID),
			LinkId:      int64(dbItem.LinkId),
			Uid:         dbItem.Uid,
			Title:       dbItem.Title,
			Description: dbItem.Description,
			Url:         dbItem.Url,
			PreviewUrl:  dbItem.PreviewUrl,
			Price:       dbItem.Price,
			NeedNotify:  dbItem.NeedNotify,
			CreatedAt:   dbItem.CreatedAt,
		}
	}

	return &model.Link{
		ID:        int64(link.ID),
		Name:      link.Name,
		Url:       link.Url,
		CreatedAt: link.CreatedAt,
		ParsedAt:  parsedAt,
		Items:     items,
	}
}
