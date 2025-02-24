package token_adapter

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

// MockTransaction implements Transaction interface for testing
type MockTransaction struct {
	mock.Mock
}

func (m *MockTransaction) GetHashes() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockTransaction) GetDestinationAddresses() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

// MockToken implements Token interface for testing
type MockToken struct {
	mock.Mock
}

func (m *MockToken) BuildTransaction(txInfo *TransactionInfo) (Transaction, error) {
	args := m.Called(txInfo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(Transaction), args.Error(1)
}

// Test cases
func TestRegisterTokenCreator(t *testing.T) {
	// Clear registry before testing
	tokenRegistry = make(map[string]TokenCreator)

	tests := []struct {
		name        string
		tokenID     string
		creator     TokenCreator
		expectError bool
	}{
		{
			name:        "Valid registration",
			tokenID:     "ETH",
			creator:     func(tokenID string) Token { return &MockToken{} },
			expectError: false,
		},
		{
			name:        "Empty tokenID",
			tokenID:     "",
			creator:     func(tokenID string) Token { return &MockToken{} },
			expectError: true,
		},
		{
			name:        "Nil creator",
			tokenID:     "BTC",
			creator:     nil,
			expectError: true,
		},
		{
			name:        "Duplicate registration",
			tokenID:     "ETH",
			creator:     func(tokenID string) Token { return &MockToken{} },
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RegisterTokenCreator(tt.tokenID, tt.creator)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewToken(t *testing.T) {
	// Clear and setup registry
	tokenRegistry = make(map[string]TokenCreator)
	mockCreator := func(tokenID string) Token { return &MockToken{} }
	RegisterTokenCreator("ETH", mockCreator)

	tests := []struct {
		name        string
		tokenID     string
		expectError bool
	}{
		{
			name:        "Valid token",
			tokenID:     "ETH",
			expectError: false,
		},
		{
			name:        "Invalid token",
			tokenID:     "INVALID",
			expectError: true,
		},
		{
			name:        "Empty token",
			tokenID:     "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := NewToken(tt.tokenID)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, token)
			}
		})
	}
}

func TestBuildTransaction(t *testing.T) {
	mockToken := new(MockToken)
	mockTx := new(MockTransaction)

	// Setup test data
	txInfo := &TransactionInfo{}

	// Setup expectations
	mockTx.On("GetHashes").Return([]string{"0xabc"}, nil)
	mockTx.On("GetDestinationAddresses").Return([]string{"0xdef"}, nil)
	mockToken.On("BuildTransaction", txInfo).Return(mockTx, nil)

	// Test BuildTransaction
	tx, err := mockToken.BuildTransaction(txInfo)
	assert.NoError(t, err)
	assert.NotNil(t, tx)

	// Test GetHashes
	hashes, err := tx.GetHashes()
	assert.NoError(t, err)
	assert.Equal(t, []string{"0xabc"}, hashes)

	// Test GetDestinationAddresses
	addresses, err := tx.GetDestinationAddresses()
	assert.NoError(t, err)
	assert.Equal(t, []string{"0xdef"}, addresses)

	// Verify all expectations were met
	mockToken.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestGetSupportedTokenIDs(t *testing.T) {
	// Clear and setup registry
	tokenRegistry = make(map[string]TokenCreator)
	mockCreator := func(tokenID string) Token { return &MockToken{} }

	// Register some tokens
	RegisterTokenCreator("ETH", mockCreator)
	RegisterTokenCreator("BTC", mockCreator)

	// Get supported tokens
	tokens := GetSupportedTokenIDs()

	// Verify results
	assert.Equal(t, 2, len(tokens))
	assert.Contains(t, tokens, "ETH")
	assert.Contains(t, tokens, "BTC")
}

func TestIsTokenIDSupported(t *testing.T) {
	// Clear and setup registry
	tokenRegistry = make(map[string]TokenCreator)
	mockCreator := func(tokenID string) Token { return &MockToken{} }
	RegisterTokenCreator("ETH", mockCreator)

	tests := []struct {
		name     string
		tokenID  string
		expected bool
	}{
		{
			name:     "Supported token",
			tokenID:  "ETH",
			expected: true,
		},
		{
			name:     "Unsupported token",
			tokenID:  "BTC",
			expected: false,
		},
		{
			name:     "Empty token",
			tokenID:  "",
			expected: false,
		},
		{
			name:     "Lowercase token",
			tokenID:  "eth",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsTokenIDSupported(tt.tokenID)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUnregisterTokenCreator(t *testing.T) {
	// Clear and setup registry
	tokenRegistry = make(map[string]TokenCreator)
	mockCreator := func(tokenID string) Token { return &MockToken{} }
	RegisterTokenCreator("ETH", mockCreator)

	// Verify token is registered
	assert.True(t, IsTokenIDSupported("ETH"))

	// Unregister token
	UnregisterTokenCreator("ETH")

	// Verify token is no longer registered
	assert.False(t, IsTokenIDSupported("ETH"))
}
