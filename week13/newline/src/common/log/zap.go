package log

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"newline.com/newline/src/common/config"
	"newline.com/newline/src/common/utils"
	"os"
	"path/filepath"
	"sync"
)

var (
	ZapLogger                      *Logger
	sp                             = string(filepath.Separator)
	errWS, warnWS, infoWS, debugWS zapcore.WriteSyncer       // IO输出
	debugConsoleWS                 = zapcore.Lock(os.Stdout) // 控制台标准输出
	errorConsoleWS                 = zapcore.Lock(os.Stderr)
	serviceName                    string
)

func Init(svcName string) {
	ZapLogger = &Logger{
		Opts: &Options{},
	}
	serviceName = svcName
	//internal.Register(initLogger)
	initLogger()
}

type Logger struct {
	*zap.Logger
	sync.RWMutex
	Opts      *Options `json:"opts"`
	zapConfig zap.Config
	inited    bool
}

func initLogger() {

	ZapLogger.Lock()
	defer ZapLogger.Unlock()

	if ZapLogger.inited {
		ZapLogger.Info("[initLogger] logger Inited")
		return
	}
	//config.Init()

	ZapLogger.loadCfg()
	ZapLogger.init()
	ZapLogger.Info("[initLogger] zap plugin initializing completed")
	ZapLogger.inited = true
}

// GetLogger returns logger
func GetLogger() (ret *Logger) {
	return ZapLogger
}

func (l *Logger) init() {

	l.setSyncers()
	var err error

	l.Logger, err = l.zapConfig.Build(l.cores())

	if err != nil {
		panic(err)
	}

	defer l.Logger.Sync()
}

func (l *Logger) loadCfg() {

	//c := config2.C()

	//err := c.Path("zap", l.Opts)
	//if err != nil {
	//	panic(err)
	//}
	//config.Get("zap")

	if config.Get("app.environment").String() == "production" {
		l.zapConfig = zap.NewProductionConfig()
	} else {
		l.zapConfig = zap.NewDevelopmentConfig()
	}

	//l.zapConfig.InitialFields = map[string]interface{}{"service name": serviceName, "service ip": utils.GetExternalIP()}
	//l.zapConfig.EncoderConfig.CallerKey = "caller"
	//// application log output path
	//if l.Opts.OutputPaths == nil || len(l.Opts.OutputPaths) == 0 {
	//	l.zapConfig.OutputPaths = []string{"stdout"}
	//}
	//
	////  error of zap-self log
	//if l.Opts.ErrorOutputPaths == nil || len(l.Opts.ErrorOutputPaths) == 0 {
	//	l.zapConfig.OutputPaths = []string{"stderr"}
	//}

	// 默认输出到程序运行目录的logs子目录
	if l.Opts.LogFileDir == "" {
		l.Opts.LogFileDir, _ = filepath.Abs(filepath.Dir(filepath.Join(".")))
		l.Opts.LogFileDir += sp + "logs" + sp
	}

	if l.Opts.AppName == "" {
		l.Opts.AppName = "app"
	}

	if l.Opts.FatalFileName == "" {
		l.Opts.FatalFileName = "fatal.log"
	}

	if l.Opts.ErrorFileName == "" {
		l.Opts.ErrorFileName = "error.log"
	}

	if l.Opts.WarnFileName == "" {
		l.Opts.WarnFileName = "warn.log"
	}

	if l.Opts.InfoFileName == "" {
		l.Opts.InfoFileName = "info.log"
	}

	if l.Opts.DebugFileName == "" {
		l.Opts.DebugFileName = "debug.log"
	}

	l.Opts.MaxSize = int(config.Get("log.max_size").Int())
	l.Opts.MaxBackups = int(config.Get("log.max_backups").Int())
	l.Opts.MaxAge = int(config.Get("log.max_age").Int())

	l.Opts.EnableKafka = config.Get("log.enable_kafka").Bool()

	//l.Opts.Level = json.Unmarshal(config.Get("log.level"), &l.Opts.Level)
	if err := json.Unmarshal([]byte(config.Get("log.level").Raw), &l.Opts.Level); err != nil {
		fmt.Println("Unmarshal log level error")
	}
}

func (l *Logger) setSyncers() {

	f := func(fN string) zapcore.WriteSyncer {
		return zapcore.AddSync(&lumberjack.Logger{
			Filename:   l.Opts.LogFileDir + sp + l.Opts.AppName + "-" + fN,
			MaxSize:    l.Opts.MaxSize,
			MaxBackups: l.Opts.MaxBackups,
			MaxAge:     l.Opts.MaxAge,
			Compress:   true,
			LocalTime:  true,
		})
	}

	errWS = f(l.Opts.ErrorFileName)
	warnWS = f(l.Opts.WarnFileName)
	infoWS = f(l.Opts.InfoFileName)
	debugWS = f(l.Opts.DebugFileName)

	return
}

