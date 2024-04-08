package services

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"devtools/backend/types"
	"fmt"
	"sync"
)

type hashService struct {
}

var hashGenerator *hashService
var onceHashGenerator sync.Once

func HashGenerator() *hashService {
	if hashGenerator == nil {
		onceHashGenerator.Do(func() {
			hashGenerator = &hashService{}
		})
	}
	return hashGenerator
}

func (h *hashService) GenerateHash(content string) (resp types.JSResp) {
	bytes := []byte(content)
	md := md5.Sum(bytes)
	sha := sha1.Sum(bytes)
	sha2 := sha256.Sum256(bytes)
	sha5 := sha512.Sum512(bytes)

	resp.Success = true
	resp.Data = map[string]string{
		"md5":    fmt.Sprintf("%x", md),
		"sha1":   fmt.Sprintf("%x", sha),
		"sha256": fmt.Sprintf("%x", sha2),
		"sha512": fmt.Sprintf("%x", sha5),
	}
	return
}
