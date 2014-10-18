/*
* File Name:	ipinfo_test.go
* Description:
* Author:	Chapman Ou <ochapman.cn@gmail.com>
* Created:	2014-10-12
 */
package ipinfo_test

import (
	"github.com/ochapman/ipinfo"
	"sync"
	"testing"
)

func TestIP(t *testing.T) {
	ips := [...]string{
		"8.8.8.8",
		"114.114.114.114",
	}
	for _, ip := range ips {
		if res, err := ipinfo.Get(ip); err != nil {
			t.Errorf("%s failed: %s\n", ip, err)
		} else {
			t.Logf("%s OK, info: %s\n", ip, res)
		}
	}
}

// Rest API acess limited: 10qps
func TestMultiReqIP(t *testing.T) {
	ip := "114.114.114.114"
	var wg sync.WaitGroup
	reqNum := 40
	wg.Add(reqNum)
	for i := 0; i < reqNum; i++ {
		go func(t *testing.T, ip string) {
			defer wg.Done()
			if res, err := ipinfo.Get(ip); err != nil {
				if err == ipinfo.ErrAPIUnavailable || err == ipinfo.ErrAPIBadWay {
					t.Logf("%s match err: %s\n", ip, err)
				} else {
					t.Errorf("%s failed: %s\n", ip, err)
				}
			} else {
				t.Logf("%s OK, info: %s\n", ip, res)
			}
		}(t, ip)
	}
	wg.Wait()
}

func TestInvalidIP(t *testing.T) {
	ips := [...]string{
		"114.114.114",
		"114.114.114.114.114",
	}
	for _, ip := range ips {
		if _, err := ipinfo.Get(ip); err != nil {
			if err != ipinfo.ErrInvalidIP {
				t.Errorf("%s failed: %s\n", ip, err)
			}
		}
	}
}

func TestPrivateIP(t *testing.T) {
	ips := [...]string{
		"192.168.0.1",
		"10.6.188.115",
	}
	for _, ip := range ips {
		if res, err := ipinfo.Get(ip); err != nil {
			if err == ipinfo.ErrPrivateIP {
				t.Logf("%s match err: %s\n", ip, err)
			} else {
				t.Errorf("%s failed: %s\n", ip, err)
			}
		} else {
			t.Logf("%s OK, info: %s\n", ip, res)
		}
	}
}

func TestNoInputIP(t *testing.T) {
	if _, err := ipinfo.Get(""); err != nil {
		if err == ipinfo.ErrNoInputIP {
			t.Logf("TestNoInputIP match err: %s\n", err)
		} else {
			t.Errorf("TestNoInputIP failed: %s\n", err)
		}
	} else {
		t.Errorf("TestNoInputIP failed: %s\n", err)
	}
}
