package engine

type ValidationType int

const (
	SuccessResponse ValidationType = iota
	FailureResponse
	ErrorResponse
)

type ValidationResponse struct {
	Type    ValidationType
	Message string
}
