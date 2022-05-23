package ssh

import (
	"fmt"
	"strings"
)

type SSH struct {
	Key  string
	User string
	Host string
}

func (s SSH) String() string {
	var str strings.Builder
	fmt.Fprintf(&str, "Key:  %v\n", s.Key)
	fmt.Fprintf(&str, "User: %v\n", s.User)
	fmt.Fprintf(&str, "Host: %v\n", s.Host)
	return str.String()
}
