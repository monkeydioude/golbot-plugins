package golmods

import "testing"

func TestHealthcheckIsOK(t *testing.T) {
	if healthCheck() != "OK" {
		t.Fail()
	}
}
