package configs

import _ "embed"

// x25519KeyPairsBytes is the x25519 key pairs yaml file.
//
//go:embed x25519_key_pairs.yaml
var x25519KeyPairsBytes []byte
