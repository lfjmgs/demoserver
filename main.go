package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type IPInfo struct {
	Code int `json:"code"`
	Data IP  `json:"data"`
}

type IP struct {
	Country   string `json:"country"`
	CountryId string `json:"country_id"`
	Area      string `json:"area"`
	AreaId    string `json:"area_id"`
	Region    string `json:"region"`
	RegionId  string `json:"region_id"`
	City      string `json:"city"`
	CityId    string `json:"city_id"`
	Isp       string `json:"isp"`
	IP        string `json:"ip"`
}

type Out struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type OutIPInfo struct {
	Out
	Data IP `json:"data"`
}

func TabaoAPI(ip string) *IPInfo {
	url := "http://ip.taobao.com/service/getIpInfo.php?ip="
	url += ip

	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	var result IPInfo
	if err := json.Unmarshal(out, &result); err != nil {
		return nil
	}

	return &result
}

func writeJson(w http.ResponseWriter, r interface{}) {
	bytes, err := json.Marshal(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func ipinfoHandler(w http.ResponseWriter, r *http.Request) {
	addr := r.RemoteAddr
	log.Println("RemoteAddr: " + addr)
	ip := string([]rune(addr)[:strings.LastIndex(addr, ":")])
	ipinfo := TabaoAPI(ip)
	var result OutIPInfo
	if ipinfo == nil {
		result.Status = "001"
		result.Message = "获取IP信息失败"
	} else {
		result.Status = "000"
		result.Message = "获取IP信息成功"
		result.Data = ipinfo.Data
		result.Data.IP = ip
	}
	writeJson(w, result)
}

func apierrHandler(w http.ResponseWriter, r *http.Request) {
	result := Out{Status: "001", Message: "请求失败"}
	writeJson(w, result)
}

func main() {
	http.HandleFunc("/ipinfo", ipinfoHandler)
	http.HandleFunc("/apierr", apierrHandler)
	log.Fatal(http.ListenAndServe(":10000", nil))
}
