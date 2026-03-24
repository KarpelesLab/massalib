package massalib_test

import (
	"testing"

	"github.com/KarpelesLab/massalib"
)

func TestPubkeyParse(t *testing.T) {
	// from https://docs.massa.net/docs/learn/operation-format-execution#public-key
	// Example of legal public key: P1t4JZwHhWNLt4xYabCbukyVNxSbhYPdF6wCYuRmDuHD784juxd
	pk := &massalib.PublicKey{}
	err := pk.UnmarshalText([]byte("P1t4JZwHhWNLt4xYabCbukyVNxSbhYPdF6wCYuRmDuHD784juxd"))
	if err != nil {
		t.Errorf("failed to parse key: %s", err)
	}
}
