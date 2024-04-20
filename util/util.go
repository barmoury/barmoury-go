package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"runtime/debug"
	"strings"
	"time"

	"github.com/barmoury/barmoury-go/cache"
)

func CacheWriteAlong[T any](bufferSize uint64, dateLastFlushed time.Time, cache cache.Cache[T], entry T) bool {
	cache.Cache(entry)
	diff := DateDiffInMinutes(dateLastFlushed, time.Now())
	return bufferSize >= cache.MaxBufferSize() || diff >= cache.IntervalBeforeFlush()
}

func DateDiffInMinutes(a time.Time, b time.Time) uint64 {
	d := b.Sub(a)
	return uint64(d.Minutes())
}

func StrFormat(str string, args ...any) string {
	return fmt.Sprintf(str, args...)
}

func GetOrDefault[T any](v T, d T) T {
	if reflect.ValueOf(v).IsZero() {
		return d
	}
	return v
}

func TranverseIterable(m interface{}, fn func(any, any)) {
	for k, v := range m.(map[string]map[string]any) {
		fn(k, v)
	}
}

func ValueInSlice[T comparable](list []T, a T) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func SplitByRegex(str string, regex string) []string {
	re := regexp.MustCompile(regex)
	sp := re.Split(str, -1)
	set := []string{}
	for i := range sp {
		set = append(set, sp[i])
	}
	return set
}

func SlicesInterception[T string | uint64](t1 []T, t2 []T) []T {
	rs := []T{}
	cm := map[T]struct{}{}

	for _, a := range t1 {
		cm[a] = struct{}{}
	}
	for _, a := range t2 {
		if _, ok := cm[a]; ok {
			rs = append(rs, a)
		}
	}
	return rs
}

func SlicesIntercepts[T string | uint64](t1 []T, t2 []T) bool {
	cm := map[T]struct{}{}

	for _, a := range t1 {
		cm[a] = struct{}{}
	}
	for _, a := range t2 {
		if _, ok := cm[a]; ok {
			return true
		}
	}
	return false
}

func If[T any](c bool, t T, f T) T {
	if c {
		return t
	}
	return f
}

func ReplaceByRegex(s string, r string, v string) string {
	var re = regexp.MustCompile(r)
	s = re.ReplaceAllString(s, v)
	return s
}

func PatternToRegex(p string) *regexp.Regexp {
	p = ReplaceByRegex(p, `:([\w])+`, "**")
	//t := If(!strings.Contains(p, "**"), `!((.)+)`, "")
	f := ReplaceByRegex(ReplaceByRegex(ReplaceByRegex(p, `\*\*`, `(.)+`), `\?`, `(.)`), `\*`, `((?!(\/)).)+`) + "" //t
	r, _ := regexp.Compile(f)
	return r
}

func StringToInterface(s string) interface{} {
	return reflect.ValueOf(s).Interface()
}

func InterfaceToString(s interface{}) string {
	return fmt.Sprintf("%s", reflect.ValueOf(s).Interface())
}

func StackTrace(skip int) []string {
	st := []string{}
	ds := strings.Split(StrFormat("%s", debug.Stack()), "\n")
	st = append(st, ds[0])
	st = append(st, ds[(skip*2)+1:]...)
	return st
}

func StackTraceAsString(skip int) string {
	c := ""
	st := StackTrace(skip)
	for _, s := range st {
		c += s + "\n"
	}
	return c
}

func RequestForResponse[T any](method string, url string, body []byte, headers map[string]string) (int, T) {
	r, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	for k, v := range headers {
		r.Header.Add(k, v)
	}
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	var resBody T
	derr := json.NewDecoder(res.Body).Decode(&resBody)
	if derr != nil {
		panic(derr)
	}
	return res.StatusCode, resBody
}

func PostRequestForResponse[T any](url string, body []byte, headers map[string]string) (int, T) {
	return RequestForResponse[T]("POST", url, body, headers)
}

func ToSnakeCase(str string) string {
	s := regexp.MustCompile("(.)([A-Z][a-z]+)").ReplaceAllString(str, "${1}_${2}")
	s = regexp.MustCompile("([a-z0-9])([A-Z])").ReplaceAllString(s, "${1}_${2}")
	return strings.ToLower(s)
}
