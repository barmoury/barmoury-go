package main

import (
	"fmt"
	"time"

	"github.com/barmoury/barmoury-go/api/model"
	"github.com/barmoury/barmoury-go/audit"
	"github.com/barmoury/barmoury-go/cache"
	"github.com/barmoury/barmoury-go/crypto"
	"github.com/barmoury/barmoury-go/log"
	"github.com/barmoury/barmoury-go/util"
)

type User struct {
	//model model.Model
	Name string
}

func testModel() {
	model1 := new(model.Model)
	model2 := model.Model{}
	model3 := model.Model{Id: 40}
	//user1 := new(User)
	//user2 := User{Name: "Two"}
	model2.Resolve(nil, nil, nil)
	model2.Id = 30
	fmt.Println("The model1:", *model1.Resolve(model2, nil, nil))
	fmt.Println("The model1:", *model1.Resolve(model3, nil, nil))
	fmt.Println("The model2:", model2)
	fmt.Println("The model3:", model3)
}

func testCache() {
	listCache := cache.ListCacheImpl[uint]{}
	listCache.Cache(1)
	listCache.Cache(2)
	listCache.Cache(3)
	cachedValues := listCache.GetCached()
	listCache.Cache(4)
	cachedValues = append(cachedValues, 10)
	fmt.Println(listCache.MaxBufferSize(), listCache.IntervalBeforeFlush(), listCache, cachedValues)
}

func testTimeDiff() {
	t1 := time.Now().Add(-(time.Second * 341))
	fmt.Println(t1)

	diffMinutes := util.DateDiffInMinutes(t1, time.Now())
	fmt.Println(diffMinutes, "minutes")
}

type AuditorImpl struct {
	environment string
	serviceName string
	cache       cache.Cache[audit.Audit[string]]
	audit.Auditor[string]
}

func NewAuditorImpl() AuditorImpl {
	return AuditorImpl{
		environment: "local",
		serviceName: "barmoury",
		Auditor:     audit.NewAuditor[string](),
		cache:       &cache.ListCacheImpl[audit.Audit[string]]{},
	}
}

func (c *AuditorImpl) GetCache() cache.Cache[audit.Audit[string]] {
	return c.cache
}

func (c *AuditorImpl) PreAudit(a audit.Audit[string]) {
	a.Group = c.serviceName
	a.Environment = c.environment
}

func (c *AuditorImpl) Flush() {
	audits := c.GetCache().GetCached()
	fmt.Println("PREPARING TO FLUSH", audits)
}

func testAuditor() {
	auditor := NewAuditorImpl()
	auditor.Audit(&auditor, audit.Audit[string]{Action: "TESTING1"})
	auditor.Audit(&auditor, audit.Audit[string]{Action: "TESTING2"})
	auditor.Audit(&auditor, audit.Audit[string]{Action: "TESTING3"})
	auditor.Audit(&auditor, audit.Audit[string]{Action: "TESTING4"})
	auditor.Audit(&auditor, audit.Audit[string]{Action: "TESTING5"})
	auditor.Audit(&auditor, audit.Audit[string]{Action: "TESTING6"})
	auditor.Audit(&auditor, audit.Audit[string]{Action: "TESTING7"})
	auditor.Audit(&auditor, audit.Audit[string]{Action: "TESTING8"})
}

func testLog1() {
	fmt.Println(log.INFO)
	fmt.Println(log.WARN)
	fmt.Println(log.ERROR)
	fmt.Println(log.TRACE)
	fmt.Println(log.FATAL)
	fmt.Println(log.PANIC)
	fmt.Println(log.VERBOSE)
}

type LoggerImpl struct {
	environment string
	serviceName string
	cache       cache.Cache[log.Log]
	log.Logger
}

func NewLoggerImpl() LoggerImpl {
	return LoggerImpl{
		environment: "local",
		serviceName: "barmoury",
		Logger:      log.NewLogger(),
		cache:       &cache.ListCacheImpl[log.Log]{},
	}
}

func (c *LoggerImpl) GetCache() cache.Cache[log.Log] {
	return c.cache
}

func (c *LoggerImpl) PreLog(l log.Log) {
	l.Group = c.serviceName
	//fmt.Println(l.Level, ":", l.Content)
}

func (c *LoggerImpl) Flush() {
	logs := c.GetCache().GetCached()
	fmt.Println("PREPARING TO FLUSH LOGS", logs, "\n")
}

func testLog2() {
	logger := NewLoggerImpl()
	logger.Log(&logger, log.Log{Level: log.VERBOSE, Content: "This is the log for general"})
	logger.Info(&logger, "This is the log for the level %s", "info")
	logger.Warn(&logger, "This is the log for the level %s", "warn")
	logger.Trace(&logger, "This is the log for the level %s", "trace")
	logger.Verbose(&logger, "This is the log for the level %s", "verbose")
	//logger.Fatal(&logger, "This is the log for the level %s", "fatal")
	//logger.Panic(&logger, "This is the log for the level %s", "panic")
}

func testEncryptor2[T any](encryptor crypto.IEncryptor[T], value T) {
	fmt.Println("VALUE  :", value)
	e, ok := encryptor.Encrypt(value)
	if ok {
		fmt.Println("ENCRYPT:", e)
	}
	d, ok := encryptor.Decrypt(e)
	if ok {
		fmt.Println("DECRYPT:", d)
	}
}

func testEncryptor() {
	encrptor := crypto.ZlibCompressor[log.Log]{}
	value := log.Log{Level: log.PANIC, Content: "hello encryption"}
	testEncryptor2[log.Log](&encrptor, value)
}

func main() {
	//testModel()
	//testCache()
	//testTimeDiff()
	//testAuditor()
	//testLog1()
	//testLog2()
	testEncryptor()
}
