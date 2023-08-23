package adsig

import (
	"git.bode.fun/adsig/config"
	"git.bode.fun/adsig/internal/util"
	"github.com/go-ldap/ldap/v3"
)

type Group struct {
	Name      string
	Templates []Template
	Members   []*ldap.Entry
}

func (g Group) ContainsEmail(email string) bool {
	email = util.NormalizeEmail(email)

	for _, member := range g.Members {
		memberEmail := util.NormalizeEmail(member.GetAttributeValue("mail"))

		if memberEmail == email {
			return true
		}
	}

	return false
}

func GroupsFromConfig(cnf config.Config, conn *ldap.Conn) ([]Group, error) {
	templates, err := templatesFromConfig(cnf)
	if err != nil {
		return nil, err
	}

	groups := make([]Group, 0)

	for cnfGroupName, cnfGroup := range cnf.Groups {
		group := Group{
			Name:      cnfGroupName,
			Templates: make([]Template, 0),
			Members:   make([]*ldap.Entry, 0),
		}

		entries, err := searchMembersForGroup(conn, cnfGroup.BaseDN, cnfGroup.AdFilter)
		if err != nil {
			return nil, err
		}

		group.Members = filterMembersByEmailDenylist(entries, cnfGroup.ExcludeEmails)

		group.Templates = filterTemplatesByName(templates, cnfGroup.Templates)

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
