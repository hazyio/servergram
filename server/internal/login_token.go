package internal

import (
	"sync"
	"time"
)

type LoginTokenInformation struct {
	Expires time.Time
}
type LoginToken struct {
	mu     sync.RWMutex
	tokens map[string]LoginTokenInformation
}

func NewLoginToken() *LoginToken {
	return &LoginToken{tokens: make(map[string]LoginTokenInformation)}
}

func (tokens *LoginToken) Add(token string, expires time.Time) {
	tokens.mu.Lock()
	defer tokens.mu.Unlock()
	for existing, info := range tokens.tokens {
		if info.Expires.Before(time.Now()) {
			delete(tokens.tokens, existing)
		}
	}
	tokens.tokens[token] = LoginTokenInformation{Expires: expires}
}

func (tokens *LoginToken) Valid(token string) bool {
	tokens.mu.RLock()
	info, ok := tokens.tokens[token]
	tokens.mu.RUnlock()
	return ok && info.Expires.After(time.Now())
}
