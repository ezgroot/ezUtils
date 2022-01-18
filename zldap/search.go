package zldap

import (
	"fmt"

	"github.com/go-ldap/ldap"
)

// DefaultSearchPageSize the default page size of search.
var DefaultSearchPageSize uint32

// Query info from AD server.
func Query(c *ldap.Conn, searchbase string, filter string, attributeList []string) (*ldap.SearchResult, error) {
	if c == nil {
		return nil, fmt.Errorf("connection is nil")
	}

	searchRequest := ldap.NewSearchRequest(
		searchbase,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		filter,
		attributeList,
		nil,
	)

	sr, err := c.SearchWithPaging(searchRequest, DefaultSearchPageSize)
	if err != nil {
		return nil, err
	}

	return sr, nil
}
