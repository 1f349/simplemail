package simplemail

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestFromAddress_UnmarshalText(t *testing.T) {
	t.Run("json", func(t *testing.T) {
		var fromAddress FromAddress
		assert.NoError(t, json.Unmarshal([]byte(`"Jane Doe <jane@example.com>"`), &fromAddress))
		assert.Equal(t, "Jane Doe", fromAddress.Address.Name)
		assert.Equal(t, "jane@example.com", fromAddress.Address.Address)
	})
	t.Run("yaml", func(t *testing.T) {
		var fromAddress FromAddress
		assert.NoError(t, yaml.Unmarshal([]byte(`"Jane Doe <jane@example.com>"`), &fromAddress))
		assert.Equal(t, "Jane Doe", fromAddress.Address.Name)
		assert.Equal(t, "jane@example.com", fromAddress.Address.Address)
	})
}
