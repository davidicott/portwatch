package config

import (
	"testing"
)

func TestWebexConfig_Defaults(t *testing.T) {
	d := webexDefaults()
	if d.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if d.Token != "" {
		t.Errorf("expected empty Token, got %q", d.Token)
	}
	if d.RoomID != "" {
		t.Errorf("expected empty RoomID, got %q", d.RoomID)
	}
}

func TestWebexConfig_Fields(t *testing.T) {
	nc := &NotifierConfig{}
	if fn, ok := notifierDefaults["webex"]; ok {
		fn(nc)
	} else {
		t.Fatal("webex not registered in notifierDefaults")
	}
	if nc.Webex == nil {
		t.Fatal("expected Webex config to be populated")
	}
	nc.Webex.Enabled = true
	nc.Webex.Token = "abc123"
	nc.Webex.RoomID = "room-xyz"

	if !nc.Webex.Enabled {
		t.Error("expected Enabled=true")
	}
	if nc.Webex.Token != "abc123" {
		t.Errorf("expected Token=abc123, got %q", nc.Webex.Token)
	}
	if nc.Webex.RoomID != "room-xyz" {
		t.Errorf("expected RoomID=room-xyz, got %q", nc.Webex.RoomID)
	}
}
