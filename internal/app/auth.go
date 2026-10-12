package app

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	db "garago/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

const sessionCookieName = "garago_session"

type authSettings struct {
	secret       []byte
	sessionTTL   time.Duration
	cookieSecure bool
}

type authClaims struct {
	UserID      int32  `json:"user_id"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	AuthVersion int64  `json:"auth_version"`
	jwt.RegisteredClaims
}

type signupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type forgotPasswordRequest struct {
	Email string `json:"email"`
}

type resetPasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

var (
	redisClient *redis.Client
	authConfig  authSettings
)

func initializeAuth(ctx context.Context, cfg config) error {
	options, err := redis.ParseURL(cfg.redisURL)
	if err != nil {
		return fmt.Errorf("invalid REDIS_URL: %w", err)
	}
	redisClient = redis.NewClient(options)
	if err := redisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}
	authConfig = authSettings{
		secret:       []byte(cfg.jwtSecret),
		sessionTTL:   cfg.sessionTTL,
		cookieSecure: cfg.cookieSecure,
	}
	return nil
}

func newSessionID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func createSessionToken(ctx context.Context, user db.User) (string, string, error) {
	sessionID, err := newSessionID()
	if err != nil {
		return "", "", err
	}
	authVersion, err := getAuthVersion(ctx, user.ID)
	if err != nil {
		return "", "", err
	}
	now := time.Now()
	claims := authClaims{
		UserID:      user.ID,
		Email:       user.Email,
		Role:        user.Role,
		AuthVersion: authVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(authConfig.sessionTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        sessionID,
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(authConfig.secret)
	return token, sessionID, err
}

func authVersionKey(userID int32) string { return fmt.Sprintf("auth-version:%d", userID) }

func getAuthVersion(ctx context.Context, userID int32) (int64, error) {
	value, err := redisClient.Get(ctx, authVersionKey(userID)).Result()
	if errors.Is(err, redis.Nil) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	version, err := strconv.ParseInt(value, 10, 64)
	if err != nil || version < 0 {
		return 0, errors.New("invalid authentication version")
	}
	return version, nil
}

func invalidateUserSessions(ctx context.Context, userID int32) error {
	return redisClient.Incr(ctx, authVersionKey(userID)).Err()
}

func setSessionCookie(c *gin.Context, token string) {
	c.SetCookie(sessionCookieName, token, int(authConfig.sessionTTL.Seconds()), "/", "", authConfig.cookieSecure, true)
}

func clearSessionCookie(c *gin.Context) {
	c.SetCookie(sessionCookieName, "", -1, "/", "", authConfig.cookieSecure, true)
}

func publicUser(id int32, name, email, role string) gin.H {
	return gin.H{"id": id, "name": name, "email": email, "role": role}
}

func signup(c *gin.Context, queries *db.Queries) {
	var request signupRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signup data"})
		return
	}
	request.Name = strings.TrimSpace(request.Name)
	request.Email = strings.ToLower(strings.TrimSpace(request.Email))
	if request.Name == "" || !strings.Contains(request.Email, "@") || len(request.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name, a valid email, and a password of at least 8 characters are required"})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to secure password"})
		return
	}
	user, err := queries.CreateUser(c, db.CreateUserParams{Name: request.Name, Email: request.Email, PasswordHash: string(hash), Role: "admin"})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			c.JSON(http.StatusConflict, gin.H{"error": "An account with this email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user": publicUser(user.ID, user.Name, user.Email, user.Role)})
}

func login(c *gin.Context, queries *db.Queries) {
	var request loginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login data"})
		return
	}
	user, err := queries.GetUserByEmail(c, strings.ToLower(strings.TrimSpace(request.Email)))
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}
	token, sessionID, err := createSessionToken(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}
	if err := redisClient.Set(c, "session:"+sessionID, user.ID, authConfig.sessionTTL).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store session"})
		return
	}
	setSessionCookie(c, token)
	c.JSON(http.StatusOK, gin.H{"user": publicUser(user.ID, user.Name, user.Email, user.Role)})
}

func authenticationRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenText, err := c.Cookie(sessionCookieName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			return
		}
		claims := &authClaims{}
		token, err := jwt.ParseWithClaims(tokenText, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return authConfig.secret, nil
		})
		if err != nil || !token.Valid || claims.ID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired session"})
			return
		}
		if err := redisClient.Get(c, "session:"+claims.ID).Err(); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Session is no longer active"})
			return
		}
		currentVersion, err := getAuthVersion(c, claims.UserID)
		if err != nil || claims.AuthVersion != currentVersion {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Session is no longer active"})
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("session_id", claims.ID)
		c.Next()
	}
}

func passwordResetTokenKey(token string) string {
	hash := sha256.Sum256([]byte(token))
	return "password-reset:" + hex.EncodeToString(hash[:])
}

func passwordResetThrottleKey(email string) string {
	hash := sha256.Sum256([]byte(email))
	return "password-reset-throttle:" + hex.EncodeToString(hash[:])
}

func newPasswordResetToken() (string, error) { return newSessionID() }

func forgotPassword(c *gin.Context, queries *db.Queries) {
	var request forgotPasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A valid email address is required"})
		return
	}
	email := strings.ToLower(strings.TrimSpace(request.Email))
	if !strings.Contains(email, "@") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A valid email address is required"})
		return
	}

	throttleKey := passwordResetThrottleKey(email)
	count, err := redisClient.Incr(c, throttleKey).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to process password reset request"})
		return
	}
	if count == 1 && redisClient.Expire(c, throttleKey, 15*time.Minute).Err() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to process password reset request"})
		return
	}
	if count > 3 {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Please wait before requesting another password reset email"})
		return
	}

	user, err := queries.GetUserByEmail(c, email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to process password reset request"})
		return
	}
	if err == nil {
		token, tokenErr := newPasswordResetToken()
		if tokenErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to process password reset request"})
			return
		}
		if tokenErr = redisClient.Set(c, passwordResetTokenKey(token), user.ID, 15*time.Minute).Err(); tokenErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to process password reset request"})
			return
		}
		resetURL, _ := passwordResetURL(token)
		err := passwordMailer.SendPasswordReset(c, user.Email, resetURL)
		if err != nil {
			_ = redisClient.Del(c, passwordResetTokenKey(token)).Err()
			logPasswordResetDeliveryFailure(user.ID)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to process password reset request"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "If an account exists for that email, a password reset link has been sent."})
}

func resetPassword(c *gin.Context, queries *db.Queries) {
	var request resetPasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password reset data"})
		return
	}
	request.Token = strings.TrimSpace(request.Token)
	if request.Token == "" || len(request.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A valid reset link and a password of at least 8 characters are required"})
		return
	}
	storedUserID, err := redisClient.GetDel(c, passwordResetTokenKey(request.Token)).Int64()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This password reset link is invalid or has expired"})
		return
	}
	if storedUserID < 1 || storedUserID > 2147483647 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This password reset link is invalid or has expired"})
		return
	}
	userID := int32(storedUserID)
	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to secure password"})
		return
	}
	if err = queries.UpdateUserPassword(c, db.UpdateUserPasswordParams{PasswordHash: string(hash), ID: userID}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
		return
	}
	if err = invalidateUserSessions(c, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully. Please sign in with your new password."})
}

func currentUser(c *gin.Context, queries *db.Queries) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}
	user, err := queries.GetUserByID(c, userID.(int32))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Account not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": publicUser(user.ID, user.Name, user.Email, user.Role)})
}

func logout(c *gin.Context) {
	if sessionID, exists := c.Get("session_id"); exists {
		_ = redisClient.Del(c, "session:"+sessionID.(string)).Err()
	}
	clearSessionCookie(c)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}
