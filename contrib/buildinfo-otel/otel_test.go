package buildinfootel

import (
	"testing"

	"go.opentelemetry.io/otel/attribute"
)

func TestAttributes_ContainsExpectedKeys(t *testing.T) {
	attrs := Attributes()

	want := map[attribute.Key]bool{
		KeyVersion:   true,
		KeyCommit:    true,
		KeyBranch:    true,
		KeyTime:      true,
		KeyGoVersion: true,
		KeyGOOS:      true,
		KeyGOARCH:    true,
		KeyModified:  true,
	}
	for _, a := range attrs {
		delete(want, a.Key)
	}
	for k := range want {
		t.Errorf("Attributes missing key %q", k)
	}
}

func TestAttributes_ExtraAppendedAndOverrides(t *testing.T) {
	override := attribute.String(KeyVersion, "1.2.3-overridden")
	extra := attribute.String("custom.key", "custom-value")

	attrs := Attributes(override, extra)

	// Both the original build.version and the override should be present;
	// downstream consumers (e.g. resource.New + resource.Merge) decide which wins.
	versionCount := 0
	customFound := false
	for _, a := range attrs {
		if a.Key == KeyVersion {
			versionCount++
		}
		if a.Key == "custom.key" && a.Value.AsString() == "custom-value" {
			customFound = true
		}
	}
	if versionCount != 2 {
		t.Errorf("KeyVersion count: got %d, want 2 (one from build, one from override)", versionCount)
	}
	if !customFound {
		t.Errorf("custom.key not appended")
	}
}

func TestAttributes_TypedValues(t *testing.T) {
	attrs := Attributes()
	for _, a := range attrs {
		switch a.Key {
		case KeyModified:
			if a.Value.Type() != attribute.BOOL {
				t.Errorf("%s: type got %v, want BOOL", a.Key, a.Value.Type())
			}
		default:
			if a.Value.Type() != attribute.STRING {
				t.Errorf("%s: type got %v, want STRING", a.Key, a.Value.Type())
			}
		}
	}
}
