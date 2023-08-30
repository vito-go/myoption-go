package service

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"math/rand"
	"myoption/iface/myerr"
	"net"
	"net/http"
	"regexp"
	"sync"
	"time"
)

func main() {
	addr, err := GetIpAddrByd777("202.99.231.3")
	if err != nil {
		panic(err)
	}

	fmt.Println(addr)
}

var getIpWith = [...]func(ip string) (string, error){
	GetIpAddrByipshu,
	GetIpAddrByipshudi,
	GetIpAddrByiplookup,
	GetIpAddrByhao86,
	GetIpAddrByd777,
}

type _ipLocalCache struct {
	mux  sync.Mutex
	data map[string]string
}

var ipLocalCache = _ipLocalCache{
	mux:  sync.Mutex{},
	data: make(map[string]string),
}

func GetIpAddrWithRandom(ip string) (string, error) {
	ipLocalCache.mux.Lock()
	defer ipLocalCache.mux.Unlock()
	address, ok := ipLocalCache.data[ip]
	if ok {
		return address, nil
	}
	rand.Seed(time.Now().UnixNano())
	address, err := getIpWith[rand.Intn(len(getIpWith))](ip)
	if err != nil {
		return "", err
	}
	ipLocalCache.data[ip] = address
	return address, nil
}
func getIpAddrBy(ipUrl string, reg *regexp.Regexp) (string, error) {
	method := "GET"
	tr := &http.Transport{
		// local error: tls: no renegotiation
		TLSClientConfig: &tls.Config{
			Renegotiation: tls.RenegotiateOnceAsClient,
			// You may need this if connecting to servers with self-signed certificates
			// InsecureSkipVerify: true,
		},

		//from the http.DefaultTransport
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: 3 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       3 * time.Second,
		TLSHandshakeTimeout:   3 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest(method, ipUrl, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	bodyString := string(body)
	result := reg.FindStringSubmatch(bodyString)
	if len(result) < 2 {
		return "", myerr.DataNotFound
	}
	return fmt.Sprintf("%s\n%s", result[1], ipUrl), nil
}

func GetIpAddrByipshu(ip string) (string, error) {
	url := fmt.Sprintf("https://zh-hans.ipshu.com/ipv4/%s", ip)
	reg := regexp.MustCompile("我们检测到该设备的物理位置位于(.+?)，上面的图片显示了")
	return getIpAddrBy(url, reg)
}
func GetIpAddrByipshudi(ip string) (string, error) {
	url := fmt.Sprintf("https://www.ipshudi.com/%s.htm", ip)
	reg := regexp.MustCompile(`归属地</td>[\S\s]+?<span>(.+?)</span>`)
	return getIpAddrBy(url, reg)
}
func GetIpAddrByiplookup(ip string) (string, error) {
	url := fmt.Sprintf("https://www.ip138.com/iplookup.php?ip=%s&action=2", ip)
	//ASN归属地</td>
	//<td><span>中国内蒙古赤峰喀喇沁旗</span>
	reg := regexp.MustCompile(`ASN归属地</td>[\S\s]+?<span>(.+?)</span>`)
	return getIpAddrBy(url, reg)
}
func GetIpAddrByhao86(ip string) (string, error) {
	url := fmt.Sprintf("https://ip.hao86.com/%s/", ip)
	//ASN归属地</td>
	//<td><span>中国内蒙古赤峰喀喇沁旗</span>
	reg := regexp.MustCompile(`地理定位</td>[\S\s]+?<td>(.+?)</td>`)
	return getIpAddrBy(url, reg)
}

func GetIpAddrByd777(ip string) (string, error) {
	url := fmt.Sprintf("https://ip.d777.com/%s_ip", ip)
	//ASN归属地</td>
	//<td><span>中国内蒙古赤峰喀喇沁旗</span>
	reg := regexp.MustCompile(`<br>IP详细地址:&nbsp;(.+?)<br/>`)
	return getIpAddrBy(url, reg)
}
