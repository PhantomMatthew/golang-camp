package routers

//import (
//	customer "newline.com/customer/srv/proto/customer"
//	"context"
//	"github.com/gin-gonic/gin"
//	"github.com/sirupsen/logrus"
//)
//
//
//type Customer struct {
//
//}
//
//
//
//// UsersInfo UsersInfo
//// @ID UsersInfo
//// @Tags Users
//// @Summary 所有用户信息
//// @Description
//// @Produce json
//// @Success 200 {object} models.Response
//// @Router /users/ [get]
//func (o *Customer) CustomerGetUsersInfo(svc *customer.CustomerService) gin.HandlerFunc {
//	fn := func(c *gin.Context) {
//		logrus.Print("Received Customer.UsersInfo API request")
//
//		name := c.Param("name")
//
//		response, err := (*svc).Call(context.TODO(), &customer.Request{
//			Name: name,
//		})
//
//		if err != nil {
//			c.JSON(500, err)
//		}
//
//		c.JSON(200, response)
//	}
//	return fn
//}
