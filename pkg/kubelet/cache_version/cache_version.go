/*
CacheVersion, the target is to save version info to the checkpoint file.
write checkpoint when start up. and read checkpoint when version-comparation
*/

package cache_version

import (
	"fmt"
	"k8s.io/component-base/version"
	"strconv"
	"strings"
	"sync"
	"time"

	"k8s.io/klog/v2"
)

type CacheVersion struct {
	currentVersion          string //depend runtime. such current_version:v1.9.1
	currentVersionStartTime int64  //Unix()
	historyVersion          string //depend file. such history_version:v1.7.1
	historyVersionStartTime int64  //Unix()
}

var CacheVersionInfo = PCacheVersion()

const (
	DefaultCacheversion = "v0.0.0"
	DefaultStarttime    = 946656000 //  2000/01/01 00:00:00
)

func PCacheVersion() *CacheVersion {
	var nv CacheVersion

	nv.currentVersion = DefaultCacheversion
	nv.currentVersionStartTime = DefaultStarttime
	nv.historyVersion = DefaultCacheversion
	nv.historyVersionStartTime = DefaultStarttime

	return &nv
}

func CurrentVersion() string {
	return CacheVersionInfo.currentVersion
}

func HistoryVersion() string {
	return CacheVersionInfo.historyVersion
}

func CurrentVersionStartTime() int64 {
	return CacheVersionInfo.currentVersionStartTime
}

func HistoryVersionStartTime() int64 {
	return CacheVersionInfo.historyVersionStartTime
}

func GetCacheVersion() *CacheVersion {
	return CacheVersionInfo
}

func SetCachVersion() {
	SubCachVersion()
}

func SubCachVersion() {
	mutex.Lock()
	defer mutex.Unlock()

	version_file := fmt.Sprintf("/%s/%s", "tmp", "cache.version.save")

	exist, err := PathExists(version_file)
	if nil != err {
		klog.Errorf("stat version.lock.save fail!, err: '%v'", err)
		return
	}

	if exist {
		//file exist. set history version
		//a. just restart kubelet and not upgrade version
		//b. upgrade version(lower -> higher OR higher -> lower)
		GetFileContent(version_file)
	} else {
		//not exist.  set current_version == history_version
		//a. file was deleted by someone misoperation
		//b. first setup kubelet
		PersistVersion(version_file)
	}
}

func GetFileContent(fname string) {
	str, err := ReadFile(fname)
	if err != nil {
		return
	}

	//str = version,time
	str = strings.Replace(str, " ", "", -1) //remove space in head-tail
	vec := strings.Split(str, ",")
	l := len(vec)
	if 2 != l {
		return
	}

	CacheVersionInfo.historyVersion = vec[0]

	nt, err := strconv.ParseInt(vec[1], 10, 64)
	if err == nil {
		CacheVersionInfo.historyVersionStartTime = nt
	} else {
		klog.Errorf("parse version.lock.save fail!, err: '%v'", err)
	}

	//
	CacheVersionInfo.currentVersion = version.Get().String()
	CacheVersionInfo.currentVersionStartTime = time.Now().Unix()
}

func PersistVersion(fname string) {
	var v string
	var t int64

	v = version.Get().String()
	t = time.Now().Unix()

	CacheVersionInfo.currentVersion = v
	CacheVersionInfo.currentVersionStartTime = t
	CacheVersionInfo.historyVersion = v
	CacheVersionInfo.historyVersionStartTime = t

	str := fmt.Sprintf("%s,%d", v, t)
	WriteFile(fname, str)
}

var mutex sync.Mutex

func CheckValidity() bool {
	return (!ExceptionTriggered) &&
		(DefaultCacheversion != CurrentVersion())
}