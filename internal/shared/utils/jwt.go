package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/topboyasante/pitstop/internal/core/config"
	"github.com/topboyasante/pitstop/internal/core/logger"
)

func CreateJWTTokens(config *config.Config, userID string, audience string) (string, string, int64, error) {
	logger.Debug("Creating JWT tokens", "userID", userID, "audience", audience)

	accessTokenExp := time.Now().Add(time.Minute * 30).Unix()
	accessTokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,                  // Subject (user identifier)
		"iss": config.Server.JWTIssuer, // Issuer
		"aud": audience,                // Audience (intended recipient)
		"exp": accessTokenExp,          // Expiration time
		"iat": time.Now().Unix(),       // Issued at
	})

	accessTokenString, err := accessTokenClaims.SignedString([]byte(config.Server.JWTSecret))
	if err != nil {
		logger.Error("Failed to sign access token", "error", err, "userID", userID)
		return "", "", 0, err
	}

	refreshTokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,                                     // Subject (user identifier)
		"iss": config.Server.JWTIssuer,                    // Issuer
		"aud": audience,                                   // Audience (intended recipient)
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), // Expiration time = 30 days
		"iat": time.Now().Unix(),                          // Issued at
	})

	refreshTokenString, err := refreshTokenClaims.SignedString([]byte(config.Server.JWTSecret))
	if err != nil {
		logger.Error("Failed to sign refresh token", "error", err, "userID", userID)
		return "", "", 0, err
	}

	logger.Info("JWT tokens created successfully", "userID", userID, "audience", audience, "accessTokenExp", accessTokenExp)
	return accessTokenString, refreshTokenString, accessTokenExp, nil
}

func ValidateJWTToken(config *config.Config, tokenString string) (*jwt.Token, error) {
	logger.Debug("Validating JWT token")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.Error("Unexpected signing method", "method", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Server.JWTSecret), nil
	})

	if err != nil {
		logger.Error("Failed to parse JWT token", "error", err)
		return nil, err
	}

	if !token.Valid {
		logger.Warn("Invalid JWT token provided")
		return nil, errors.New("invalid token")
	}

	logger.Debug("JWT token validated successfully")
	return token, nil
}

func ExtractClaims(token *jwt.Token) (userID string, audience string, exp int64, err error) {
	logger.Debug("Extracting claims from JWT token")

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		logger.Error("Invalid token claims format")
		return "", "", 0, errors.New("invalid token claims")
	}

	subClaim, ok := claims["sub"].(string)
	if !ok {
		logger.Error("Invalid or missing subject claim")
		return "", "", 0, errors.New("invalid subject claim")
	}

	audClaim, ok := claims["aud"].(string)
	if !ok {
		logger.Error("Invalid or missing audience claim")
		return "", "", 0, errors.New("invalid audience claim")
	}

	expClaim, ok := claims["exp"].(float64)
	if !ok {
		logger.Error("Invalid or missing expiration claim")
		return "", "", 0, errors.New("invalid expiration claim")
	}

	logger.Debug("Claims extracted successfully", "userID", subClaim, "audience", audClaim)
	return subClaim, audClaim, int64(expClaim), nil
}

func RefreshToken(config *config.Config, refreshTokenString string) (string, string, int64, error) {
	logger.Info("Refreshing JWT tokens")

	token, err := ValidateJWTToken(config, refreshTokenString)
	if err != nil {
		logger.Error("Invalid refresh token provided", "error", err)
		return "", "", 0, fmt.Errorf("invalid refresh token: %w", err)
	}

	userID, audience, _, err := ExtractClaims(token)
	if err != nil {
		logger.Error("Failed to extract claims from refresh token", "error", err)
		return "", "", 0, fmt.Errorf("failed to extract claims: %w", err)
	}

	logger.Info("Refresh token validated, creating new tokens", "userID", userID)
	return CreateJWTTokens(config, userID, audience)
}

func GetUserIDFromToken(config *config.Config, tokenString string) (string, error) {
	logger.Debug("Extracting user ID from JWT token")

	token, err := ValidateJWTToken(config, tokenString)
	if err != nil {
		logger.Error("Failed to validate token for user ID extraction", "error", err)
		return "", err
	}

	userID, _, _, err := ExtractClaims(token)
	if err != nil {
		logger.Error("Failed to extract user ID from token", "error", err)
		return "", err
	}

	logger.Debug("User ID extracted successfully", "userID", userID)
	return userID, nil
}
