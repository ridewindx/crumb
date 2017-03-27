// Package dnscache caches DNS lookups.
package dnscache

import (
	"net"
	"sync"
	"time"
	"github.com/emirpasic/gods/maps/treemap"
)

type Resolver struct {
	lock  sync.RWMutex
	cache *treemap.Map
	size  int
}

func New(cacheSize int, refreshInterval time.Duration) *Resolver {
	resolver := &Resolver{
		cache: treemap.NewWithStringComparator(),
		size: cacheSize,
	}
	if refreshInterval > 0 {
		go resolver.autoRefresh(refreshInterval) // TODO: stop
	}
	return resolver
}

func (r *Resolver) Fetch(address string) ([]net.IP, error) {
	r.lock.RLock()
	ips, exists := r.cache.Get(address)
	r.lock.RUnlock()
	if exists {
		return ips.([]net.IP), nil
	}

	return r.Lookup(address)
}

func (r *Resolver) FetchOne(address string) (net.IP, error) {
	ips, err := r.Fetch(address)
	if err != nil || len(ips) == 0 {
		return nil, err
	}
	return ips[0], nil
}

func (r *Resolver) FetchOneString(address string) (string, error) {
	ip, err := r.FetchOne(address)
	if err != nil || ip == nil {
		return "", err
	}
	return ip.String(), nil
}

func (r *Resolver) Refresh() {
	r.lock.RLock()
	addresses := make([]string, r.cache.Size())
	i := 0
	for _, key := range r.cache.Keys() {
		addresses[i] = key.(string)
		i++
	}
	r.lock.RUnlock()

	for _, address := range addresses {
		r.Lookup(address)
		time.Sleep(time.Millisecond * 100) // TODO:
	}
}

func (r *Resolver) Lookup(address string) ([]net.IP, error) {
	ips, err := net.LookupIP(address)
	if err != nil {
		return nil, err
	}

	r.lock.Lock()
	r.cache.Put(address, ips)
	if r.size > 0 && r.cache.Size() > r.size {
		k, _ := r.cache.Min()
		r.cache.Remove(k)
	}
	r.lock.Unlock()
	return ips, nil
}

func (r *Resolver) autoRefresh(interval time.Duration) {
	for {
		time.Sleep(interval)
		r.Refresh()
	}
}
