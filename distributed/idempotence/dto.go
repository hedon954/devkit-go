package idempotence

type Response struct {
	Code    ResponseCode `json:"code"`
	ErrCode ErrCode      `json:"err_code"`
	Msg     string       `json:"msg"`
	Data    any          `json:"data"`
}

type CreateOrderReq struct {
	UserID string `json:"user_id" binding:"required"`
	ItemID string `json:"item_id" binding:"required"`
}

type PayOrderReq struct {
	UserID  string `json:"user_id" binding:"required"`
	OrderID string `json:"order_id" binding:"required"`
}
