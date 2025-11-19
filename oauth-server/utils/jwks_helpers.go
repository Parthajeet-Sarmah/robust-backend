package utils

import (
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
)

// Allowed algorithms for private_key_jwt
var allowedAlgs = map[string]bool{
	"RS256": true, "RS384": true, "RS512": true,
	"PS256": true, "PS384": true, "PS512": true,
	"ES256": true, "ES384": true, "ES512": true,
}

func ValidateJwksJSON(jwksJson string) error {
	set, err := jwk.Parse([]byte(jwksJson))
	if err != nil {
		return err
	}

	if set.Len() == 0 {
		return errors.New("jwks contains no keys")
	}

	it := set.Iterate(context.Background())
	for it.Next(context.Background()) {
		key := it.Pair().Value.(jwk.Key)

		//
		// 1) Validate "use"
		//
		if useVal, ok := key.Get("use"); ok {
			use := useVal.(string)
			if use != "sig" {
				return errors.New("unsupported JWK use: must be 'sig'")
			}
		}

		//
		// 2) Extract underlying public key
		//
		var pub any
		if err := key.Raw(&pub); err != nil {
			return errors.New("failed to extract public key: " + err.Error())
		}

		switch p := pub.(type) {

		case *rsa.PublicKey:
			if p.N.BitLen() < 2048 {
				return errors.New("RSA key too small; must be >= 2048 bits")
			}

		case *ecdsa.PublicKey:
			curve := p.Curve.Params().Name
			switch curve {
			case "P-256", "P-384", "P-521":
				// allowed
			default:
				return errors.New("unsupported EC curve: " + curve)
			}

		default:
			return errors.New("unsupported JWK key type: must be RSA or EC")
		}

		//
		// 3) Validate algorithm (if provided)
		//
		alg := key.Algorithm().String()
		if alg != "" && !allowedAlgs[alg] {
			return errors.New("unsupported JWK alg: " + alg)
		}

		//
		// 4) Must have kid
		//
		if key.KeyID() == "" {
			return errors.New("missing kid in JWK")
		}
	}

	return nil
}

// --- Fetch + validate JWKS URI ---
func FetchAndValidateJwks(ctx context.Context, uri string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("jwks_uri returned non-200")
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return err
	}

	return ValidateJwksJSON(string(data))
}
