package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	//"github.com/sirupsen/logrus"
)

// Postgres Postgres
type Postgres struct {
	accessKey   string
	secretKey   string
	bucket      string
	externalURL string
}

//func removeColorAndTrimSpace(s *string) {
//	flags := []string{"\033[0m", "\033[36;1m", "\033[33m[", "\033[31;1m", "\033[35m", "\033[36;31m"}
//	for _, flag := range flags {
//		*s = strings.ReplaceAll(*s, flag, "")
//	}
//	*s = strings.TrimSpace(*s)
//}
//
//func removeFilePathColorAndTrimSpace(s *string) {
//	removeColorAndTrimSpace(s)
//	*s = strings.TrimPrefix(*s, "(")
//	*s = strings.TrimSuffix(*s, ")")
//}
//
//type gormLogger struct{}
//
//func (gl *gormLogger) Print(v ...interface{}) {
//	gormFormattedLogs := gorm.LogFormatter(v...)
//	file := gormFormattedLogs[0].(string)
//	elapsed := gormFormattedLogs[2].(string)
//	l3 := gormFormattedLogs[3]
//	result := gormFormattedLogs[4].(string)
//	// bts, _ := json.MarshalIndent(gormFormattedLogs, "", "  ")
//	// fmt.Println(string(bts))
//
//	message := ""
//	isError := false
//	switch l3.(type) {
//	case error:
//		message = (l3.(error)).Error()
//		isError = true
//	case string:
//		message = l3.(string)
//	default:
//		message = fmt.Sprintf("%v", l3)
//	}
//
//	removeFilePathColorAndTrimSpace(&file)
//	removeColorAndTrimSpace(&elapsed)
//	removeColorAndTrimSpace(&result)
//	//fieldList := []logrus.Fields{
//	//	logrus.Fields{
//	//		"flag": "gorm",
//	//		"file": file,
//	//		"msg":  message,
//	//	},
//	//	logrus.Fields{
//	//		"flag": "gorm",
//	//		"file": file,
//	//		"msg":  result + elapsed,
//	//	},
//	//}
//
//
//
//
//	// bts, _ := json.MarshalIndent(fieldList, "", "  ")
//	// fmt.Println(string(bts))
//	if isError {
//		//logrus.WithFields(fieldList[0]).Error(nil)
//		zap.Logger.Info("gorm", zap.Any("file", string(file)), zap.Any("msg", message))
//
//
//	} else {
//		zap.Logger.Info("gorm", zap.Any("file", file), zap.Any("msg", result + elapsed))
//	}
//}

// InitPostgres InitPostgres
func InitPostgres(dbURL string, showLog bool, maxOpenConns, maxIdleConns int) (*gorm.DB, error) {
	db, err := gorm.Open("postgres", dbURL)
	db.DB().SetMaxOpenConns(maxOpenConns)
	db.DB().SetMaxIdleConns(maxIdleConns)
	db.LogMode(showLog)
	//db.SetLogger(&gormLogger{})
	return db, err
}

func InitMySQL(dbURL string, showLog bool, maxOpenConns, maxIdleConns int) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", dbURL)
	//audited.RegisterCallbacks(db)
	//_, err1 := loggable.Register(db)
	//if err1 != nil {
	//	panic(err1)
	//}

	db.DB().SetMaxOpenConns(maxOpenConns)
	db.DB().SetMaxIdleConns(maxIdleConns)
	db.LogMode(showLog)

	return db, err
}
