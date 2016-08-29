package lazascrap

type ProductItem struct {
	Title           string
	Price           string
	DiscountedPrice string
	Image           string
}

type Product struct {
	Items      []ProductItem
	TotalItems int
}
