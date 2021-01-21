package mstore

import (
	"io"

	"github.com/graymeta/stow"
)

// Item -
type Item struct {
	stow.Item
	Body io.Reader
}
