// Package clients membungkus koneksi gRPC ke service lain (DRY: pakai
// service-utils/grpcutil + stub gen di service-utils).
package clients

import (
	"context"

	"github.com/vikikurnia87/service-order/configs"

	userv1 "github.com/vikikurnia87/service-utils/gen/go/user/v1"
	"github.com/vikikurnia87/service-utils/grpcutil"
	"google.golang.org/grpc"
)

// UserClient adalah client gRPC ke service-user (auth + detail user).
type UserClient struct {
	conn *grpc.ClientConn
	cli  userv1.UserServiceClient
}

// NewUserClient membuka koneksi gRPC (lazy) ke service-user.
func NewUserClient(target string) (*UserClient, error) {
	conn, err := grpcutil.Dial(grpcutil.ClientConfig{
		Target:    target,
		Insecure:  true, // jaringan internal; ganti TLS untuk produksi
		EnableAPM: configs.ElasticAPMServerURL != "",
	})
	if err != nil {
		return nil, err
	}
	return &UserClient{conn: conn, cli: userv1.NewUserServiceClient(conn)}, nil
}

// ValidateToken memvalidasi JWT ke service-user dan mengembalikan konteks user+tenant.
func (c *UserClient) ValidateToken(ctx context.Context, accessToken string) (*userv1.ValidateTokenResponse, error) {
	return c.cli.ValidateToken(ctx, &userv1.ValidateTokenRequest{AccessToken: accessToken})
}

// GetUser mengambil detail user by id dari service-user.
func (c *UserClient) GetUser(ctx context.Context, id int64) (*userv1.UserResponse, error) {
	return c.cli.GetUser(ctx, &userv1.GetUserRequest{Id: id})
}

// GetUsersByUuids mengambil detail ringkas banyak user by user_uuid sekaligus,
// dikembalikan sebagai map[user_uuid]UserBrief untuk lookup O(1) saat enrich.
// Input kosong → map kosong tanpa round-trip.
func (c *UserClient) GetUsersByUuids(ctx context.Context, uuids []string) (map[string]*userv1.UserBrief, error) {
	if len(uuids) == 0 {
		return map[string]*userv1.UserBrief{}, nil
	}
	resp, err := c.cli.GetUsersByUuids(ctx, &userv1.GetUsersByUuidsRequest{UserUuids: uuids})
	if err != nil {
		return nil, err
	}
	out := make(map[string]*userv1.UserBrief, len(resp.GetUsers()))
	for _, u := range resp.GetUsers() {
		out[u.GetUserUuid()] = u
	}
	return out, nil
}

// Close menutup koneksi gRPC.
func (c *UserClient) Close() error {
	return c.conn.Close()
}
