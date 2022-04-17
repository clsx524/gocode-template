package uuid

import "testing"

func TestGenerateUuid(t *testing.T) {
	uuid1 := Generate()
	uuid2 := Generate()
	if uuid1 == uuid2 {
		t.Fail()
	}
}
