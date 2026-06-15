// Package auth menyimpan konteks identitas+tenant hasil ValidateToken (gRPC ke
// service-user) di Echo context, plus getter untuk handler.
package auth

import (
	"slices"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	userv1 "github.com/vikikurnia87/service-utils/gen/go/user/v1"
)

type ctxKey string

const (
	keyUserID      ctxKey = "user_id"
	keyUserUUID    ctxKey = "user_uuid"
	keyCompanyUUID ctxKey = "company_uuid"
	keySystem      ctxKey = "system"
	keyPermissions ctxKey = "permissions"
	keyRoles       ctxKey = "roles"
)

// SetContext menaruh hasil ValidateToken ke Echo context.
func SetContext(c *echo.Context, resp *userv1.ValidateTokenResponse) {
	c.Set(string(keyUserID), resp.GetUserId())
	c.Set(string(keyUserUUID), resp.GetUserUuid())
	c.Set(string(keyCompanyUUID), resp.GetCompanyUuid())
	c.Set(string(keySystem), resp.GetSystem())
	c.Set(string(keyPermissions), resp.GetPermissions())
	c.Set(string(keyRoles), resp.GetRoleCodes())
}

// UserID mengembalikan id user (0 bila tidak ada).
func UserID(c *echo.Context) int64 {
	v, _ := c.Get(string(keyUserID)).(int64)
	return v
}

// UserIDPtr mengembalikan pointer id user untuk kolom created_by/updated_by.
func UserIDPtr(c *echo.Context) *int64 {
	if v := UserID(c); v != 0 {
		return &v
	}
	return nil
}

// CompanyUUID mengembalikan company aktif (uuid.Nil bila tidak ada/invalid).
func CompanyUUID(c *echo.Context) uuid.UUID {
	s, _ := c.Get(string(keyCompanyUUID)).(string)
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}
	return id
}

// IsSystem true bila user role tier SYSTEM (lintas tenant).
func IsSystem(c *echo.Context) bool {
	v, _ := c.Get(string(keySystem)).(bool)
	return v
}

// Permissions mengembalikan daftar permission user.
func Permissions(c *echo.Context) []string {
	v, _ := c.Get(string(keyPermissions)).([]string)
	return v
}

// HasPermission cek apakah user punya permission tertentu.
func HasPermission(c *echo.Context, perm string) bool {
	return slices.Contains(Permissions(c), perm)
}

// UserUUID mengembalikan user uuid (uuid.Nil bila tidak ada/invalid).
func UserUUID(c *echo.Context) uuid.UUID {
	s, _ := c.Get(string(keyUserUUID)).(string)
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}
	return id
}

// UserUUIDPtr mengembalikan pointer user_uuid untuk kolom created_by/updated_by.
func UserUUIDPtr(c *echo.Context) *uuid.UUID {
	id := UserUUID(c)
	if id != uuid.Nil {
		return &id
	}
	return nil
}

// UserIDString untuk konteks log/APM.
func UserIDString(c *echo.Context) string {
	return strconv.FormatInt(UserID(c), 10)
}

// BearerToken mengambil raw JWT dari header Authorization (tanpa prefix "Bearer ").
// Dipakai untuk meneruskan token user ke service lain (mis. service-media passport).
func BearerToken(c *echo.Context) string {
	h := c.Request().Header.Get("Authorization")
	const p = "Bearer "
	if len(h) <= len(p) || !strings.HasPrefix(h, p) {
		return ""
	}
	return strings.TrimSpace(h[len(p):])
}
