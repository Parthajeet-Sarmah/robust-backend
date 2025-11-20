package utils

import (
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io"
	custom_types "local/bomboclat-oauth-server/types"
	"math/big"
	"net/http"
	"os"
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

func ConstructJWKSFromPublicKey(publicKey string) map[string][]custom_types.JWK {
	// NOTE: Decode PEM to the original base64 encoding
	block, _ := pem.Decode([]byte(publicKey))

	if block == nil || block.Type != "PUBLIC KEY" {
		return nil
	}

	// NOTE: Convert the key to into crypto.PublicKey
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil
	}

	// NOTE: Cast key to rsa.PublicKey
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil
	}

	// NOTE: Get the modulus (N) and encode it to base64
	n := base64.RawURLEncoding.EncodeToString(rsaPub.N.Bytes())

	// NOTE: Get the exponent (E), convert it to big-endian order, and encode it to base64
	eBytes := make([]byte, 0)
	for e := rsaPub.E; e > 0; e >>= 8 {
		eBytes = append([]byte{byte(e & 0xff)}, eBytes...)
	}
	e := base64.URLEncoding.EncodeToString(eBytes)

	fjwk := custom_types.JWK{
		Kty: "RSA",
		Kid: os.Getenv("JWK_KEY_ID"),
		N:   n,
		E:   e,
		Alg: "RSA256",
		Use: "sig",
	}

	jwks := map[string][]custom_types.JWK{
		"keys": []custom_types.JWK{fjwk},
	}

	return jwks
}

type JWKSet struct {
	Keys []custom_types.JWK `json:"keys"`
}

func BuildRSAPublicKey(jwkJSON []byte, kid string) (*rsa.PublicKey, error) {
	var set JWKSet
	if err := json.Unmarshal(jwkJSON, &set); err != nil {
		return nil, err
	}

	var key custom_types.JWK
	found := false
	for _, k := range set.Keys {
		if k.Kid == kid {
			key = k
			found = true
			break
		}
	}
	if !found {
		return nil, errors.New("kid not found in JWK set")
	}

	if key.Kty != "RSA" {
		return nil, errors.New("unsupported key type")
	}

	nb, err := base64.RawURLEncoding.DecodeString(key.N)
	if err != nil {
		return nil, err
	}

	eb, err := base64.RawURLEncoding.DecodeString(key.E)
	if err != nil {
		return nil, err
	}

	n := new(big.Int).SetBytes(nb)

	e := new(big.Int).SetBytes(eb).Int64()
	if e <= 0 {
		return nil, errors.New("invalid RSA exponent")
	}

	pub := &rsa.PublicKey{
		N: n,
		E: int(e),
	}

	return pub, nil
}
