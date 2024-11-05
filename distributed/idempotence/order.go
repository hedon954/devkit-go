package idempotence

type Order struct {
	ID        string         `json:"id"`
	KeyParams OrderKeyParams `json:"key_params"`
	Items     []string       `json:"items"`
	State     OrderState     `json:"state"`
}

type OrderKeyParams struct {
	UserID   string  `json:"user_id"`
	ItemID   string  `json:"item_id"`
	Price    int64   `json:"price"`
	Discount float64 `json:"discount"`
	CouponID string  `json:"coupon_id"`
}

type OrderState int

const (
	OrderStateInit OrderState = iota
	OrderStatePending
	OrderStatePaid
	OrderStateSuccess
	OrderStateFailed
)
