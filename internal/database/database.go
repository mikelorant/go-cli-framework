package database

import (
	"fmt"
	"strings"
)

type Database struct {
	User     string
	Password string
	Host     string
	Port     string
}

func (d Database) String() string {
	var str strings.Builder
	fmt.Fprintf(&str, "User:     %v\n", d.User)
	fmt.Fprintf(&str, "Password: %v\n", d.Password)
	fmt.Fprintf(&str, "Host:     %v\n", d.Host)
	fmt.Fprintf(&str, "Port:     %v\n", d.Port)
	return str.String()
}
