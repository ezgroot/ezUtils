package zldap

import (
	"fmt"
	"time"

	"github.com/go-ldap/ldap"
)

// Connect to AD or LDAP server, default port is 389.
func Connect(domain string, timeout time.Duration) (*ldap.Conn, error) {
	c, err := ldap.Dial("tcp", domain)
	if err != nil {
		return nil, err
	}

	c.SetTimeout(timeout * time.Millisecond)

	return c, nil
}

// Login AD or LDAP server.
func Login(c *ldap.Conn, account string, password string) (*ldap.Conn, error) {
	if c == nil {
		return nil, fmt.Errorf("connection is nil")
	}

	err := c.Bind(account, password)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// Close connection.
func Close(c *ldap.Conn) {
	if c == nil {
		return
	}

	c.Close()
}

func init() {
	ldap.DefaultTimeout = time.Duration(1500) * time.Millisecond

	DefaultSearchPageSize = 256
}
