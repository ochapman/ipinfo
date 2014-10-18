/*
* File Name:	ipinfo.go
* Description:  Get IP information by Taobao REST API(ip.taobao.com)
* Author:	Chapman Ou <ochapman.cn@gmail.com>
* Created:	2014-10-12
 */
package ipinfo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const taobaoReqURL = "http://ip.taobao.com/service/getIpInfo.php?ip="

var (
	ErrInvalidIP      = errors.New("invalid IP")
	ErrNoInputIP      = errors.New("no input IP")
	ErrPrivateIP      = errors.New("private IP")
	ErrAPINoData      = errors.New("API no data")
	ErrAPIUnavailable = errors.New("API unavailable")
	ErrAPIBadWay      = errors.New("API bad way")
	ErrUnknown        = errors.New("unknown error")
)

//
// Response Sample:
// Sucess:
// {
//	"code":0,
//   	"data":
//		{
//		"country":"\u4e2d\u56fd",
//		"country_id":"CN",
//		"area":"\u534e\u5357",
//		"area_id":"800000",
//		"region":"\u5e7f\u4e1c\u7701",
//		"region_id":"440000",
//		"city":"\u6df1\u5733\u5e02",
//		"city_id":"440300",
//		"county":"",
//		"county_id":"-1",
//		"isp":"\u7535\u4fe1",
//		"isp_id":"100017",
//		"ip":"202.104.103.41"
//		}
// }
//
// Invalid:
// {
//	"code": 1,
//	"data": "invalid ip"
// }
//
// {
//	"code":0,
//	"data":
//		{
//		"country":"\u672a\u5206\u914d\u6216\u8005\u5185\u7f51IP",
//		"country_id":"IANA",
//		"area":"",
//		"area_id":"",
//		"region":"",
//		"region_id":"",
//		"city":"",
//		"city_id":"",
//		"county":"",
//		"county_id":"",
//		"isp":"",
//		"isp_id":"",
//		"ip":"0.114.114.114"
//		}
// }
//
//
// Service Temporarily Unavailable sample:
//
// <!DOCTYPE HTML PUBLIC "-//IETF//DTD HTML 2.0//EN">
// <html>
// <head><title>503 Service Temporarily Unavailable</title></head>
// <body bgcolor="white">
// <h1>503 Service Temporarily Unavailable</h1>
// <p>The server is temporarily unable to service your request due
// to maintenance downtime or capacity problems. Please try again later.
// <hr/>Powered by Tengine
// </body>
// </html>
//
//
// Bad Gateway
// <!DOCTYPE HTML PUBLIC "-//IETF//DTD HTML 2.0//EN">
// <html>
// <head><title>502 Bad Gateway</title></head>
// <body bgcolor="white">
// <h1>502 Bad Gateway</h1>
// <p>The proxy server received an invalid response from
// an upstream server.<hr/>Powered by Tengine
// </body>
// </html>

type IPInfo struct {
	Data data `json:data`
}

type responseCode struct {
	Code int `json:code`
}

func (r IPInfo) String() string {
	return fmt.Sprintf("%s %s %s %s %s",
		r.Data.IP, r.Data.Country, r.Data.Region, r.Data.City, r.Data.Isp)
}

type dataValid IPInfo

type dataInvalid struct {
	Data string `json:data`
}

type data struct {
	Country    string `json:country`
	Country_id string `json:country_id`
	Area       string `json:area`
	Area_id    string `json:area_id`
	Region     string `json:region`
	Region_id  string `json:region_id`
	City       string `json:city`
	City_id    string `json:city_id`
	County     string `json:conty`
	County_id  string `json:conty_id`
	Isp        string `json:isp`
	Isp_id     string `json:isp_id`
	IP         string `json:ip`
}

func Get(ip string) (info IPInfo, err error) {
	if ip == "" {
		err = ErrNoInputIP
		return
	}
	resp, err := http.Get(taobaoReqURL + ip)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	defer func() {
		if info.Data.Country == "未分配或者内网IP" {
			err = ErrPrivateIP
		}
	}()

	var buf []byte
	if buf, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	if len(buf) == 0 {
		err = ErrAPINoData
		return
	}
	if isAPIUnavailable(buf) {
		err = ErrAPIUnavailable
		return
	}
	if isAPIBadWay(buf) {
		err = ErrAPIBadWay
		return
	}
	var rc responseCode
	if err = json.Unmarshal(buf, &rc); err != nil {
		return
	}
	switch rc.Code {
	case 0:
		var d dataValid
		if err = json.Unmarshal(buf, &d); err != nil {
			return
		}
		info = IPInfo(d)
		return
	case 1:
		var d dataInvalid
		if err = json.Unmarshal(buf, &d); err != nil {
			return
		}
		err = ErrInvalidIP
		return
	default:
		err = ErrUnknown
		return
	}
	return
}

func isAPIUnavailable(buf []byte) bool {
	unvailable := []byte("503 Service Temporarily Unavailable")
	return bytes.Contains(buf, unvailable)
}

func isAPIBadWay(buf []byte) bool {
	badway := []byte("502 Bad Gateway")
	return bytes.Contains(buf, badway)
}
