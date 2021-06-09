package main

import (
	"fmt"
	"net/url"
	"strings"
)

func tagLi(u, name string) string {

	uu, _ := url.Parse(u)
	s := uu.EscapedPath()
	return fmt.Sprintf("<li><a href='/github/repo?u=%s'>%s</a></li>", s, name)
}

func tagUl(lii []string) string {
	return fmt.Sprintf("<ul>%s</ul>", strings.Join(lii, ""))
}
