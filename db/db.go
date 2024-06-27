package db

import (
	"ACT_GO/db/entities"
	"encoding/json"
	"fmt"
	"github.com/thoas/go-funk"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	DB *gorm.DB
)

func init() {
	var err error

	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=localhost user=postgres password=qwerty dbname=act sslmode=disable", // data source name, refer https://github.com/jackc/pgx
		PreferSimpleProtocol: true,                                                                      // disables implicit prepared statement usage. By default pgx automatically uses the extended protocol
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		log.Fatalf("DB init: %v", err)
	}

}

func TruncateLogs() *gorm.DB {
	return DB.Exec("truncate table app.log restart IDENTITY")
}

func Add_Log(log *entities.Log) {
	if log.StackTrace == "" {
		log.StackTrace = stacktrace()
	}
	if log.Module == "" {
		log.Module = caller()
	}

	DB.Create(log)
}

func Log_Cleanup() {
	DB.Raw("truncate table app.log")
}

func stacktrace() string {
	buf := make([]byte, 1<<20)
	stacklen := runtime.Stack(buf, false)
	return fmt.Sprintf("%s", buf[:stacklen])
}

func caller() string {
	pc, file, no, ok := runtime.Caller(2)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		return fmt.Sprintf("%s (%s line %d)", details.Name(), file, no)
	}
	return ""
}

func serialize_request_and_response(args ...interface{}) (req_json string, resp_json string) {
	var req interface{} = nil
	var resp interface{} = nil

	if len(args) > 0 {
		req = args[0]
	}
	if len(args) > 1 {
		resp = args[1]
	}

	req_json = ""
	if req != nil {
		switch (req).(type) {
		case string:
			req_json = (req).(string)
		default:
			bytes, _ := json.Marshal(req)
			req_json = string(bytes)
		}
	}

	resp_json = ""
	if resp != nil {
		switch (resp).(type) {
		case string:
			resp_json = (resp).(string)
		default:
			bytes, _ := json.Marshal(resp)
			resp_json = string(bytes)
		}
	}
	return
}

var m sync.Mutex

func AddMessage(msg string, args ...interface{}) {
	AddMessageWithModule(msg, "", args...)
}

func AddMessageWithModule(msg string, module string, args ...interface{}) {

	req_json, resp_json := serialize_request_and_response(args...)

	if module == "" {
		module = caller()
	}

	log_rec := &entities.Log{
		Category:   entities.LogPost,
		Message:    msg,
		Module:     module,
		StackTrace: stacktrace(),
		Request:    req_json,
		Response:   resp_json,
	}

	m.Lock()
	DB.Create(log_rec)
	time.Sleep(100 * time.Millisecond)
	m.Unlock()

}

func AddError(err error, tag string, args ...interface{}) {

	req_json, resp_json := serialize_request_and_response(args...)

	log_rec := &entities.Log{
		Category:   entities.LogError,
		Message:    err.Error(),
		Module:     caller(),
		Tag:        tag,
		StackTrace: stacktrace(),
		Request:    req_json,
		Response:   resp_json,
	}

	DB.Create(log_rec)
}

func AddErrors(errs []error, args ...interface{}) {

	req_json, resp_json := serialize_request_and_response(args...)

	log_rec := &entities.Log{
		Category:   entities.LogError,
		Message:    strings.Join(funk.Map(errs, func(e error) string { return e.Error() }).([]string), "; "),
		Module:     caller(),
		StackTrace: stacktrace(),
		Request:    req_json,
		Response:   resp_json,
	}

	DB.Create(log_rec)
}
