package topology

import (
	"context"
	"testing"

	"github.com/mongodb/mongo-go-driver/core/address"
	"github.com/mongodb/mongo-go-driver/core/connection"
	"github.com/mongodb/mongo-go-driver/core/description"
	"github.com/mongodb/mongo-go-driver/core/wiremessage"
	"github.com/stretchr/testify/require"
)

type netErr struct {
}

func (n netErr) Error() string {
	return "error"
}

func (n netErr) Timeout() bool {
	return false
}

func (n netErr) Temporary() bool {
	return false
}

type connect struct {
	err *connection.NetworkError
}

func (c connect) WriteWireMessage(ctx context.Context, wm wiremessage.WireMessage) error {
	return *c.err
}
func (c connect) ReadWireMessage(ctx context.Context) (wiremessage.WireMessage, error) {
	return nil, *c.err
}
func (c connect) Close() error {
	return nil
}
func (c connect) Alive() bool {
	return true
}
func (c connect) Expired() bool {
	return false
}
func (c connect) ID() string {
	return ""
}

// Test case for sconn processErr
func TestConnectionProcessErrSpec(t *testing.T) {
	ctx := context.Background()
	s, err := NewServer(address.Address("localhost"))
	require.NoError(t, err)

	desc := s.Description()
	require.Nil(t, desc.LastError)

	s.connectionstate = connected

	innerErr := netErr{}
	connectErr := connection.NetworkError{"blah", innerErr}
	c := connect{&connectErr}
	sc := sconn{c, s, 1}
	err = sc.WriteWireMessage(ctx, nil)
	require.NotNil(t, err)
	desc = s.Description()
	require.NotNil(t, desc.LastError)
	require.Equal(t, desc.Kind, (description.ServerKind)(description.Unknown))
}
