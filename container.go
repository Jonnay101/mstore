package mstore

import (
	"fmt"
	"strings"

	"github.com/graymeta/stow"
)

// Container is a wrapper for either a bucket or blob
type Container struct {
	stow.Container
}

// GetItemByKey -
func (c *Container) GetItemByKey(key string) (itm *Item, err error) {
	itm = &Item{}

	walkItemsFunc := func(i stow.Item, err error) error {
		if strings.HasPrefix(i.Name(), key) {
			itm.Body, err = i.Open()
			if err != nil {
				return err
			}

		}
		return nil
	}

	if err = stow.Walk(c, stow.NoPrefix, 100, walkItemsFunc); err != nil {
		return nil, err
	}

	if itm.Body == nil {
		return nil, fmt.Errorf("key: %s not found in %s container", key, c.Name())
	}

	return
}
