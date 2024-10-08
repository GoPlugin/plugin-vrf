package vrf

import (
	"github.com/goplugin/plugin-vrf/internal/dkg"
	dkg_contract "github.com/goplugin/plugin-vrf/internal/dkg/contract"
)

type KeyProvider interface {
	KeyLookup(p dkg_contract.KeyID) (kd dkg.KeyData)
}
