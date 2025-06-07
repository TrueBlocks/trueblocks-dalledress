package facets

import "sync/atomic"

type LoadState struct {
	fetching int32 // atomic bool (0 = false, 1 = true)
	loaded   int32 // atomic bool
	expected int32 // atomic int
}

func NewLoadState() *LoadState {
	return &LoadState{}
}

func (s *LoadState) StartFetching() bool {
	return atomic.CompareAndSwapInt32(&s.fetching, 0, 1)
}

func (s *LoadState) StopFetching() {
	atomic.StoreInt32(&s.fetching, 0)
	atomic.StoreInt32(&s.loaded, 1)
}

func (s *LoadState) IsFetching() bool {
	return atomic.LoadInt32(&s.fetching) == 1
}

func (s *LoadState) IsLoaded() bool {
	return atomic.LoadInt32(&s.loaded) == 1
}

func (s *LoadState) Reset() {
	atomic.StoreInt32(&s.fetching, 0)
	atomic.StoreInt32(&s.loaded, 0)
	atomic.StoreInt32(&s.expected, 0)
}
