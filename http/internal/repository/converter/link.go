package converter

import (
	"github.com/dfg007star/avito_informer/http/internal/model"
	repoModel "github.com/dfg007star/avito_informer/http/internal/repository/model"
	"time"
)

func RepoLinkToModel(link *repoModel.Link) *model.Link {
	createdAt := time.Now()
	if link.CreatedAt.Valid {
		createdAt = link.CreatedAt.Time
	}

	var parsedAt *time.Time
	if link.ParsedAt.Valid {
		t := link.ParsedAt.Time
		parsedAt = &t
	}

	items := make([]*model.Item, len(link.Items))
	for i, dbItem := range link.Items {
		created := time.Now()
		if dbItem.CreatedAt.Valid {
			created = dbItem.CreatedAt.Time
		}
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
			CreatedAt:   created,
		}
	}

	return &model.Link{
		ID:        link.ID,
		Name:      link.Name,
		Url:       link.Url,
		CreatedAt: createdAt,
		ParsedAt:  parsedAt,
		Items:     items,
	}
}
