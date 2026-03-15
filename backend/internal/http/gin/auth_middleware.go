package gin

import (
	"github.com/alexbsec/AustralianChess/backend/internal/auth"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authorizer auth.IAuthorizer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session, err := authorizer.Authorize(ctx.Request.Context(), ctx.Request)
		if err != nil {
			ctx.AbortWithStatus(401)
			return
		}

		ctx.Set("userId", session.UserId)
		ctx.Set("sessionId", session.Id)
		ctx.Next()
	}
}

