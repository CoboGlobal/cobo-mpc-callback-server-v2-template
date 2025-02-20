package token_adapter

import (
	"fmt"
	"strings"
	"sync"
)

type TokenCreator func(tokenID string) Token

var (
	registryLock sync.RWMutex
	// tokenRegistry save token id and creator
	tokenRegistry = make(map[string]TokenCreator)
)

func RegisterTokenCreator(tokenID string, creator TokenCreator) error {
	if tokenID == "" {
		return fmt.Errorf("token_id cannot be empty")
	}
	if creator == nil {
		return fmt.Errorf("creator function cannot be nil")
	}

	registryLock.Lock()
	defer registryLock.Unlock()

	tokenID = strings.ToUpper(strings.TrimSpace(tokenID))

	if _, exists := tokenRegistry[tokenID]; exists {
		return fmt.Errorf("token creator with token_id %s is already registered", tokenID)
	}

	tokenRegistry[tokenID] = creator
	return nil
}

func NewToken(tokenID string) (Token, error) {
	registryLock.RLock()
	defer registryLock.RUnlock()

	tokenID = strings.ToUpper(strings.TrimSpace(tokenID))

	creator, ok := tokenRegistry[tokenID]
	if !ok {
		return nil, fmt.Errorf("unsupported token_id: %s", tokenID)
	}

	return creator(tokenID), nil
}

func GetSupportedTokenIDs() []string {
	registryLock.RLock()
	defer registryLock.RUnlock()

	tokenIDs := make([]string, 0, len(tokenRegistry))
	for tokenID := range tokenRegistry {
		tokenIDs = append(tokenIDs, tokenID)
	}
	return tokenIDs
}

func IsTokenIDSupported(tokenID string) bool {
	registryLock.RLock()
	defer registryLock.RUnlock()

	tokenID = strings.ToUpper(strings.TrimSpace(tokenID))
	_, ok := tokenRegistry[tokenID]
	return ok
}

func UnregisterTokenCreator(tokenID string) {
	registryLock.Lock()
	defer registryLock.Unlock()

	tokenID = strings.ToUpper(strings.TrimSpace(tokenID))
	delete(tokenRegistry, tokenID)
}
