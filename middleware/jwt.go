package middleware

import (
    "strings"
    "net/http"
    "github.com/dgrijalva/jwt-go"  // Pastikan sudah mengimpor jwt-go
    "github.com/labstack/echo/v4"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        // Ambil token dari header Authorization
        token := c.Request().Header.Get("Authorization")
        if token == "" {
            return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Authorization header missing"})
        }
        
        // Hapus "Bearer " jika ada
        token = strings.TrimPrefix(token, "Bearer ")

        // Parsing dan verifikasi token
        claims := &jwt.MapClaims{}
        parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
            return []byte("JmySuperSecretKey12345"), nil // Ganti dengan secret key yang kamu gunakan
        })
        
        if err != nil || !parsedToken.Valid {
            return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid or expired token"})
        }

        // Set claims ke dalam context untuk digunakan di handler selanjutnya
        c.Set("user", claims)
        
        return next(c)
    }
}
