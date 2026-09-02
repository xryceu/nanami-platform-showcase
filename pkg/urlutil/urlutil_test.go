package urlutil

import "testing"

func TestBuild(t *testing.T) {
	got, err := Build("http", "control-plane", "8080")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "http://control-plane:8080" {
		t.Fatalf("expected http://control-plane:8080, got %q", got)
	}

	got, err = Build("https", "cp.example.com", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "https://cp.example.com" {
		t.Fatalf("expected https://cp.example.com, got %q", got)
	}

	got, err = Build("https", "cp.example.com", "443")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "https://cp.example.com" {
		t.Fatalf("expected https://cp.example.com, got %q", got)
	}

	got, err = Build("http", "cp.example.com", "80")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "http://cp.example.com" {
		t.Fatalf("expected http://cp.example.com, got %q", got)
	}

	if _, err := Build("", "cp.example.com", ""); err == nil {
		t.Fatalf("expected error for missing protocol")
	}

	if _, err := Build("https", "", ""); err == nil {
		t.Fatalf("expected error for missing host")
	}
}
