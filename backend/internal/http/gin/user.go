package gin

import (
	"context"
	"net/http"

	"github.com/alexbsec/AustralianChess/backend/users"
	"github.com/gin-gonic/gin"
)

func UserHandler(
	ctx context.Context,
	routerGroup *gin.RouterGroup,
	userService users.IService,
) {
	routerGroup.Handle("POST", "/user/login", MakeUserLogin(userService))
	routerGroup.Handle("POST", "/user/register", MakeUserRegister(userService))
}

func MakeUserLogin(userService users.IService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userDTO, err := parseRequestBody[UserLoginDTO](ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		response, err := userService.LoginUser(ctx, userDTO.Username, userDTO.Password)
		if err != nil {
			code := errToHTTPStatus(err)
			ctx.JSON(code, gin.H{"error": "invalid credentials"})
			return
		}

		ctx.JSON(http.StatusOK, response)
	}
}

func MakeUserRegister(userService users.IService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userDTO, err := parseRequestBody[UserRegisterDTO](ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		if !arePasswordsEqual(userDTO.Password, userDTO.ConfirmPassword) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "passwords do not match"})
			return
		}

		response, err := userService.NewUser(ctx, userDTO.Username, userDTO.Password)
		if err != nil {
			code := errToHTTPStatus(err)
			ctx.JSON(code, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, response)
	}
}

func arePasswordsEqual(password, confirmPassword string) bool {
	return password == confirmPassword
}

func errToHTTPStatus(err error) int {
	switch err {
	case users.ErrInvalidCredentials:
		return http.StatusUnauthorized
	case users.ErrInvalidPlayerId:
		return http.StatusForbidden
	case users.ErrPasswordTooShort, users.ErrPlayerIdAlreadyExists:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
