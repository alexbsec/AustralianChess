package engine

import "math"

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

const (
	maxScore = math.MaxInt32
	minScore = math.MinInt32

	ttExact      = 0
	ttLowerBound = 1
	ttUpperBound = 2
)
