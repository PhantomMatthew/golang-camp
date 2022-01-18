package config

import (
	"encoding/json"
	"flag"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Path struct {
	DataID string `json:"dataId"`
	Group  string `json:"group"`
}

type AcmConfig struct {
	Config     constant.ClientConfig `json:"config"`
	ConfigPath struct {
		Staging    Path `json:"staging"`
		Production Path `json:"production"`
	} `json:"configPath"`
}

func loadConfig() []byte {
	configPtr := flag.String("config", "", "config file path")
	flag.Parse()
	configPath := strings.TrimSpace(*configPtr)
	if len(configPath) == 0 {
		panic("miss config file")
	}
	kv := strings.Split(configPath, "path:")
	if len(kv) == 2 {
		configPath, _ = filepath.Abs(kv[1])
		buf, err := ioutil.ReadFile(configPath)
		if err != nil {
			panic(err)
		}
		return buf
	} else {
		acmConfigPath, _ := filepath.Abs("config/acm.json")
		buf, err := ioutil.ReadFile(acmConfigPath)
		if err != nil {
			panic(err)
		}
		var acmConfig AcmConfig
		json.Unmarshal(buf, &acmConfig)

		configClient, err := clients.CreateConfigClient(map[string]interface{}{
			"clientConfig": acmConfig.Config,
		})

		if err != nil {
			panic(err)
		}

		var path Path

		switch kv[0] {
		case "staging":
			path = acmConfig.ConfigPath.Staging
		case "production":
			path = acmConfig.ConfigPath.Production
		}

		content, err := configClient.GetConfig(vo.ConfigParam{
			DataId: path.DataID,
			Group:  path.Group,
		})

		return []byte(content)
	}

	//configPtr := flag.String("config", "", "config file path")
	//flag.Parse()
	//
	//configPath := strings.TrimSpace(*configPtr)
	////configPath := "./config/dev.json"

}

var config gjson.Result

//type logrusClassicTextFormatter struct {
//	base logrus.JSONFormatter
//}
//
//func (f *logrusClassicTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
//	colored := Get("app.environment").String() == "staging"
//	levelColor := "36m"
//	if entry.Level == logrus.DebugLevel {
//		levelColor = "32m"
//	} else if entry.Level == logrus.ErrorLevel {
//		levelColor = "31m"
//	}
//
//	caller := entry.Caller
//	timeNow := entry.Time.Format("15:04:05")
//	level := entry.Level.String()
//	filePath := caller.File + ":" + strconv.Itoa(caller.Line)
//	funcTitle := caller.Function
//	msg := entry.Message
//
//	if entry.Data["flag"] == "gorm" {
//		filePath = entry.Data["file"].(string)
//		funcTitle = entry.Data["flag"].(string) + ".func"
//		msg = entry.Data["msg"].(string)
//	}
//
//	r := fmt.Sprintf("[%s] [%-5s] %s %s %s\n",
//		timeNow, level, filePath, funcTitle, msg,
//	)
//	if colored {
//		r = fmt.Sprintf("[%s] \033[%s[%-5s]\033[0m %s \033[34m%s\033[0m \033[36;1m%s\033[0m\n",
//			timeNow, levelColor, level, filePath, funcTitle, msg,
//		)
//	}
//	//show all color
//	//tr := []string{}
//	//for idx := 0; idx < 100; idx++ {
//	//	tr = append(tr, fmt.Sprintf("\033[%dm [%dm] \033[0m", idx, idx))
//	//	tr = append(tr, fmt.Sprintf("\033[%d;1m [%d;1m] \033[0m", idx, idx))
//	//}
//	//r = strings.Join(tr, "")
//	return []byte(r), nil
//}

// Init Init
func Init() {
	config = gjson.ParseBytes(loadConfig())
	//initLogger()
}

//func initLogger() {
//	env := Get("app.environment").String()
//	if env == "staging" {
//		logrus.SetFormatter(&logrus.JSONFormatter{})
//		logrus.SetOutput(armory.Log.DailyRotateLog(armory.Pilot.AppPath(Get("log.runtimeLog").String())))
//	} else {
//		logrus.SetFormatter(&logrusClassicTextFormatter{})
//		logrus.SetOutput(os.Stdout)
//	}
//
//	logLevel := logrus.InfoLevel
//	switch Get("log.level").String() {
//	case "panic":
//		logLevel = logrus.PanicLevel
//	case "fatal":
//		logLevel = logrus.FatalLevel
//	case "error":
//		logLevel = logrus.ErrorLevel
//	case "warn":
//		logLevel = logrus.WarnLevel
//	case "info":
//		logLevel = logrus.InfoLevel
//	case "debug":
//		logLevel = logrus.DebugLevel
//	case "trace":
//		logLevel = logrus.TraceLevel
//	}
//	logrus.SetLevel(logLevel)
//	logrus.SetReportCaller(true)
//
//	logger := &lumberjack.Logger{
//		Filename:   Get("log.runtimeLog").String(),
//		// 日志文件最大 size, 单位是 MB
//		MaxSize:    Get("log.max_size").Index, // megabytes
//		// 最大过期日志保留的个数
//		MaxBackups: Get("log.max_backups").Index,
//		// 保留过期文件的最大时间间隔,单位是天
//		MaxAge:     Get("log.max_age").Index,   //days
//		// 是否需要压缩滚动日志, 使用的 gzip 压缩
//		Compress:   Get("log.compress").Bool(), // disabled by default
//	}
//	logrus.SetOutput(logger)
//}

// Get Get
func Get(path string) gjson.Result {
	return config.Get(path)
}
