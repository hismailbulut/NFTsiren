package widgets

/*
import (
	"time"
)

type Animation struct {
	begin    time.Time // seconds
	duration time.Duration
	progress float32
}

func (anim *Animation) Start(now time.Time, duration time.Duration) {
	if now.IsZero() {
		panic("animation: now must not zero")
	}
	if duration <= 0 {
		panic("animation: duration must bigger than zero")
	}
	anim.begin = now
	anim.duration = duration
	anim.progress = 0
}

func (anim *Animation) Update(now time.Time) {
	if anim.duration == 0 {
		return
	}
	elapsed := now.Sub(anim.begin)
	anim.progress = clamp(float32(elapsed.Seconds()/anim.duration.Seconds()), 0, 1)
}

// Progress returns current animation progress between [0-1]
// 1 means it is finished or not started
func (anim *Animation) Progress() float32 {
	if anim.duration == 0 {
		return 1
	}
	return anim.progress
}
*/