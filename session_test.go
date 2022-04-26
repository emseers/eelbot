package eelbot_test

import (
	"testing"

	"github.com/emseers/eelbot"
	"github.com/stretchr/testify/require"
)

func TestSession(t *testing.T) {
	s, err := eelbot.NewSession("token")
	require.NoError(t, err)
	require.NotNil(t, s)
}
