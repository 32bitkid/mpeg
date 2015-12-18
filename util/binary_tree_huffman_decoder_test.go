package util_test

import t "testing"
import "bytes"
import "github.com/32bitkid/mpeg/util"

func TestBTHD_Simple(t *t.T) {

	// 0xa5 => 0b 1 01 00 1 01 => True, False, Maybe, True, False

	br := util.NewBitReader(bytes.NewReader([]byte{0xa5}))

	init := util.HuffmanTable{
		{"1", "True"},
		{"01", "False"},
		{"00", "Maybe"},
	}

	hd := util.NewBinaryTreeHuffmanDecoder(init)

	actual, err := hd.Decode(br)
	if err != nil {
		t.Error(err)
	} else if expected := "True"; actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}

	actual, err = hd.Decode(br)
	if err != nil {
		t.Error(err)
	} else if expected := "False"; actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}

	actual, err = hd.Decode(br)
	if err != nil {
		t.Error(err)
	} else if expected := "Maybe"; actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}

	actual, err = hd.Decode(br)
	if err != nil {
		t.Error(err)
	} else if expected := "True"; actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}

	actual, err = hd.Decode(br)
	if err != nil {
		t.Error(err)
	} else if expected := "False"; actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}

}
