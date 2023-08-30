package configs

import _ "embed"

//go:embed x25519_key_pairs.yaml
var x25519KeyPairsBytes []byte
