// Package cacheevict provides some cache eviction policy algorithms.
package cacheevict

type Cache interface {
	Add(string, any)
	Get(string) (any, bool)
}
