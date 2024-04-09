package services

import (
	"context"
	"devtools/backend/storage"
	"devtools/backend/types"
	"encoding/json"
	"fmt"
	"github.com/deckarep/golang-set/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type ipService struct {
	ctx           context.Context
	settings      *storage.SettingsStorage
	preferences   *storage.PreferencesStorage
	locationCache map[string]*types.DevLocation
}

var IPService *ipService
var onceIPService sync.Once

func IPServices(settings *storage.SettingsStorage, preferences *storage.PreferencesStorage) *ipService {
	if IPService == nil {
		onceIPService.Do(func() {
			IPService = &ipService{
				settings:      settings,
				preferences:   preferences,
				locationCache: make(map[string]*types.DevLocation),
			}
		})
	}
	return IPService
}

func (s *ipService) Start(ctx context.Context) {
	s.ctx = ctx
}

func (s *ipService) GetIP() (resp types.JSResp) {
	var wg sync.WaitGroup
	wg.Add(5)
	var taobaoIp string
	var ipIp string
	var upaiIp string
	var cloudFlareIp string
	var ipifyIp string
	go func() {
		defer wg.Done()
		taobaoIp = s.getTaobaoResult()
	}()
	go func() {
		defer wg.Done()
		ipIp = s.getIPIPResult()
	}()
	go func() {
		defer wg.Done()
		upaiIp = s.getUPaiResult()
	}()
	go func() {
		defer wg.Done()
		cloudFlareIp = s.getCloudFlareV4Result()
	}()
	go func() {
		defer wg.Done()
		ipifyIp = s.getIpifyV4Result()
	}()

	wg.Wait()
	ipSet := mapset.NewSet(taobaoIp, ipIp, upaiIp, cloudFlareIp, ipifyIp)

	for ip := range ipSet.Iter() {
		s.getIpDetail(ip)
	}
	ips := []*types.DevIP{
		{
			Source:   "TaoBao",
			IP:       taobaoIp,
			Location: s.locationCache[taobaoIp],
		},
		{
			Source:   "IPify IPv4",
			IP:       ipifyIp,
			Location: s.locationCache[ipifyIp],
		},
		{
			Source:   "UPaiYun",
			IP:       upaiIp,
			Location: s.locationCache[upaiIp],
		},
		{
			Source:   "CloudFlare IPv4",
			IP:       cloudFlareIp,
			Location: s.locationCache[cloudFlareIp],
		},
		{
			Source:   "IPIP",
			IP:       ipIp,
			Location: s.locationCache[ipIp],
		},
	}
	resp.Data = map[string]any{
		"ips": ips,
	}
	resp.Success = true
	return
}

func (s *ipService) getTaobaoResult() string {
	url := "https://www.taobao.com/help/getip.php?callback=ipCallback"
	response, err := http.Get(url)
	if err != nil {
		runtime.LogErrorf(s.ctx, "get ip from taobao failed %s", err.Error())
		return ""
	}
	//var data map[string]string
	responseData, err := io.ReadAll(response.Body)
	runtime.LogDebugf(s.ctx, "get taobao response :%s", responseData)
	ip := responseData[16 : len(responseData)-3]
	runtime.LogDebugf(s.ctx, "get taobao ip :%s", ip)
	//err = json.Unmarshal(responseString, &data)
	//if err != nil {
	//	runtime.LogErrorf(s.ctx, "unmarshal taobao ip result failed %s", err.Error())
	//	return ""
	//}
	//ip := data["ip"]
	return string(ip)
}

func (s *ipService) getIPIPResult() string {
	url := "https://myip.ipip.net/json"
	response, err := http.Get(url)
	if err != nil {
		runtime.LogErrorf(s.ctx, "get ip from ipip failed %s", err.Error())
		return ""
	}
	var data map[string]interface{}
	responseData, err := io.ReadAll(response.Body)
	err = json.Unmarshal(responseData, &data)
	if err != nil {
		runtime.LogErrorf(s.ctx, "unmarshal ipip ip result failed %s", err.Error())
		return ""
	}
	ip := data["data"].(map[string]interface{})["ip"].(string)
	return ip
}

