package ldapsearch

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"strings"

	"ldap-email-exporter/cmd/ldap-email-exporter/config"
)

type LDAPSearch struct {
	ldap    *ldap.Conn
	baseDN  string
	filter  []string
	exclude []string
	Result  []string
}

func New(config configuration.LDAPConfig) (*LDAPSearch, error) {
	// Connect to LDAP
	l, err := ldap.DialURL(config.URL)
	if err != nil {
		return nil, err
	}

	// Bind to LDAP
	err = l.Bind(config.BindDN, config.Secret)
	if err != nil {
		return nil, err
	}

	// Return LDAPSearch
	return &LDAPSearch{
		ldap:    l,
		baseDN:  config.BaseDN,
		filter:  config.Filters,
		exclude: config.Exclude,
	}, err
}

func (l *LDAPSearch) Search() error {
	// Create search request for each filter in l.filter
	for _, filter := range l.filter {
		searchRequest := ldap.NewSearchRequest(
			l.baseDN,
			ldap.ScopeWholeSubtree,
			ldap.NeverDerefAliases,
			0,
			0,
			false,
			filter,
			[]string{"mail"},
			[]ldap.Control{},
		)

		// Search for users with filters by pages
		sr, err := l.ldap.SearchWithPaging(searchRequest, 100)
		//sr, err := l.ldap.Search(searchRequest)

		if err != nil {
			return err
		}

		fmt.Println("#Got ", len(sr.Entries), " entries")
		for _, entry := range sr.Entries {
			l.append(entry.GetAttributeValue("mail"))
		}
	}
	return nil
}

func (l *LDAPSearch) append(value string) {
	// Trim whitespace in value
	value = strings.TrimSpace(value)

	// Check that value is not empty
	if value == "" {
		return
	}

	// Check that value's domain is not in domains exclude list
	for _, domain := range l.exclude {
		if strings.HasSuffix(value, domain) {
			return
		}
	}

	// Append value to result
	l.Result = append(l.Result, value)
}

// Close closes the connection to the LDAP server.
func (l *LDAPSearch) Close() {
	l.ldap.Close()
}

// Print prints the result of the search.
func (l *LDAPSearch) Print() {
	for _, value := range l.Result {
		fmt.Println(value)
	}
}
