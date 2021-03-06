// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/dxasu/gf.

package gtimer

import (
	"time"

	"github.com/dxasu/gf/container/gtype"
)

// Entry is the timing job entry to wheel.
type Entry struct {
	name              string
	wheel             *wheel      // Belonged wheel.
	job               JobFunc     // The job function.
	singleton         *gtype.Bool // Singleton mode.
	status            *gtype.Int  // Job status.
	times             *gtype.Int  // Limit running times.
	intervalTicks     int64       // The interval ticks of the job.
	createTicks       int64       // Timer ticks when the job installed.
	createMs          int64       // The timestamp in milliseconds when job installed.
	intervalMs        int64       // The interval milliseconds of the job.
	installIntervalMs int64       // Interval when first installation in milliseconds.
}

// JobFunc is the job function.
type JobFunc = func()

// addEntry adds a timing job to the wheel.
func (w *wheel) addEntry(interval time.Duration, job JobFunc, singleton bool, times int, status int) *Entry {
	if times <= 0 {
		times = defaultTimes
	}
	var (
		intervalMs    = interval.Nanoseconds() / 1e6
		intervalTicks = intervalMs / w.intervalMs
	)
	if intervalTicks == 0 {
		// If the given interval is lesser than the one of the wheel,
		// then sets it to one tick, which means it will be run in one interval.
		intervalTicks = 1
	}
	var (
		nowMs    = time.Now().UnixNano() / 1e6
		nowTicks = w.ticks.Val()
		entry    = &Entry{
			wheel:             w,
			job:               job,
			times:             gtype.NewInt(times),
			status:            gtype.NewInt(status),
			createTicks:       nowTicks,
			intervalTicks:     intervalTicks,
			singleton:         gtype.NewBool(singleton),
			createMs:          nowMs,
			intervalMs:        intervalMs,
			installIntervalMs: intervalMs,
		}
	)
	// Install the job to the list of the slot.
	w.slots[(nowTicks+intervalTicks)%w.number].PushBack(entry)
	return entry
}

// addEntryByParent adds a timing job with parent entry.
// The parameter `rollOn` specifies if just rolling on the entry, which was not met the runnable requirement
// and not executed previously. This is true often when the job internal is too long.
func (w *wheel) addEntryByParent(rollOn bool, nowMs, interval int64, parent *Entry) *Entry {
	intervalTicks := interval / w.intervalMs
	if intervalTicks == 0 {
		intervalTicks = 1
	}
	nowTicks := w.ticks.Val()
	entry := &Entry{
		name:              parent.name,
		wheel:             w,
		job:               parent.job,
		times:             parent.times,
		status:            parent.status,
		intervalTicks:     intervalTicks,
		singleton:         parent.singleton,
		createTicks:       nowTicks,
		createMs:          nowMs,
		intervalMs:        interval,
		installIntervalMs: parent.installIntervalMs,
	}
	if rollOn {
		entry.createMs = parent.createMs
		if parent.wheel.level == w.level {
			entry.createTicks = parent.createTicks
		}
	}
	w.slots[(nowTicks+intervalTicks)%w.number].PushBack(entry)
	return entry
}

// Status returns the status of the job.
func (entry *Entry) Status() int {
	return entry.status.Val()
}

// SetStatus custom sets the status for the job.
func (entry *Entry) SetStatus(status int) int {
	return entry.status.Set(status)
}

// Start starts the job.
func (entry *Entry) Start() {
	entry.status.Set(StatusReady)
}

// Stop stops the job.
func (entry *Entry) Stop() {
	entry.status.Set(StatusStopped)
}

//Reset reset the job.
func (entry *Entry) Reset() {
	entry.status.Set(StatusReset)
}

// Close closes the job, and then it will be removed from the timer.
func (entry *Entry) Close() {
	entry.status.Set(StatusClosed)
}

// IsSingleton checks and returns whether the job in singleton mode.
func (entry *Entry) IsSingleton() bool {
	return entry.singleton.Val()
}

// SetSingleton sets the job singleton mode.
func (entry *Entry) SetSingleton(enabled bool) {
	entry.singleton.Set(enabled)
}

// SetTimes sets the limit running times for the job.
func (entry *Entry) SetTimes(times int) {
	entry.times.Set(times)
}

// Run runs the job.
func (entry *Entry) Run() {
	entry.job()
}

// check checks if the job should be run in given ticks and timestamp milliseconds.
func (entry *Entry) check(nowTicks int64, nowMs int64) (runnable, addable bool) {
	switch entry.status.Val() {
	case StatusStopped:
		return false, true
	case StatusClosed:
		return false, false
	case StatusReset:
		return false, true
	}
	// Firstly checks using the ticks, this may be low precision as one tick is a little bit long.
	//if entry.name == "1" {
	//	intlog.Print("check:", nowTicks-entry.createTicks, nowTicks, entry.createTicks, entry.intervalTicks)
	//}
	if diff := nowTicks - entry.createTicks; diff > 0 && diff%entry.intervalTicks == 0 {
		// If not the lowest level wheel.
		if entry.wheel.level > 0 {
			diffMs := nowMs - entry.createMs
			switch {
			// Add it to the next slot, which means it will run on next interval.
			case diffMs < entry.wheel.timer.intervalMs:
				entry.wheel.slots[(nowTicks+entry.intervalTicks)%entry.wheel.number].PushBack(entry)
				return false, false

			// Normal rolls on the job.
			case diffMs >= entry.wheel.timer.intervalMs:
				// Calculate the leftover milliseconds,
				// if it is greater than the minimum interval, then re-install it.
				if leftMs := entry.intervalMs - diffMs; leftMs > entry.wheel.timer.intervalMs {
					// Re-calculate and re-installs the job proper slot.
					entry.wheel.timer.doAddEntryByParent(false, nowMs, leftMs, entry)
					return false, false
				}
			}
		}
		// Singleton mode check.
		if entry.IsSingleton() {
			// Note that it is atomic operation to ensure concurrent safety.
			if entry.status.Set(StatusRunning) == StatusRunning {
				return false, true
			}
		}
		// Limit running times.
		times := entry.times.Add(-1)
		if times <= 0 {
			// Note that it is atomic operation to ensure concurrent safety.
			if entry.status.Set(StatusClosed) == StatusClosed || times < 0 {
				return false, false
			}
		}
		// This means it does not limit the running times.
		// I know it's ugly, but it is surely high performance for running times limit.
		if times < 2000000000 && times > 1000000000 {
			entry.times.Set(defaultTimes)
		}
		//if entry.name == "1" {
		//	intlog.Print("runnable:", nowTicks-entry.createTicks, nowTicks, entry.createTicks, entry.createTicks, entry.interval)
		//}
		return true, true
	}
	return false, true
}