func (s *ipService) getUPaiResult() string {
	url := fmt.Sprintf("https://pubstatic.b0.upaiyun.com/?_upnode&t=%d", time.Now().Unix())
	response, err := http.Get(url)
	if err != nil {
		runtime.LogErrorf(s.ctx, "get ip from upaiyun failed %s", err.Error())
		return ""
	}
	var data map[string]interface{}
	responseData, err := io.ReadAll(response.Body)
	err = json.Unmarshal(responseData, &data)
	if err != nil {
		runtime.LogErrorf(s.ctx, "unmarshal upaiyun ip result failed %s", err.Error())
		return ""
	}
	ip := data["remote_addr"].(string)
	if data["remote_addr_location"] != nil {
		runtime.LogDebugf(s.ctx, "get remote_addr_location:%v", data["remote_addr_location"])
		address := data["remote_addr_location"].(map[string]interface{})
		location := &types.DevLocation{
			Country:  address["country"].(string),
			Province: address["province"].(string),
			City:     address["city"].(string),
			ISP:      address["isp"].(string),
		}
		s.locationCache[ip] = location
	}
	return ip
}

func (s *ipService) getCloudFlareV4Result() string {
	url := "https://1.0.0.1/cdn-cgi/trace"
	response, err := http.Get(url)
	if err != nil {
		runtime.LogErrorf(s.ctx, "get ip from upaiyun failed %s", err.Error())
		return ""
	}
	responseData, err := io.ReadAll(response.Body)
	runtime.LogDebugf(s.ctx, "ip from cloud flare result %s", responseData)
	lines := strings.Split(string(responseData), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ip=") {
			ip := line[3:]
			return ip
		}
	}
	return ""
}

func (s *ipService) getIpifyV4Result() string {
	url := "https://api4.ipify.org?format=json"
	response, err := http.Get(url)
	if err != nil {
		runtime.LogErrorf(s.ctx, "get Ipify from upaiyun failed %s", err.Error())
		return ""
	}
	var data map[string]string
	responseData, err := io.ReadAll(response.Body)
	err = json.Unmarshal(responseData, &data)
	if err != nil {
		runtime.LogErrorf(s.ctx, "unmarshal ipip ip result failed %s", err.Error())
		return ""
	}
	ip := data["ip"]
	return ip
}

func (s *ipService) getIpDetail(ip string) *types.DevLocation {
	if location, ok := s.locationCache[ip]; ok {
		return location
	}
	location := s.fetchIpApiDetail(ip)
	s.locationCache[ip] = location
	return location
}

func (s *ipService) fetchIpApiDetail(ip string) *types.DevLocation {
	lang := s.preferences.GetLang()
	url := fmt.Sprintf("http://ip-api.com/json/%s?fields=66842623&lang=%s", ip, lang)
	response, err := http.Get(url)
	if err != nil {
		runtime.LogErrorf(s.ctx, "get Ipify from upaiyun failed %s", err.Error())
		return nil
	}
	var data map[string]string
	responseData, err := io.ReadAll(response.Body)
	err = json.Unmarshal(responseData, &data)
	location := &types.DevLocation{
		Country:  data["country"],
		Province: data["regionName"],
		City:     data["city"],
		ISP:      data["org"],
		ASN:      data["as"],
	}
	return location
}

//func (s *ipService) fetchIpDetail(ip string) string {
//	url := fmt.Sprintf("https://ipinfo.io/%s", ip)
//	token := s.settings.GetIPInfoToken()
//	if len(token) > 0 {
//		url = fmt.Sprintf("%s?token=%s", url, token)
//	}
//	response, err := http.Get(url)
//	if err != nil {
//		runtime.LogErrorf(s.ctx, "get Ipify from upaiyun failed %s", err.Error())
//		return ""
//	}
//	var data map[string]string
//	responseData, err := io.ReadAll(response.Body)
//	err = json.Unmarshal(responseData, &data)
//	if err != nil {
//		runtime.LogErrorf(s.ctx, "unmarshal ipip ip result failed %s", err.Error())
//		return ""
//	}
//	//ip := data["ip"]
//	return ip
//}