func (l *Logger) cores() zap.Option {

	//l.zapConfig.InitialFields = map[string]interface{}{"service name": serviceName, "service ip": utils.GetExternalIP()}

	fileEncoder := zapcore.NewJSONEncoder(l.zapConfig.EncoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(l.zapConfig.EncoderConfig)

	fatalPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.FatalLevel && zapcore.FatalLevel-l.zapConfig.Level.Level() > -1
	})
	errPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.ErrorLevel && zapcore.ErrorLevel-l.zapConfig.Level.Level() > -1
	})
	warnPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel && zapcore.WarnLevel-l.zapConfig.Level.Level() > -1
	})
	infoPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel && zapcore.InfoLevel-l.zapConfig.Level.Level() > -1
	})
	debugPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.DebugLevel && zapcore.DebugLevel-l.zapConfig.Level.Level() > -1
	})

	var topicErrors zapcore.WriteSyncer
	var kafkaEncoder zapcore.Encoder

	l.Opts.EnableKafka = config.Get("log.enable_kafka").Bool()

	if l.Opts.EnableKafka {
		var (
			kl  LogKafka
			err error
		)
		kl.Topic = "service_log_topic"
		// 设置日志输入到Kafka的配置
		kconfig := sarama.NewConfig()
		//等待服务器所有副本都保存成功后的响应
		kconfig.Producer.RequiredAcks = sarama.WaitForAll
		//随机的分区类型
		kconfig.Producer.Partitioner = sarama.NewRandomPartitioner
		//是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoReponse这里才有用.
		kconfig.Producer.Return.Successes = true
		kconfig.Producer.Return.Errors = true

		kafkaAddrs := []string{}
		kafkaAddr := config.Get("kafka.host").String()
		kafkaAddrs = append(kafkaAddrs, kafkaAddr)
		kl.Producer, err = sarama.NewSyncProducer(kafkaAddrs, kconfig)
		if err != nil {
			ZapLogger.Info("connect kafka failed: %+v\n", zap.Error(err))
			os.Exit(-1)
		}
		topicErrors = zapcore.AddSync(&kl)
		// 打印在kafka
		if config.Get("app.environment").String() == "production" {
			kafkaEncoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
		} else {
			kafkaEncoder = zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig())

		}
	}

	cores := []zapcore.Core{}

	fields := []zapcore.Field{zap.Namespace("@service"), zap.String("ip", utils.GetExternalIP()), zap.String("name", serviceName)}

	if config.Get("app.environment").String() == "production" {
		cores = append(cores, zapcore.NewCore(fileEncoder, errWS, fatalPriority).With(fields))
		cores = append(cores, zapcore.NewCore(fileEncoder, errWS, errPriority).With(fields))
		cores = append(cores, zapcore.NewCore(fileEncoder, warnWS, warnPriority).With(fields))
		cores = append(cores, zapcore.NewCore(fileEncoder, infoWS, infoPriority).With(fields))
		cores = append(cores, zapcore.NewCore(consoleEncoder, debugConsoleWS, fatalPriority).With(fields))
		cores = append(cores, zapcore.NewCore(consoleEncoder, debugConsoleWS, errPriority).With(fields))
		cores = append(cores, zapcore.NewCore(consoleEncoder, debugConsoleWS, warnPriority).With(fields))
		cores = append(cores, zapcore.NewCore(consoleEncoder, debugConsoleWS, infoPriority).With(fields))
	} else {
		cores = append(cores, zapcore.NewCore(fileEncoder, errWS, fatalPriority).With(fields))
		cores = append(cores, zapcore.NewCore(fileEncoder, errWS, errPriority).With(fields))
		cores = append(cores, zapcore.NewCore(fileEncoder, warnWS, warnPriority).With(fields))
		cores = append(cores, zapcore.NewCore(fileEncoder, infoWS, infoPriority).With(fields))
		cores = append(cores, zapcore.NewCore(fileEncoder, debugWS, debugPriority).With(fields))

		cores = append(cores, zapcore.NewCore(consoleEncoder, debugConsoleWS, fatalPriority).With(fields))
		cores = append(cores, zapcore.NewCore(consoleEncoder, debugConsoleWS, errPriority).With(fields))
		cores = append(cores, zapcore.NewCore(consoleEncoder, debugConsoleWS, warnPriority).With(fields))
		cores = append(cores, zapcore.NewCore(consoleEncoder, debugConsoleWS, infoPriority).With(fields))
		cores = append(cores, zapcore.NewCore(consoleEncoder, debugConsoleWS, debugPriority).With(fields))
	}

	if kafkaEncoder != nil && topicErrors != nil {
		if config.Get("app.environment").String() == "production" {
			cores = append(cores, zapcore.NewCore(kafkaEncoder, topicErrors, fatalPriority).With(fields))
			cores = append(cores, zapcore.NewCore(kafkaEncoder, topicErrors, errPriority).With(fields))
			cores = append(cores, zapcore.NewCore(kafkaEncoder, topicErrors, warnPriority).With(fields))
			cores = append(cores, zapcore.NewCore(kafkaEncoder, topicErrors, infoPriority).With(fields))
		} else {
			cores = append(cores, zapcore.NewCore(kafkaEncoder, topicErrors, fatalPriority).With(fields))
			cores = append(cores, zapcore.NewCore(kafkaEncoder, topicErrors, errPriority).With(fields))
			cores = append(cores, zapcore.NewCore(kafkaEncoder, topicErrors, warnPriority).With(fields))
			cores = append(cores, zapcore.NewCore(kafkaEncoder, topicErrors, infoPriority).With(fields))
			cores = append(cores, zapcore.NewCore(kafkaEncoder, topicErrors, debugPriority).With(fields))
		}
	}

	return zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(cores...)
	})
}

type LogKafka struct {
	Producer sarama.SyncProducer
	Topic    string
}

func (lk *LogKafka) Write(p []byte) (n int, err error) {
	msg := &sarama.ProducerMessage{}
	msg.Topic = lk.Topic
	msg.Value = sarama.ByteEncoder(p)
	_, _, err = lk.Producer.SendMessage(msg)
	if err != nil {
		return
	}
	return
}
