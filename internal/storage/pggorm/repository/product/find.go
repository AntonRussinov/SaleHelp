package product

import (
	"context"
	"fmt"
	"salehelp/domain"
	"salehelp/internal/storage/pggorm/model"

	"github.com/jackc/pgtype"
	"gorm.io/gorm"
)

func (r *Repository) GetByID(ctx context.Context, id int) (*domain.Product, error) {
	var m model.Product

	if err := r.DB.WithContext(ctx).
		First(&m, "id = ?", id).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("product get by id: %w", err)
	}

	return m.ToDomain(), nil
}

// GetByShopAndExternalID получает продукт по (shop_id, external_id)
// это твой уникальный ключ
func (r *Repository) GetByShopAndExternalID(
	ctx context.Context,
	shopID int,
	externalID string,
) (*domain.Product, error) {
	var m model.Product

	if err := r.DB.WithContext(ctx).
		First(&m, "shop_id = ? AND external_id = ?", shopID, externalID).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, fmt.Errorf("product get by shop_id+external_id: %w", err)
	}

	return m.ToDomain(), nil
}

// ListByShop список продуктов конкретного магазина c пагинацией
func (r *Repository) ListByShop(
	ctx context.Context,
	shopID int,
	limit, offset int,
) ([]*domain.Product, error) {
	var models []model.Product

	q := r.DB.WithContext(ctx).
		Where("shop_id = ?", shopID).
		Order("id ASC")

	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}

	if err := q.Find(&models).Error; err != nil {
		return nil, fmt.Errorf("product list by shop: %w", err)
	}

	result := make([]*domain.Product, 0, len(models))
	for i := range models {
		result = append(result, models[i].ToDomain())
	}

	return result, nil
}

func (r *Repository) ListIDsByShopID(ctx context.Context, shopID int) ([]int, error) {
	type row struct {
		ID int `gorm:"column:id"`
	}

	var rows []row
	if err := r.DB.WithContext(ctx).
		Raw(`SELECT id FROM products WHERE shop_id = ? ORDER BY id`, shopID).
		Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("product list ids by shop_id: %w", err)
	}

	out := make([]int, 0, len(rows))
	for _, rr := range rows {
		out = append(out, rr.ID)
	}
	return out, nil
}

func (r *Repository) GetByIDs(ctx context.Context, ids []int) (map[int]*domain.Product, error) {
	out := make(map[int]*domain.Product, len(ids))
	if len(ids) == 0 {
		return out, nil
	}

	arr := pgtype.Int8Array{}
	if err := arr.Set(ids); err != nil {
		return nil, fmt.Errorf("product get by ids: set array: %w", err)
	}

	var rows []model.Product
	if err := r.DB.WithContext(ctx).
		Where("id = ANY(?)", &arr).
		Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("product get by ids: %w", err)
	}

	for i := range rows {
		d := rows[i].ToDomain()
		out[d.ID] = d
	}

	return out, nil
}

// GetTitlesByBarcodes returns a sku→title map for the given barcodes belonging to the shop.
// Used as fallback when NmID from the Statistics API doesn't match the product's external_id.
func (r *Repository) GetTitlesByBarcodes(ctx context.Context, shopID int, barcodes []string) (map[string]string, error) {
	if len(barcodes) == 0 {
		return map[string]string{}, nil
	}

	type row struct {
		SKU   string `gorm:"column:sku"`
		Title string `gorm:"column:title"`
	}

	var rows []row
	if err := r.DB.WithContext(ctx).
		Table("product_size_skus pssk").
		Select("pssk.sku, p.title").
		Joins("JOIN product_sizes ps ON ps.id = pssk.product_size_id").
		Joins("JOIN products p ON p.id = ps.product_id").
		Where("p.shop_id = ? AND pssk.sku IN ?", shopID, barcodes).
		Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("product get titles by barcodes: %w", err)
	}

	out := make(map[string]string, len(rows))
	for _, row := range rows {
		if row.Title != "" {
			out[row.SKU] = row.Title
		}
	}
	return out, nil
}

func (r *Repository) CountByShop(ctx context.Context, shopID int) (int, error) {
	var cnt int64
	if err := r.DB.WithContext(ctx).
		Table("products").
		Where("shop_id = ?", shopID).
		Count(&cnt).Error; err != nil {
		return 0, fmt.Errorf("count products by shop: %w", err)
	}
	return int(cnt), nil
}
