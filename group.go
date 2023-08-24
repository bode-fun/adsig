package adsig

import (
	"git.bode.fun/adsig/config"
	"git.bode.fun/adsig/internal/util"
	"github.com/go-ldap/ldap/v3"
)

type Group struct {
	Name       string
	Signatures []Signature
	Members    []*ldap.Entry
}

func (g Group) MemberByEmail(email string) (ok bool, member *ldap.Entry) {
	email = util.NormalizeEmail(email)

	for _, member := range g.Members {
		memberEmail := util.NormalizeEmail(member.GetAttributeValue("mail"))

		if memberEmail == email {
			return true, member
		}
	}

	return false, nil
}

func GroupsFromConfig(cnf config.Config, conn *ldap.Conn) ([]Group, error) {
	templates, err := SignaturesFromConfig(cnf)
	if err != nil {
		return nil, err
	}

	groups := make([]Group, 0)

	for cnfGroupName, cnfGroup := range cnf.Groups {
		group := Group{
			Name:       cnfGroupName,
			Signatures: make([]Signature, 0),
			Members:    make([]*ldap.Entry, 0),
		}

		entries, err := searchMembersForGroup(conn, cnfGroup.BaseDN, cnfGroup.AdFilter)
		if err != nil {
			return nil, err
		}

		group.Members = filterMembersByEmailDenylist(entries, cnfGroup.ExcludeEmails)

		group.Signatures = filterSignaturesByName(templates, cnfGroup.Templates)

		groups = append(groups, group)
	}

	return groups, nil
}

func searchMembersForGroup(conn *ldap.Conn, baseDN, filter string) ([]*ldap.Entry, error) {
	searchRequest := new(ldap.SearchRequest)

	searchRequest.BaseDN = baseDN
	searchRequest.Scope = ldap.ScopeWholeSubtree
	searchRequest.Filter = filter

	searchRes, err := conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	return searchRes.Entries, nil
}

func filterMembersByEmailDenylist(members []*ldap.Entry, denylist []string) []*ldap.Entry {
	filteredMembers := make([]*ldap.Entry, 0)

	for _, member := range members {
		memberEmail := util.NormalizeEmail(member.GetAttributeValue("mail"))
		if memberEmail != "" {
			inExcludeList := false

			for _, emailToExclude := range denylist {
				if emailToExclude == memberEmail {
					inExcludeList = true

					break
				}
			}

			if !inExcludeList {
				filteredMembers = append(filteredMembers, member)
			}
		}
	}

	return filteredMembers
}
