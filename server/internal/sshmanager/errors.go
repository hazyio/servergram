package sshmanager

import "fmt"

type UnknownHostError struct{ Host, Fingerprint, KeyType string }

func (e *UnknownHostError) Error() string {
	return fmt.Sprintf("unknown host key for %s:%s", e.Host, e.Fingerprint)
}

type HostKeyMismatchError struct{ Host, Expected, Got string }

func (e *HostKeyMismatchError) Error() string {
	return fmt.Sprintf("host key mismatch for %s: expected %s, got %s", e.Host, e.Expected, e.Got)
}

