package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/zhack01/goapi/conn"
	"github.com/zhack01/goapi/model"
)

func RequiredAuth(ctx *gin.Context) {
	//tokenString, err := ctx.Cookie("Authorization")
	tokensval := ctx.GetHeader("Authorization")
	tokenString := strings.Split(tokensval, " ")[1]
	fmt.Println(tokenString)
	// if err != nil {
	// 	ctx.AbortWithStatus(http.StatusUnauthorized)
	// }

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECERET_KEY")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized, model.Status{
					Code:    http.StatusUnauthorized,
					Message: "Token already expired!",
				},
			)
		}
		var user model.User
		conn.DB.Select([]string{"id", "email", "password", "username"}).Where("id = ?", claims["sub"]).Find(&user)
		if user.ID == 0 {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized, model.Status{
					Code:    http.StatusUnauthorized,
					Message: "Token not found!",
				},
			)
		}
		// Get the requested endpoint path
		endpoint := ctx.FullPath()

		// Get user's IP address
		clientIP := ctx.ClientIP()

		// Get the domain from the request
		domain := ctx.Request.Host

		// i want to save the endpoint, userId, ipaddress and which domain does request came from
		// Save this information or perform any necessary logging
		// For example, you can save this information to a log file or database
		SaveRequestInformation(endpoint, user.ID, clientIP, domain)

		// Check if the endpoint and userID combination exists in the database
		status, restrictionStatus := CheckPageAccess(endpoint, user.ID)
		if status == http.StatusUnauthorized {
			// If the combination exists, restrict access to the page
			ctx.AbortWithStatusJSON(status, restrictionStatus)
			return
		}

		ctx.Next()
	} else {
		ctx.AbortWithStatusJSON(
			http.StatusUnauthorized, model.Status{
				Code:    http.StatusUnauthorized,
				Message: "Invalid token!",
			},
		)
	}

}

func CheckPageAccess(endpoint string, userID uint) (int, model.Status) {
	// Remove the "/api" prefix from the endpoint
	endpoint = strings.TrimPrefix(endpoint, "/api/")

	var access string
	// Execute the updated query to check endpoint and permission type based on user ID
	err := conn.DB.Raw(`
		SELECT 
			(CASE WHEN rp.permission_type = 1000 AND apa.permission_type = 1001 THEN 'false' ELSE 'true' END) AS access
		FROM pages p
		LEFT JOIN api_page_access apa ON p.page_id = apa.page_id
		LEFT JOIN endpoints e ON apa.api_id = e.api_id
		LEFT JOIN role_permissions rp ON rp.page_id = p.page_id
		LEFT JOIN roles r ON rp.role_id = r.role_id
		LEFT JOIN user_profile up ON up.role_id = r.role_id
		WHERE up.user_id = ? AND e.path = ?
		GROUP BY p.page_id, e.path, apa.permission_type, rp.permission_type`, userID, endpoint).Row().Scan(&access)

	if err != nil {
		fmt.Println("endpoint err: ", err)
		// Return a message indicating restrictions for the page
		return http.StatusUnauthorized, model.Status{
			Code:    http.StatusUnauthorized,
			Message: "Access denied.",
		}
	}

	// Compare permissionType with result.permission_type
	if access == "false" {
		return http.StatusUnauthorized, model.Status{
			Code:    http.StatusUnauthorized,
			Message: "Access denied.",
		}
	}

	// If the query found a match for the endpoint and permission type, return an empty status and message
	return 0, model.Status{}
}

func SaveRequestInformation(endpoint string, userID uint, clientIP string, domain string) {
	// Perform the logic to save this information to a database or log file
	// For example:
	// Create a struct to hold this information
	// requestInfo := model.ActivityLogs{
	// 	Endpoint:  endpoint,
	// 	UserID:    userID,
	// 	ClientIP:  clientIP,
	// 	Domain:    domain,
	// 	CreatedAt: time.Now(),
	// }

	// Start a transaction
	// tx := conn.DB.Begin()

	// Save this information to the database using your ORM or log it to a file
	// For database:
	// conn.DB.Create(&requestInfo)
	// For logging to a file:
	// log.Printf("Endpoint: %s, UserID: %d, IP: %s, Domain: %s\n", endpoint, userID, clientIP, domain)

	// Commit the transaction if everything succeeds
	// tx.Commit()
}
