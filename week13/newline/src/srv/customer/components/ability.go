package components

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"moul.io/zapgorm"
	mail "newline.com/newline/src/common/mail"
	sms "newline.com/newline/src/common/sms"
	"time"

	"github.com/jinzhu/gorm"

	"newline.com/newline/src/common/cache"
	"newline.com/newline/src/common/config"
	"newline.com/newline/src/common/database"
	"newline.com/newline/src/common/log"
)

var serviceName string

var MainDB *gorm.DB

var MySQLDb *gorm.DB

var Mongo *mongo.Client

// Redis Redis
var Redis cache.Redis

var smsAliyun sms.SMSAliyun

var mailer mail.Mailer

func init() {
	config.Init()
	serviceName = config.Get("name.customer").String()
	if serviceName == "" {
		serviceName = "go.micro.service.customerdefault"
	}

	log.Init(serviceName)
}

//func initPostgres() {
//	postgresqlURL := config.Get("persistence.postgresql.url").String()
//	postgresqlShowLog := config.Get("persistence.postgresql.showLog").Bool()
//	var err error
//	MainDB, err = database.InitPostgres(postgresqlURL, postgresqlShowLog)
//	if err != nil {
//		logrus.Error(err)
//		time.Sleep(5 * time.Second)
//		initPostgres()
//	}
//}

func initMySQL() {
	url := config.Get("persistence.mysql.url").String()
	showLog := config.Get("persistence.mysql.showLog").Bool()
	var err error
	maxOpenConns := int(config.Get("persistence.max_open_connections").Int())
	maxIdleConns := int(config.Get("persistence.max_open_connections").Int())

	MainDB, err = database.InitMySQL(url, showLog, maxOpenConns, maxIdleConns)
	MainDB.SetLogger(zapgorm.New(log.GetLogger().Logger))

	if err != nil {
		log.GetLogger().Error("Init MYSQL", zap.Error(err))
		time.Sleep(5 * time.Second)
		initMySQL()
	}
}

// Init Init
func Init() {
	//initPostgres()

	initMySQL()

	//initMongoDB()

	redisHost := config.Get("persistence.redis.host").String()
	Redis.Init(redisHost)

	smsAliyunAccessKeyID := config.Get("provider.aliyun-sms.accessKeyID").String()
	smsAliyunAccessKeySecret := config.Get("provider.aliyun-sms.accessKeySecret").String()
	smsAliyunRegionID := config.Get("provider.aliyun-sms.regionID").String()
	smsAliyunCaptchaSignName := config.Get("provider.aliyun-sms.captchaSignName").String()
	smsAliyunCaptchaTemplateCode := config.Get("provider.aliyun-sms.captchaTemplateCode").String()
	smsAliyun.Init(smsAliyunAccessKeyID, smsAliyunAccessKeySecret, smsAliyunRegionID, smsAliyunCaptchaSignName, smsAliyunCaptchaTemplateCode)

	mailerSMTPHost := config.Get("provider.qqmail.SMTPHost").String()
	mailerSMTPPort := config.Get("provider.qqmail.SMTPPort").Int()
	mailerAccount := config.Get("provider.qqmail.account").String()
	mailerPassword := config.Get("provider.qqmail.password").String()
	mailer.Init(mailerSMTPHost, int(mailerSMTPPort), mailerAccount, mailerPassword)
}

func GetServiceName() string {
	return serviceName
}
