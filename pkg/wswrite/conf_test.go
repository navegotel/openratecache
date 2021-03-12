package wswrite

import (
	"testing"
)

func TestSettings(t *testing.T) {
	CreateInitialSettings("demo.conf")
	s, err := LoadSettings("demo.conf")
	if err != nil {
		t.Error(err)
	}
	if s.Port != 2511 {
		t.Errorf("Expected value is 2511, got value %v", s.Port)
	}
	//os.Remove("demo.conf")

}
