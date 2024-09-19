package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"
)

// 单个cpu利用率，[0～100]
var cpuUsage atomic.Int64

// 内存占用，单位Mb
var memoryUsage atomic.Int64

var memoryUsageBuf []byte

func init() {
	cpuUsage.Store(0)
	memoryUsage.Store(0)
	cpuUsageTask()
}

func main() {
	http.HandleFunc("/", handleDemoAPI)
	http.HandleFunc("/cpu", handleCpuUsgApi)
	http.HandleFunc("/memory", handleMemoryUsgApi)
	http.HandleFunc("/storage", handleStorageUsgApi)
	http.HandleFunc("/curl", handleCurlApi)
	http.HandleFunc("/http_status", handleHttpStatusAPI)
	http.HandleFunc("/delay", handleDelayAPI)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func handleHttpStatusAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	status := r.URL.Query().Get("status")
	statusCode, _ := strconv.Atoi(status)
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"code":0,"message":"success"}`))
}

func handleDelayAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	ms := r.URL.Query().Get("ms")
	msValue, _ := strconv.Atoi(ms)
	time.Sleep(time.Duration(msValue) * time.Millisecond)
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"code":0,"message":"success"}`))
}

func handleCurlApi(w http.ResponseWriter, r *http.Request) {
	targetUrl := r.URL.Query().Get("target_url")
	targetUrl, _ = url.QueryUnescape(targetUrl)
	resp, err := http.DefaultClient.Get(targetUrl)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	defer resp.Body.Close()
	w.WriteHeader(resp.StatusCode)
	data, _ := io.ReadAll(resp.Body)
	w.Write(data)
	return
}

// var startTime = time.Now()
func handleDemoAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// 存活探针
	//if time.Now().Sub(startTime) > 15*time.Minute {
	//	w.WriteHeader(500)
	//	w.Write([]byte{})
	//	return
	//}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Hello, world!"}`))
}

func handleStorageUsgApi(w http.ResponseWriter, r *http.Request) {
	usage := r.URL.Query().Get("usage") // 单位m
	v, err := strconv.ParseInt(usage, 10, 64)
	if err != nil || v < 0 {
		w.WriteHeader(500)
		w.Write([]byte("invalid step value"))
		return
	}
	fileName := fmt.Sprintf("/data/nebula-test-%s.txt", time.Now().Format("20060102150405"))
	f, err := os.Create(fileName)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	defer f.Close()
	m := 1 * 1024 * 1024
	info := []byte("hello word!\n")
	total := (int(v) * m) / len(info)
	for i := 0; i < total; i++ {
		if _, err = f.Write(info); err != nil && err != io.EOF {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
	}
	w.WriteHeader(200)
	w.Write([]byte(`write file success.`))
}

func handleCpuUsgApi(w http.ResponseWriter, r *http.Request) {
	usage := r.URL.Query().Get("usage")
	v, err := strconv.ParseInt(usage, 10, 64)
	if err != nil || v < 0 || v > 100 {
		w.WriteHeader(500)
		w.Write([]byte("invalid step value"))
		return
	}
	cpuUsage.Store(v)
	w.WriteHeader(200)
	w.Write([]byte("ok"))
}

// 不生效多调几次
func handleMemoryUsgApi(w http.ResponseWriter, r *http.Request) {
	usage := r.URL.Query().Get("usage")
	v, err := strconv.ParseInt(usage, 10, 64)
	if err != nil || v < 0 {
		w.WriteHeader(500)
		w.Write([]byte("invalid step value"))
		return
	}
	updateMemoryUsage(int(v))
	w.WriteHeader(200)
	w.Write([]byte("ok"))
}

func updateMemoryUsage(usageMb int) {
	mb := usageMb * 1024 * 1024
	memoryUsageBuf = nil
	runtime.GC()
	time.Sleep(200 * time.Millisecond)
	memoryUsageBuf = make([]byte, mb)
}

func cpuUsageTask() {
	go func() {
		for {
			usage := float64(cpuUsage.Load()) / float64(100)
			Compute(usage)
		}
	}()
}

// Compute 使用cpu占用率达到目标值
// usage [0, 1]， CPU利用率百分比
func Compute(usage float64) {
	// 一个总的CPU利用率的统计周期为1000毫秒
	var t = 1000.0
	// 总时间转换为纳秒
	var r int64 = 1000 * 1000
	totalNanoTime := t * (float64)(r)
	// 计算时间，纳秒
	runTime := totalNanoTime * usage
	// 休眠时间，纳秒
	sleepTime := totalNanoTime - runTime
	startTime := time.Now().UnixNano()
	// 运行
	for float64(time.Now().UnixNano())-float64(startTime) < runTime {
	}
	// 休眠
	time.Sleep(time.Duration(sleepTime) * time.Nanosecond)
}
