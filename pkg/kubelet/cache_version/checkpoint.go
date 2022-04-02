/*
Just wrap read and write operation that from ioutil
*/

package cache_version

import (
	"io/ioutil"
	"os"

	"k8s.io/klog"
)

func PathExists(path string) (bool, error) {
	var err error
	var ret bool

	Block {
		Try: func() {
			_, err := os.Stat(path)
			if err == nil {
				ret = true
				err = nil
				return
			} else if os.IsNotExist(err) {
				ret = false
				err = nil
				return
			} else {
				Throw("Oh,...")
			}
		},
		Catch: func(e Exception) {
			ExceptionTriggered = true
			klog.Errorf("stat version.lock.save fail!, err: '%v'\n", e)
		},
		Finally: func() {
		},
	}.Do()

	if err == nil {
		return ret, nil
	}

	return false, err
}

func ReadFile(fname string) (string, error) {
	var err error
	var ret string

	Block {
		Try: func() {
			b, err := ioutil.ReadFile(fname)
			if err == nil {
				ret = string(b)
				err = nil
				return
			}

			Throw("Oh,...")
		},
		Catch: func(e Exception) {
			ExceptionTriggered = true
			klog.Errorf("read version.lock.save fail!, err: '%v'\n", e)
		},
		Finally: func() {
		},
	}.Do()

	if err == nil {
		return ret, nil
	}

	return "", err
}

func WriteFile(fname string, content string) error {
	var err error

	Block{
		Try: func() {
			bf := []byte(content)
			err = ioutil.WriteFile(fname, bf, 0644)
			if err == nil {
				return
			}

			Throw("Oh,...")
		},
		Catch: func(e Exception) {
			ExceptionTriggered = true
			klog.Errorf("write version.lock.save fail!, err: '%v', content:%s\n", e, content)
		},
		Finally: func() {
		},
	}.Do()

	return err
}