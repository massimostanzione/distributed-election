// ref.: https://github.com/golang/go/issues/12721
package misc

import (
	"time"
)

func TimerFictitiousInit(timer *time.Timer) *time.Timer {
	timer = time.NewTimer(time.Second)
	timer.Stop()
	return timer
}

func TickerFictitiousInit(ticker *time.Ticker) *time.Ticker {
	ticker = time.NewTicker(time.Second)
	ticker.Stop()
	return ticker
}
