package crypto

import (
	"crypto/rand"
	"testing"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/venus/pkg/crypto"
	_ "github.com/filecoin-project/venus/pkg/crypto/secp" // to run init()
)

func TestVerify(t *testing.T) {
	randMsg := make([]byte, 32)
	_, _ = rand.Read(randMsg)

	privkey, _ := crypto.Generate(crypto.SigTypeSecp256k1)
	pubkey, _ := crypto.ToPublic(crypto.SigTypeSecp256k1, privkey)
	addr, _ := address.NewSecp256k1Address(pubkey)
	sig, _ := crypto.Sign(randMsg, privkey, crypto.SigTypeSecp256k1)

	if Verify(sig, addr, randMsg) != nil {
		t.Errorf("error in Verify function")
	}
}

func TestSign(t *testing.T) {
	randMsg := make([]byte, 32)
	_, _ = rand.Read(randMsg)

	privkey, _ := crypto.Generate(crypto.SigTypeSecp256k1)
	pubkey, _ := crypto.ToPublic(crypto.SigTypeSecp256k1, privkey)
	addr, _ := address.NewSecp256k1Address(pubkey)
	sig, _ := Sign(privkey, crypto.SigTypeSecp256k1, randMsg)

	if crypto.Verify(sig, addr, randMsg) != nil {
		t.Errorf("error in Sign function")
	}
}
