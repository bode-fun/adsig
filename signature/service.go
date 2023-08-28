package signature

import (
	"git.bode.fun/adsig"
	"git.bode.fun/adsig/config"
	"github.com/dgraph-io/badger/v3"
	"github.com/go-ldap/ldap/v3"
)

type SignatureService struct {
	DB       *badger.DB
	LdapConn *ldap.Conn
	Config   config.Config
}

func (s *SignatureService) RenderSignatureForAccount(account string) ([]signatureForAccount, error) {
	groups, err := adsig.GroupsFromConfig(s.Config, s.LdapConn)
	if err != nil {
		return nil, err
	}

}
