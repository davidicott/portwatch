package ratelimit_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/ratelimit"
)

func TestAllow_UnderLimit(t *testing.T) {
	l := ratelimit.New(3, time.Minute)
	for i := 0; i < 3; i++ {
		if !l.Allow("tcp:22") {
			t.Fatalf("expected allow on call %d", i+1)
		}
	}
}

func TestAllow_ExceedsLimit(t *testing.T) {
	l := ratelimit.New(2, time.Minute)
	l.Allow("tcp:80")
	l.Allow("tcp:80")
	if l.Allow("tcp:80") {
		t.Fatal("expected deny after limit exceeded")
	}
}

func TestAllow_SeparateKeys(t *testing.T) {
	l := ratelimit.New(1, time.Minute)
	if !l.Allow("tcp:22") {
		t.Fatal("expected allow for tcp:22")
	}
	if !l.Allow("tcp:80") {
		t.Fatal("expected allow for tcp:80 (different key)")
	}
}

func TestAllow_WindowExpiry(t *testing.T) {
	l := ratelimit.New(1, 50*time.Millisecond)
	l.Allow("tcp:443")
	time.Sleep(60 * time.Millisecond)
	if !l.Allow("tcp:443") {
		t.Fatal("expected allow after window expired")
	}
}

func TestReset_ClearsBucket(t *testing.T) {
	l := ratelimit.New(1, time.Minute)
	l.Allow("tcp:22")
	l.Reset("tcp:22")
	if !l.Allow("tcp:22") {
		t.Fatal("expected allow after reset")
	}
}
