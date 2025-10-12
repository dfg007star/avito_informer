package item

import (
	"context"
	"fmt"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/dfg007star/avito_informer/collector/internal/model"
)

func (r *repository) CreateItems(ctx context.Context, items []*model.Item) error {
	if len(items) == 0 {
		return nil
	}

	var uids []string
	for _, item := range items {
		uids = append(uids, item.Uid)
	}

	existingUIDs := make(map[string]struct{})
	selectQuery, selectArgs, err := squirrel.Select("uid").From("items").Where(squirrel.Eq{"uid": uids}).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build select query: %w", err)
	}

	rows, err := r.db.Query(ctx, selectQuery, selectArgs...)
	if err != nil {
		return fmt.Errorf("failed to query existing uids: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			return fmt.Errorf("failed to scan existing uid: %w", err)
		}
		existingUIDs[uid] = struct{}{} // Add existing UID to map
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows iteration error: %w", err)
	}

	var newItems []*model.Item
	for _, item := range items {
		if _, exists := existingUIDs[item.Uid]; !exists {
			newItems = append(newItems, item)
		}
	}

	if len(newItems) == 0 {
		log.Println("no new items to create")
		return nil
	}

	log.Printf("creating %d new items", len(newItems))

	queryBuilder := squirrel.Insert("items").
		Columns("uid", "link_id", "title", "price", "description", "url", "preview_url")

	for _, item := range newItems {
		queryBuilder = queryBuilder.Values(item.Uid, item.LinkId, item.Title, item.Price, item.Description, item.Url, item.PreviewUrl)
	}

	query, args, err := queryBuilder.
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build insert query: %w", err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to exec insert query: %w", err)
	}

	return nil
}
