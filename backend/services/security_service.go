package services

import (
	"devtools/backend/types"
	"golang.org/x/crypto/bcrypt"
	"sync"
)

type securityService struct {
	strength int
}

var securities *securityService
var onceSecurities sync.Once

func Securities() *securityService {
	if securities == nil {
		onceSecurities.Do(func() {
			securities = &securityService{
				strength: 10,
			}
		})
	}
	return securities
}

func (s *securityService) GenerateBCPassword(plainPassword string) (resp types.JSResp) {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), s.strength)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	resp.Success = true
	resp.Data = string(passwordBytes)
	return
}
