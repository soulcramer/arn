package arn

// Purchase ...
type Purchase struct {
	ID       string `json:"id"`
	UserID   string `json:"userId"`
	ItemID   string `json:"itemId"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
	Currency string `json:"currency"`
	Date     string `json:"date"`
}

// Item returns the item the user bought.
func (purchase *Purchase) Item() *Item {
	item, _ := GetItem(purchase.ItemID)
	return item
}

// NewPurchase creates a new Purchase object with a generated ID.
func NewPurchase(userID string, itemID string, quantity int, price int, currency string) *Purchase {
	return &Purchase{
		ID:       GenerateID("Purchase"),
		UserID:   userID,
		ItemID:   itemID,
		Quantity: quantity,
		Price:    price,
		Currency: currency,
		Date:     DateTimeUTC(),
	}
}

// StreamPurchases returns a stream of all anime.
func StreamPurchases() (chan *Purchase, error) {
	objects, err := DB.All("Purchase")
	return objects.(chan *Purchase), err
}

// MustStreamPurchases returns a stream of all anime.
func MustStreamPurchases() chan *Purchase {
	stream, err := StreamPurchases()
	PanicOnError(err)
	return stream
}

// AllPurchases returns a slice of all anime.
func AllPurchases() ([]*Purchase, error) {
	var all []*Purchase

	stream, err := StreamPurchases()

	if err != nil {
		return nil, err
	}

	for obj := range stream {
		all = append(all, obj)
	}

	return all, nil
}
