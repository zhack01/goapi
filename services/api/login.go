package api

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/zhack01/goapi/conn"
	"github.com/zhack01/goapi/model"
	"golang.org/x/crypto/bcrypt"
)

func (s *ServiceAPI) LogIn() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request model.LoginRequest
		if err := ctx.ShouldBindJSON(&request); err != nil {
			log.Println(err)
			ctx.JSON(http.StatusNotAcceptable, Response[UserData]{
				Status: Status{
					Code:    http.StatusNotAcceptable,
					Message: "Error in decoding request",
				},
				Data: nil,
			})
			return
		}
		var user model.User
		//email or username as auth
		if request.Email != "" {
			conn.DB.Select([]string{"id", "email", "password", "username"}).Where("email = ?", request.Email).Find(&user)
		} else if request.Username != "" {
			conn.DB.Select([]string{"id", "email", "password", "username"}).Where("username = ?", request.Username).Find(&user)
		}
		if user.ID == 0 {
			ctx.JSON(http.StatusNotAcceptable, Response[UserData]{
				Status: Status{
					Code:    http.StatusNotFound,
					Message: "Invalid Email or Password.",
				},
				Data: nil,
			})
			return
		}

		// Validate password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
			ctx.JSON(http.StatusNotAcceptable, Response[UserData]{
				Status: Status{
					Code:    http.StatusForbidden,
					Message: "Invalid Email or Password.",
				},
				Data: nil,
			})
			return
		}

		// Create and sign token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.ID,
			"exp": time.Now().Add(time.Hour).Unix(),
		})

		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
		if err != nil {
			ctx.JSON(http.StatusNotAcceptable, Response[UserData]{
				Status: Status{
					Code:    http.StatusForbidden,
					Message: "Failed to create Token",
				},
				Data: nil,
			})
			return
		}

		// Fetch user data and role permissions
		var userDatas struct {
			UserID      uint   `json:"user_id"`
			Name        string `json:"name"`
			Username    string `json:"username"`
			Email       string `json:"email"`
			AccountType string `json:"accountType" gorm:"column:accountType"`
			Role        string `json:"role" gorm:"column:role"`
			RoleID      uint   `json:"role_id"`
			OperatorId  int    `json:"operator_id"`
			BrandId     int    `json:"brand_id"`
			AgentId     int    `json:"agent_id"`
			Status      int    `json:"status_id" gorm:"column:status_id"`
		}

		result := conn.DB.Table("users u").
			Select("ifnull(u.operator_id,0) operator_id, u.id as user_id, u.name, u.username, u.email, ut.name as accountType, u.status_id, r.name as `role`, r.role_id, u.brand_id, u.client_id as agent_id").
			Joins("INNER JOIN user_profile up ON u.id = up.user_id").
			Joins("INNER JOIN roles r ON up.role_id = r.role_id").
			Joins("INNER JOIN user_type ut ON r.user_type_id = ut.user_type_id").
			Where("u.id = ?", user.ID).
			Scan(&userDatas)

		if result.Error != nil {
			ctx.JSON(http.StatusNotAcceptable, Response[UserData]{
				Status: Status{
					Code:    http.StatusForbidden,
					Message: "Failed to fetch user data",
				},
				Data: nil,
			})
			return
		}

		// Update the Response struct
		type Response struct {
			Status Status      `json:"status"`
			Data   interface{} `json:"data"`
		}

		// Create the response structure
		response := struct {
			AccountType string `json:"accountType" gorm:"column:accountType"`
			Name        string `json:"name"`
			Username    string `json:"username"`
			Email       string `json:"email"`
			Role        string `json:"role" gorm:"column:role"`
			RoleID      uint   `json:"role_id"`
			UserID      uint   `json:"user_id"`
			OperatorId  int    `json:"operator_id"`
			BrandId     int    `json:"brand_id"`
			AgentId     int    `json:"agent_id"`
			Status      int    `json:"status_id"`
			Token       string `json:"token"`
		}{
			AccountType: userDatas.AccountType,
			Name:        userDatas.Name,
			Username:    userDatas.Username,
			Email:       userDatas.Email,
			Role:        userDatas.Role,
			UserID:      userDatas.UserID,
			OperatorId:  userDatas.OperatorId,
			BrandId:     userDatas.BrandId,
			AgentId:     userDatas.AgentId,
			Status:      userDatas.Status,
			Token:       tokenString,
		}

		// Return the response
		ctx.JSON(http.StatusOK, Response{
			Status: Status{
				Code:    http.StatusOK,
				Message: "Success Login!",
			},
			Data: response, // Use the response struct directly as the Data field
		})

	}
}
