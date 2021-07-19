package portmididrv

import "time"

// Option is an option that can be passed to the driver
type Option func(*Driver)

// SleepingTime sets the duration for sleeping between reads when polling on in port
// The default sleeping time is 0.1ms
func SleepingTime(d time.Duration) Option {
	return func(i *Driver) {
		i.sleepingTime = d
	}
}

/*
// BufferSize sets the size of the buffer when reading from in port
// The default buffersize is 1024
func BufferSizeRead(buffersize int) Option {
	return func(i *driver) {
		i.buffersizeRead = buffersize
	}
}

// BufferSize sets the size of the buffer when reading from in port
// The default buffersize is 1024
func BufferSizeIn(buffersize int64) Option {
	return func(i *driver) {
		i.buffersizeIn = buffersize
	}
}

// BufferSize sets the size of the buffer when reading from in port
// The default buffersize is 1024
func BufferSizeOut(buffersize int64) Option {
	return func(i *driver) {
		i.buffersizeOut = buffersize
	}
}
*/
