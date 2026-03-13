package gin

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"github.com/gin-gonic/gin"
)

func parseRequestBody[T any](ctx *gin.Context) (T, error) {
	var zero T

	whatType := reflect.TypeOf((*T)(nil)).Elem()
	if whatType.Kind() != reflect.Struct {
		return zero, fmt.Errorf("parseRequestBody requires a struct type, got %s", whatType.Kind())
	}

	data, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return zero, fmt.Errorf("failed to read request body: %w", err)
	}

	var requestBody T
	err = json.Unmarshal(data, &requestBody)
	if err != nil {
		return zero, fmt.Errorf("failed to unmarshal request body: %w", err)
	}

	return requestBody, nil
}
