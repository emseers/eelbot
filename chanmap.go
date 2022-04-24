package eelbot

import (
	"sync"
	"time"
)

// A chanMap is a wrapper around a sync.Map to act like a map[string]struct{}.
type chanMap struct {
	m *sync.Map
}

func newChanMap() *chanMap {
	return &chanMap{new(sync.Map)}
}

func (c *chanMap) hasChannel(channelID string) (ok bool) {
	_, ok = c.m.Load(channelID)
	return
}

func (c *chanMap) addChannel(channelID string, resetTime time.Duration) {
	c.m.Store(channelID, struct{}{})
	if resetTime > 0 {
		go func() {
			time.Sleep(resetTime)
			c.deleteChannel(channelID)
		}()
	}
}

func (c *chanMap) deleteChannel(channelID string) {
	c.m.Delete(channelID)
}
