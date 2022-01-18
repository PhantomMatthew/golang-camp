package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"moul.io/zapgorm"
	"newline.com/newline/src/common/database"
	"newline.com/newline/src/srv/customer/model"

	"math/rand"
	"reflect"
	"time"

	"newline.com/newline/src/common/config"
)

// go run cmd/migrate.go --config ../config/dev.json

var (
	mainDB *gorm.DB
	logger *zap.Logger
)

//func loadConfig() []byte {
//	configPtr := flag.String("c", "", "config file path")
//	flag.Parse()
//
//	configPath := strings.TrimSpace(*configPtr)
//	if len(configPath) == 0 {
//		panic("miss config file")
//	}
//
//	configPath, _ = filepath.Abs(configPath)
//	buf, err := ioutil.ReadFile(configPath)
//	if err != nil {
//		panic(err)
//	}
//	return buf
//}

func init() {
	logger, _ = zap.NewDevelopment()

	rand.Seed(time.Now().UnixNano())
	logger.Info("Initializing config")
	config.Init()
	logger.Info("Initializing db")

	mainDB = loadMySQLDB(config.Get("persistence.mysql.url").String(), true)
	mainDB.SetLogger(zapgorm.New(logger))
}

func loadPostgresDB(URL string, showLog bool) *gorm.DB {
	db, err := database.InitPostgres(URL, showLog, 20, 20)
	if err != nil {
		logger.Error("Initializing Postgres Error", zap.Error(err))
	}
	return db
}

func loadMySQLDB(URL string, showLog bool) *gorm.DB {
	logger.Info("db info", zap.String("url", URL))
	db, err := database.InitMySQL(URL, showLog, 20, 20)
	logger.Info("db info", zap.Any("db", db))
	if err != nil {
		logger.Error("Initializing MySQL Error", zap.Error(err))
	}
	return db
}

func loadModelFieldComments(model interface{}) map[string]string {
	comments := map[string]string{}
	modelType := reflect.TypeOf(model)
	scope := mainDB.NewScope(model)
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		if field.Type.Kind() == reflect.Struct {
			cs := loadModelFieldComments(reflect.New(field.Type).Elem().Interface())
			for f, c := range cs {
				comments[f] = c
			}
		}
		columnComment := field.Tag.Get("comment")
		if len(columnComment) > 0 {
			filed, _ := scope.FieldByName(field.Name)
			comments[filed.DBName] = columnComment
		}
	}
	return comments
}

func setupComments(modelList ...interface{}) {
	for _, model := range modelList {
		modelType := reflect.TypeOf(model)
		scope := mainDB.NewScope(model)
		tableName := scope.TableName()
		tableComment := ""
		for i := 0; i < modelType.NumField(); i++ {
			field := modelType.Field(i)
			tc := field.Tag.Get("table-comment")
			if len(tc) > 0 {
				tableComment = tc
			}
		}
		if len(tableComment) > 0 {
			mainDB.Exec(fmt.Sprintf("COMMENT ON TABLE %s IS '%s'", tableName, tableComment))
		}
		fieldComments := loadModelFieldComments(model)
		for fieldName, fieldComment := range fieldComments {
			mainDB.Exec(fmt.Sprintf("COMMENT ON COLUMN %s.%s IS '%s'", tableName, fieldName, fieldComment))
		}
	}
}

func main() {
	logger.Info("Start main")

	modelList := []interface{}{
		model.CustBasicInfo{},

	}
	//fmt.Println("33")
	logger.Info("Start Auto Migration")
	// mainDB.DropTableIfExists(modelList...)
	db := mainDB.AutoMigrate(modelList...)
	//fmt.Println(db.Error)
	if db.Error != nil {
		logger.Error("DB error", zap.Error(db.Error))
	}

	logger.Info("Start Setup Comments")
	setupComments(modelList...)
}
