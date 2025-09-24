package errors_test

import (
	"testing"

	"github.com/0xPolygon/polygon-cli/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeSmcErrorCode(t *testing.T) {
	reason, err := errors.DecodeSmcErrorCode("0x1a070d9a")
	require.NoError(t, err)
	assert.Equal(t, "InitSequencedBatchDoesNotMatchError", reason)

	reason, err = errors.DecodeSmcErrorCode("0x52ad525a")
	require.NoError(t, err)
	assert.Equal(t, "InvalidPessimisticProofError", reason)

	reason, err = errors.DecodeSmcErrorCode("0x9aad315a")
	require.NoError(t, err)
	assert.Equal(t, "0x9aad315a (unknown selector)", reason)

	var emptyInterface interface{}
	_, err = errors.DecodeSmcErrorCode(emptyInterface)
	require.Error(t, err)

	code := "0x3bbd317c"
	reason, err = errors.DecodeSmcErrorCode(code)
	require.NoError(t, err)
	assert.Equal(t, "0x3bbd317c (unknown selector)", reason)
}
