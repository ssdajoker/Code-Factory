package github

import (
	"testing"
)

func TestNewOAuthFlow(t *testing.T) {
	flow := NewOAuthFlow("test-client-id")
	if flow == nil {
		t.Fatal("NewOAuthFlow returned nil")
	}
	if flow.config.ClientID != "test-client-id" {
		t.Errorf("ClientID = %v, want 'test-client-id'", flow.config.ClientID)
	}
}

func TestOAuthFlowNoClientID(t *testing.T) {
	flow := NewOAuthFlow("")
	_, err := flow.InitiateDeviceFlow()
	if err == nil {
		t.Error("expected error when client ID is empty")
	}
}

func TestDeviceCodeStruct(t *testing.T) {
	dc := &DeviceCode{
		DeviceCode:      "abc123",
		UserCode:        "ABCD-1234",
		VerificationURI: "https://github.com/login/device",
		ExpiresIn:       900,
		Interval:        5,
	}

	if dc.DeviceCode != "abc123" {
		t.Errorf("DeviceCode = %v, want 'abc123'", dc.DeviceCode)
	}
	if dc.UserCode != "ABCD-1234" {
		t.Errorf("UserCode = %v, want 'ABCD-1234'", dc.UserCode)
	}
	if dc.Interval < 5 {
		t.Errorf("Interval = %v, should be >= 5", dc.Interval)
	}
}

func TestOAuthScopes(t *testing.T) {
	flow := NewOAuthFlow("test")
	scopes := flow.config.Scopes

	expected := []string{"repo", "read:user"}
	if len(scopes) != len(expected) {
		t.Errorf("expected %d scopes, got %d", len(expected), len(scopes))
	}

	for i, s := range expected {
		if scopes[i] != s {
			t.Errorf("scope[%d] = %v, want %v", i, scopes[i], s)
		}
	}
}
