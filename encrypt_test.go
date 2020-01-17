package main

import (
	"testing"
)

func TestEncrypt(t *testing.T) {
	enc := Encrypt("uims.jlu.edu.cn")
	if enc != "77726476706e69737468656265737421e5fe4c8f693a6445300d8db9d6562d" {
		t.Errorf("Wrong encrypt result: %s\n", enc)
	}
}

func TestDecrypt(t *testing.T) {
	dec := Decrypt("77726476706e69737468656265737421e5fe4c8f693a6445300d8db9d6562d")
	if dec != "uims.jlu.edu.cn" {
		t.Errorf("Wrong decrypt result: %s\n", dec)
	}
}
