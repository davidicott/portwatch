package config

import "testing"

func TestHipChatConfig_Defaults(t *testing.T) {
	d := hipChatDefaults()
	if d.Enabled {
		t.Error("expected enabled=false by default")
	}
	if d.ServerURL != "https://api.hipchat.com" {
		t.Errorf("unexpected default server_url: %s", d.ServerURL)
	}
	if d.RoomID != "" {
		t.Errorf("expected empty room_id, got %s", d.RoomID)
	}
	if d.Token != "" {
		t.Errorf("expected empty token, got %s", d.Token)
	}
}

func TestHipChatConfig_Fields(t *testing.T) {
	c := HipChatConfig{
		Enabled:   true,
		ServerURL: "https://chat.example.com",
		RoomID:    "99",
		Token:     "secret",
	}
	if !c.Enabled {
		t.Error("expected enabled")
	}
	if c.RoomID != "99" {
		t.Errorf("unexpected room_id: %s", c.RoomID)
	}
	if c.Token != "secret" {
		t.Errorf("unexpected token: %s", c.Token)
	}
}

func TestHipChatConfig_DefaultsNotMutated(t *testing.T) {
	d1 := hipChatDefaults()
	d2 := hipChatDefaults()
	d1.ServerURL = "https://modified.example.com"
	d1.Token = "mutated"
	if d2.ServerURL != "https://api.hipchat.com" {
		t.Errorf("defaults mutated: server_url changed to %s", d2.ServerURL)
	}
	if d2.Token != "" {
		t.Errorf("defaults mutated: token changed to %s", d2.Token)
	}
}
