package password

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/scnewma/auditor/pwned"
)

func PwnedCount(client pwned.Client, password string) (int, error) {
	hash := sha1hash(password)
	prefix := prefix(hash)

	pwnedHashes, err := client.HashesInRange(prefix)
	if err != nil {
		return 0, fmt.Errorf("could not determine if password is pwned: %v", err)
	}

	pwnedHash := find(pwnedHashes, prefix, hash)
	if pwnedHash == nil {
		return 0, nil
	}
	return pwnedHash.Count, nil
}

func find(pwnedPasswords []pwned.Hash, prefix, hash string) *pwned.Hash {
	for _, p := range pwnedPasswords {
		if strings.EqualFold(prefix+p.Value, hash) {
			return &p
		}
	}

	return nil
}

func contains(pwnedPasswords []pwned.Hash, prefix, hash string) bool {
	return find(pwnedPasswords, prefix, hash) != nil
}

func prefix(shaPass string) string {
	return shaPass[0:5]
}

func sha1hash(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
