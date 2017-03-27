package dnscache

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

func TestFetchReturnsAndErrorOnInvalidLookup(t *testing.T) {
	_, err := New(0, 0).Lookup("invalid.crumb.io")
	if err == nil {
		return
	}
	expected := "lookup invalid.crumb.io: no such host"
	if err.Error() != expected {
		t.Errorf("Expecting %q error, got %q", expected, err.Error())
	}
}

func TestCallingLookupAddsTheItemToTheCache(t *testing.T) {
	r := New(0, 0)
	ips, _ := r.Lookup("dnscache.go.test.crumb.io")
	cached, _ := r.cache.Get("dnscache.go.test.crumb.io")
	assert.Equal(t, cached, ips)
}

func TestFetchLoadsValueFromTheCache(t *testing.T) {
	r := New(0, 0)
	ips := []net.IP{net.ParseIP("1.1.2.3")}
	r.cache.Put("invalid.crumb.io", ips)
	fetched, _ := r.Fetch("invalid.crumb.io")
	assert.Equal(t, fetched, ips)
}

func TestFetchOneLoadsTheFirstValue(t *testing.T) {
	r := New(0, 0)
	ips := []net.IP{net.ParseIP("1.1.2.3"), net.ParseIP("100.100.102.103")}
	r.cache.Put("something.crumb.io", ips)
	ip, _ := r.FetchOne("something.crumb.io")
	assert.Equal(t, ip, ips[0])
}

func TestFetchOneStringLoadsTheFirstValue(t *testing.T) {
	r := New(0, 0)
	ips := []net.IP{net.ParseIP("1.1.2.3"), net.ParseIP("100.100.102.103")}
	r.cache.Put("something.crumb.io", ips)
	ip, _ := r.FetchOneString("something.crumb.io")
	assert.Equal(t, ip, ips[0].String())
}

func TestItReloadsTheIpsAtAGivenInterval(t *testing.T) {
	r := New(0, time.Millisecond)
	r.cache.Put("www.baidu.com", nil)
	time.Sleep(5*time.Millisecond)
	_, exists := r.cache.Get("www.baidu.com")
	if !exists {
		t.Error("Expecting not-nil ips, got nil")
	}
}
