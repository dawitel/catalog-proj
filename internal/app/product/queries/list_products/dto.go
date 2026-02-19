package list_products

type Item struct {
	ID                        string
	Name                      string
	Description               string
	Category                  string
	BasePriceNumerator        int64
	BasePriceDenominator      int64
	EffectivePriceNumerator   int64
	EffectivePriceDenominator int64
	Status                    string
}

type Result struct {
	Items     []Item
	NextToken string
}
