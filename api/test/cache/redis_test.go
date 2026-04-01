package cache_test

import (
	"context"
	"testing"
	"time"

	"main/cache"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func newTestClient(t *testing.T) (*redis.Client, *miniredis.Miniredis) {
	t.Helper()
	mr := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	return rdb, mr
}

// GetJSON / SetJSON — nil client (caching disabled)

func TestGetJSON_NilClient(t *testing.T) {
	_, ok := cache.GetJSON[string](context.Background(), nil, "key")
	if ok {
		t.Fatal("expected ok=false with nil client")
	}
}

func TestSetJSON_NilClient(t *testing.T) {
	if err := cache.SetJSON(context.Background(), nil, "key", "value", time.Minute); err != nil {
		t.Fatalf("expected no error with nil client, got %v", err)
	}
}

// GetJSON / SetJSON — round-trip

func TestSetJSON_GetJSON_String(t *testing.T) {
	rdb, _ := newTestClient(t)
	ctx := context.Background()

	if err := cache.SetJSON(ctx, rdb, "str", "hello", time.Minute); err != nil {
		t.Fatalf("SetJSON: %v", err)
	}

	got, ok := cache.GetJSON[string](ctx, rdb, "str")
	if !ok {
		t.Fatal("expected ok=true")
	}
	if got != "hello" {
		t.Fatalf("got %q, want %q", got, "hello")
	}
}

func TestSetJSON_GetJSON_Struct(t *testing.T) {
	type payload struct {
		UserID string `json:"user_id"`
		Score  int    `json:"score"`
	}

	rdb, _ := newTestClient(t)
	ctx := context.Background()

	want := payload{UserID: "123", Score: 42}
	if err := cache.SetJSON(ctx, rdb, "score:1:123", want, time.Minute); err != nil {
		t.Fatalf("SetJSON: %v", err)
	}

	got, ok := cache.GetJSON[payload](ctx, rdb, "score:1:123")
	if !ok {
		t.Fatal("expected ok=true")
	}
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestGetJSON_MissingKey(t *testing.T) {
	rdb, _ := newTestClient(t)
	_, ok := cache.GetJSON[string](context.Background(), rdb, "nonexistent")
	if ok {
		t.Fatal("expected ok=false for missing key")
	}
}

func TestGetJSON_InvalidJSON(t *testing.T) {
	rdb, mr := newTestClient(t)
	mr.Set("bad", "not-json{{{")

	_, ok := cache.GetJSON[map[string]string](context.Background(), rdb, "bad")
	if ok {
		t.Fatal("expected ok=false for invalid JSON")
	}
}

// TTL expiry

func TestSetJSON_Expiry(t *testing.T) {
	rdb, mr := newTestClient(t)
	ctx := context.Background()

	if err := cache.SetJSON(ctx, rdb, "expiring", "value", 1*time.Second); err != nil {
		t.Fatalf("SetJSON: %v", err)
	}

	mr.FastForward(2 * time.Second)

	_, ok := cache.GetJSON[string](ctx, rdb, "expiring")
	if ok {
		t.Fatal("expected key to have expired")
	}
}

// Delete

func TestDelete_NilClient(t *testing.T) {
	// Should not panic
	cache.Delete(context.Background(), nil, "key1", "key2")
}

func TestDelete_RemovesKeys(t *testing.T) {
	rdb, _ := newTestClient(t)
	ctx := context.Background()

	_ = cache.SetJSON(ctx, rdb, "k1", "v1", time.Minute)
	_ = cache.SetJSON(ctx, rdb, "k2", "v2", time.Minute)

	cache.Delete(ctx, rdb, "k1", "k2")

	if _, ok := cache.GetJSON[string](ctx, rdb, "k1"); ok {
		t.Fatal("k1 should have been deleted")
	}
	if _, ok := cache.GetJSON[string](ctx, rdb, "k2"); ok {
		t.Fatal("k2 should have been deleted")
	}
}

func TestDelete_NonexistentKey(t *testing.T) {
	rdb, _ := newTestClient(t)
	// Should not error or panic
	cache.Delete(context.Background(), rdb, "ghost")
}
