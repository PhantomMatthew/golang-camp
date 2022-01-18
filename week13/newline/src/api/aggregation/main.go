package main

import (
	"aizinger.com/newline/src/common/log"
	"aizinger.com/newline/src/common/utils"
	"context"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/client"
	_ "github.com/micro/go-micro/v2/config"
	"go.uber.org/zap"

	//"github.com/sirupsen/logrus"
	"github.com/wantg/armory"
	"math/rand"
	"time"

	"aizinger.com/dzzs/docs"
	"aizinger.com/dzzs/reactors"
	config "aizinger.com/newline/src/common/config"
	customer "aizinger.com/newline/src/srv/customer/proto/customer"
	datavalidation "aizinger.com/newline/src/srv/validation/proto/datavalidation"
	"github.com/micro/go-micro/v2/web"
)

var (
	cl customer.CustomerService
	dv datavalidation.DatavalidationService
)

type Customer struct {
}
type Datavalidation struct {
}

func init() {
	rand.Seed(time.Now().UnixNano())
	//config.Init()
	reactors.Init()

	gin.DefaultWriter = armory.Log.DailyRotateLog(armory.Pilot.AppPath(config.Get("log.accessLog").String()))
	if config.Get("app.environment").String() == "production" {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultErrorWriter = armory.Log.DailyRotateLog(armory.Pilot.AppPath(config.Get("log.errorLog").String()))
	}
}

func main() {
	listen := config.Get("app.listen").String()

	service := web.NewService(
		web.Name("go.micro.api.customer"),
		web.Address(listen),
		web.Registry(utils.GetRegistry()),
	)

	service.Init()

	cl = customer.NewCustomerService("go.micro.service.customer", client.DefaultClient)
	dv = datavalidation.NewDatavalidationService("go.micro.service.datavalidation", client.DefaultClient)

	customer := new(Customer)
	datavalidation := new(Datavalidation)
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.Static("/assets", armory.Pilot.AppPath(config.Get("app.assets").String()))

	docConfig := map[string]string{
		"title":       config.Get("app.title").String(),
		"description": config.Get("app.description").String(),
		"version":     config.Get("app.version").String(),
		"listen":      config.Get("app.listen").String(),
		"username":    config.Get("doc.username").String(),
		"password":    config.Get("doc.password").String(),
	}
	if config.Get("app.environment").String() != "production" {
		docConfig["auth"] = "false"
	}
	docs.HandlerDoc(r, docConfig)

	//r.Use(middleware.AuthIdentify)

	r.GET("/customer/users", customer.CustomerGetUsersInfo)
	r.GET("/datavalidation/check", datavalidation.CheckCustomerCount)
	additionalRouters(r)

	service.Handle("/", r)

	//listen := config.Get("app.listen").String()
	//logrus.Info("listen on http://" + listen)
	//r.Run(listen)

	if err := service.Run(); err != nil {
		log.GetLogger().Error("service run error", zap.Error(err))
	}
}

func additionalRouters(r *gin.Engine) {
}

func (i *Customer) CustomerGetUsersInfo(c *gin.Context) {
	log.GetLogger().Info("Received Customer.UsersInfo API request")

	response, err := cl.GetCustPaInfo(context.TODO(), &customer.GetCustInfoReq{
		CustId: 232,
	})

	log.GetLogger().Info("response is", zap.Any("response", response))
	if err != nil {
		c.JSON(500, err)
	}

	c.JSON(200, response)
}

func (d *Datavalidation) CheckCustomerCount(c *gin.Context) {
	response, err := dv.CustomerInfoValidation(context.TODO(), &datavalidation.CivRequest{
		First: 1,
	})
	if err != nil {
		c.JSON(500, err)
	}

	// TODO TO BE NOTIFIED: it the return value is zero, it will not be shown in json
	c.JSON(200, response)

}
