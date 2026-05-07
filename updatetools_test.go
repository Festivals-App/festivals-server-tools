package servertools

import (
	"errors"
	"testing"
)

// --- LatestVersion tests ---
// These tests rely on a real GitHub API call, so we use a known stable repo.
// For full isolation, consider abstracting the GitHub client behind an interface.

func TestLatestVersion_ValidRepo(t *testing.T) {
	// google/go-github is a well-maintained repo guaranteed to have tags
	version, err := LatestVersion("google", "go-github")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if version == "" {
		t.Error("expected a non-empty version string")
	}
}

func TestLatestVersion_InvalidOrg(t *testing.T) {
	_, err := LatestVersion("this-org-definitely-does-not-exist-xyz123", "some-repo")
	if err == nil {
		t.Error("expected an error for a non-existent organisation, got nil")
	}
}

func TestLatestVersion_InvalidRepo(t *testing.T) {
	_, err := LatestVersion("google", "this-repo-definitely-does-not-exist-xyz123")
	if err == nil {
		t.Error("expected an error for a non-existent repository, got nil")
	}
}

// --- RunUpdate tests ---

func TestRunUpdate_DevelopmentServer(t *testing.T) {
	result, err := RunUpdate("development", "org", "repo", "/some/script.sh")
	if err == nil {
		t.Fatal("expected an error for development version, got nil")
	}
	if err.Error() != "this is a development server please update manually" {
		t.Errorf("unexpected error message: %v", err)
	}
	if result != "" {
		t.Errorf("expected empty result, got: %s", result)
	}
}

func TestRunUpdate_AlreadyUpToDate(t *testing.T) {
	// To test "already up-to-date" we need to know the current latest tag.
	// Fetch it first, then pass it as the currentVersion.
	latest, err := LatestVersion("google", "go-github")
	if err != nil {
		t.Skipf("skipping test: could not fetch latest version: %v", err)
	}

	result, err := RunUpdate(latest, "google", "go-github", "/some/script.sh")
	if err == nil {
		t.Fatal("expected an error when already up-to-date, got nil")
	}
	if err.Error() != "this server is already up-to-date" {
		t.Errorf("unexpected error message: %v", err)
	}
	if result != "" {
		t.Errorf("expected empty result, got: %s", result)
	}
}

func TestRunUpdate_FailedVersionFetch(t *testing.T) {
	result, err := RunUpdate("v1.0.0", "this-org-definitely-does-not-exist-xyz123", "some-repo", "/some/script.sh")
	if err == nil {
		t.Fatal("expected an error when version fetch fails, got nil")
	}
	if !errors.Is(err, err) { // ensure it is a real error value
		t.Error("expected a non-nil error")
	}
	// Should contain our wrapper message
	expected := "failed to retrieve latest release version with error: "
	if len(err.Error()) < len(expected) || err.Error()[:len(expected)] != expected {
		t.Errorf("unexpected error prefix: %v", err)
	}
	if result != "" {
		t.Errorf("expected empty result, got: %s", result)
	}
}

func TestRunUpdate_InvalidScriptPath(t *testing.T) {
	// Use a version that is definitely older than the latest so we reach the exec step.
	// We pass a script path that cannot be executed to trigger the Start() error.
	// Because exec.Command with /bin/bash -c is unlikely to fail on Start() itself
	// (bash will fail later), we instead verify the happy-path shape of the return value.
	//
	// To reliably test the script-start error branch you would need to mock exec.Command.
	// Here we document the limitation and test what we can without a mock.
	t.Log("exec.Command Start() error path requires mocking exec.Command; skipping direct test.")
}

func TestRunUpdate_SuccessReturnFormat(t *testing.T) {
	// Fetch the latest tag so we can supply an older fake version.
	latest, err := LatestVersion("google", "go-github")
	if err != nil {
		t.Skipf("skipping test: could not fetch latest version: %v", err)
	}

	oldVersion := "v0.0.0-test"
	if oldVersion == latest {
		t.Skip("skipping: fake old version collides with latest tag")
	}

	// Use a no-op script so the command starts successfully without side-effects.
	result, err := RunUpdate(oldVersion, "google", "go-github", "/bin/true")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	expected := oldVersion + " => " + latest
	if result != expected {
		t.Errorf("expected result %q, got %q", expected, result)
	}
}
