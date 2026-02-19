package m_product

const (
	ProductID            = "product_id"
	Name                 = "name"
	Description          = "description"
	Category             = "category"
	BasePriceNumerator   = "base_price_numerator"
	BasePriceDenominator = "base_price_denominator"
	DiscountPercent      = "discount_percent"
	DiscountStartDate    = "discount_start_date"
	DiscountEndDate      = "discount_end_date"
	Status               = "status"
	Version              = "version"
	CreatedAt            = "created_at"
	UpdatedAt            = "updated_at"
	ArchivedAt           = "archived_at"
)

const TableName = "products"

func Columns() []string {
	return []string{
		ProductID, Name, Description, Category,
		BasePriceNumerator, BasePriceDenominator,
		DiscountPercent, DiscountStartDate, DiscountEndDate,
		Status, Version, CreatedAt, UpdatedAt, ArchivedAt,
	}
}
