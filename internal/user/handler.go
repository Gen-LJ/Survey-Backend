package user

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"survey-backend/internal/user/request"
	userResponse "survey-backend/internal/user/response"
	pkgResponse "survey-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

func RegisterHandler(c *gin.Context) {
	r := &pkgResponse.Response[userResponse.UserResponse]{}

	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, r.Fail(err.Error()))
		return
	}

	newUser := User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  req.Password,
		Role:      req.Role,
		CountryID: req.CountryID,
		RegionID:  req.RegionID,
	}

	err := RegisterUser(&newUser)
	if err != nil {
		// User already exists?
		if strings.Contains(err.Error(), "already exists") {
			c.JSON(http.StatusConflict, r.Fail(err.Error())) // 409 Conflict
			return
		}
		// Other errors
		c.JSON(http.StatusInternalServerError, r.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, r.OK(
		userResponse.UserResponse{
			Name:      req.Name,
			Email:     req.Email,
			Role:      req.Role,
			Points:    0,
			CountryID: req.CountryID,
			RegionID:  req.RegionID},
	))
}

func LoginHandler(c *gin.Context) {

	r := &pkgResponse.Response[userResponse.LoginResponse]{}

	var req request.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("Error Bad Request:", err.Error())
		c.JSON(http.StatusBadRequest, r.Fail(err.Error()))
		return
	}

	user, token, err := LoginUser(req.Email, req.Password)
	if err != nil {
		if err.Error() == "invalid credentials" {
			c.JSON(http.StatusUnauthorized, r.Fail(err.Error())) // 401
			return
		}
		c.JSON(http.StatusInternalServerError, r.Fail(err.Error()))
		return
	}

	data := userResponse.LoginResponse{
		Token:        token,
		UserResponse: user.ToResponse(),
	}
	c.JSON(http.StatusOK, r.OK(data))

}

func GetUserHandler(c *gin.Context) {
	r := &pkgResponse.Response[userResponse.UserResponse]{}
	currentUser := c.MustGet("currentUser").(User)

	c.JSON(http.StatusOK,
		r.OK(currentUser.ToResponse()))
}

func GetRegisterFormHandler(c *gin.Context) {
	r := &pkgResponse.Response[[]userResponse.CountryWithRegions]{}

	response, err := GetRegisterForm()
	if err != nil {
		c.JSON(http.StatusInternalServerError, r.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, r.OK(response))
}

func GrantPointsHandler(c *gin.Context) {
	r := &pkgResponse.Response[userResponse.UserResponse]{}

	var req request.GrantPointsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, r.Fail(err.Error()))
		return
	}

	updated, err := GrantPoints(req.Email, req.Points)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			c.JSON(http.StatusNotFound, r.Fail(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, r.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, r.OK(updated.ToResponse()))
}
