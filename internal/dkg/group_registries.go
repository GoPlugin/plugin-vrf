package dkg

import (
	"go.dedis.ch/kyber/v3/sign/anon"

	"github.com/goplugin/plugin-vrf/altbn_128"
	"github.com/goplugin/plugin-vrf/internal/crypto/point_translation"
)

var translatorRegistry = point_translation.TranslatorRegistry

var altBN128Pairing = &altbn_128.PairingSuite{}

var encryptionGroupRegistry = map[string]anon.Suite{
	"AltBN-128 G₁": altBN128Pairing.G1().(anon.Suite),
}
