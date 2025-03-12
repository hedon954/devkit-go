package idempotence

type ErrCode uint

const (
	ErrCodeSuccess        ErrCode = 20000
	ErrCodeDuplicateReq   ErrCode = 40001
	ErrCodeFailedCanRetry ErrCode = 40002
	ErrCodeFailedNoRetry  ErrCode = 40003
)

type ResponseCode uint

const (
	ResponseCodeSuccess ResponseCode = 200
	ResponseCodeFailed  ResponseCode = 500
)
