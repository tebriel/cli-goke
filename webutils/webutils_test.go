package webutils

import (
	"golang.org/x/net/html"
	"testing"
)

func TestGetAttr(t *testing.T) {
	attrs := []html.Attribute{html.Attribute{Key: "a", Val: "abc123"}, html.Attribute{Key: "div", Val: "321bca"}}
	attr, err := GetAttr("div", attrs)
	if attr != "321bca" || err != nil {
		t.Error()
	}
}
