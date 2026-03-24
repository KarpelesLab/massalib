package massalib_test

import (
	"testing"

	"github.com/KarpelesLab/massalib"
)

func TestAddress(t *testing.T) {
	a, err := massalib.DecodeAddress("AU128qq86hv2NzXqhowRzaeMruThQyxQQC3PgW3cgHg2ttgXMTa1A")
	if err != nil {
		t.Errorf("failed to parse addr: %s", err)
		return
	}

	if a.String() != "AU128qq86hv2NzXqhowRzaeMruThQyxQQC3PgW3cgHg2ttgXMTa1A" {
		t.Errorf("failed to reformat address: %s != AU128qq86hv2NzXqhowRzaeMruThQyxQQC3PgW3cgHg2ttgXMTa1A", a.String())
	}

	a, err = massalib.DecodeAddress("P12TZEhzNX7rNGZ271jqCLThfwFhvQQuDpUM6n8uuNiupqsiUaCs")
	if err != nil {
		t.Errorf("failed to parse addr from pubkey: %s", err)
		return
	}

	if a.String() != "AU1zyQ2XEA6ZCCtR3CDQVRC7Q1rBpNZDmQe1EN5FXrQCq9435ziH" {
		t.Errorf("failed to reformat address: %s != %s", a.String(), "AU1zyQ2XEA6ZCCtR3CDQVRC7Q1rBpNZDmQe1EN5FXrQCq9435ziH")
	}

	//log.Printf("addr = %+v", a)
}
