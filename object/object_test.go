package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "namaste duniya"}
	hello2 := &String{Value: "namaste duniya"}
	diff1 := &String{Value: "mera naam raj hai"}
	diff2 := &String{Value: "mera naam raj hai"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("strings with different content have same hash keys")
	}
}
