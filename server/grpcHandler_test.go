package server_test

import (
	"testing"

	"github.com/AntanasMaziliauskas/grpc/server/broker"
	"github.com/stretchr/testify/assert"
)

func TestImplements(t *testing.T) {
	assert.Implements(t, (*broker.BrokerService)(nil), &broker.GRPCBroker{})
}
