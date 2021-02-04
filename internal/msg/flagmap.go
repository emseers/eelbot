package msg

import (
	"sync"
	"time"
)

func newFlagMap() (fMap *flagMap) {
	fMap = &flagMap{
		syncMap: &sync.Map{},
	}
	return
}

func (fMap *flagMap) hasChannel(channelID string) (ok bool) {
	_, ok = fMap.syncMap.Load(channelID)
	return
}

func (fMap *flagMap) addChannel(channelID string) {
	fMap.syncMap.Store(channelID, struct{}{})
}

func (fMap *flagMap) addChannelWithTimedReset(channelID string, resetTime time.Duration) {
	fMap.addChannel(channelID)

	go func() {
		time.Sleep(resetTime)
		fMap.deleteChannel(channelID)
	}()
}

func (fMap *flagMap) deleteChannel(channelID string) {
	fMap.syncMap.Delete(channelID)
}
