package repository

import "sync/atomic"

type LoadState struct {
	loading  int32 // atomic bool (0 = false, 1 = true)
	loaded   int32 // atomic bool
	expected int32 // atomic int
}

func NewLoadState() *LoadState {
	return &LoadState{}
}

func (s *LoadState) StartLoading() bool {
	return atomic.CompareAndSwapInt32(&s.loading, 0, 1)
}

func (s *LoadState) StopLoading() {
	atomic.StoreInt32(&s.loading, 0)
	atomic.StoreInt32(&s.loaded, 1)
}

func (s *LoadState) IsLoading() bool {
	return atomic.LoadInt32(&s.loading) == 1
}

func (s *LoadState) IsLoaded() bool {
	return atomic.LoadInt32(&s.loaded) == 1
}

func (s *LoadState) Reset() {
	atomic.StoreInt32(&s.loading, 0)
	atomic.StoreInt32(&s.loaded, 0)
	atomic.StoreInt32(&s.expected, 0)
}
