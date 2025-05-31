package fake

import "time"

// WithChunkSize sets the size of each streaming chunk
func WithChunkSize(size int) func(*fakeModel) {
	return func(f *fakeModel) {
		f.ChunkSize = size
	}
}

// WithChunkDelay sets the delay between streaming chunks
func WithChunkDelay(delay time.Duration) func(*fakeModel) {
	return func(f *fakeModel) {
		f.ChunkDelay = delay
	}
}
