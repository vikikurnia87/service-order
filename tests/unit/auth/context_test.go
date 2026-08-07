package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/require"

	"github.com/vikikurnia87/service-order/auth"
	"github.com/vikikurnia87/service-utils/common"
	userv1 "github.com/vikikurnia87/service-utils/gen/go/user/v1"
)

// TestSetContext_PropagatesCompanyUUIDAndUserID adalah guard terhadap drift:
// kalau suatu saat SetContext berubah dan berhenti menyuntik salah satu dari
// dua key context.Context yang dipakai monitoring.ContextLogger lintas
// service, test ini gagal — bukan baru ketahuan lewat log yang diam-diam
// kehilangan dimensi observability.
func TestSetContext_PropagatesCompanyUUIDAndUserID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	resp := &userv1.ValidateTokenResponse{
		Valid:       true,
		UserId:      42,
		CompanyUuid: "11111111-1111-1111-1111-111111111111",
		UserUuid:    "22222222-2222-2222-2222-222222222222",
	}

	auth.SetContext(c, resp)

	ctx := c.Request().Context()

	gotCompanyUUID, ok := ctx.Value(common.ContextCompanyUUID).(string)
	require.True(t, ok, "context.Context harus membawa common.ContextCompanyUUID")
	require.Equal(t, "11111111-1111-1111-1111-111111111111", gotCompanyUUID)

	gotUserID, ok := ctx.Value(common.ContextUserID).(string)
	require.True(t, ok, "context.Context harus membawa common.ContextUserID")
	require.Equal(t, "42", gotUserID)
}
