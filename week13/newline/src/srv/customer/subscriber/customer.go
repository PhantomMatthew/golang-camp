package subscriber

import (
	"encoding/json"
	"github.com/micro/go-micro/v2/broker"
	"go.uber.org/zap"
	"newline.com/newline/src/srv/customer/bll"
	"newline.com/newline/src/srv/customer/handler/request"
	"newline.com/newline/src/srv/customer/model"

	//log "github.com/micro/go-micro/v2/logger"
	"newline.com/newline/src/common/log"
)

type ReceiveMsgHandler struct {
	bll *bll.CustBll
}

// New new a bll.
func NewReceiveMsgHandler() (d ReceiveMsgHandler) {
	d = ReceiveMsgHandler{
		bll: bll.NewCustBll(),
	}
	return
}

//func (e *Customer) Handle(ctx context.Context, msg *customer.Message) error {
//	log.Info("Handler Received message: ", msg.Say)
//
//	return nil
//}
//
//func Handler(ctx context.Context, msg *customer.Message) error {
//	log.Info("Function Received message: ", msg.Say)
//	return nil
//}
//func (e *Customer) Handle(ctx context.Context, msg *customer.Message) error {
//	log.Info("Handler Received message: ", msg.Say)
//
//	return nil
//}
//
//func Handler(ctx context.Context, msg *customer.Message) error {
//	log.Info("Function Received message: ", msg.Say)
//	return nil
//}

func (a ReceiveMsgHandler) ReceiveMsgHandler(p broker.Event) error {
	topic := p.Topic()
	//fmt.Println(topic)
	log.ZapLogger.Info("Handler Received message", zap.String("topic", topic))
	p.Topic()
	m := p.Message()
	msgBody := string(m.Body)
	log.ZapLogger.Info("msg body is ", zap.String("msg body", string(m.Body)))
	from := m.Header["key"]

	var action interface{}
	var custId *uint32
	var custBasicInfo *model.CustBasicInfo
	msg := FormatMsg([]byte(msgBody), from)

	//println(test)

	custBasicInfo = a.bll.UpsertCustomerByMsg(nil, &msg, from)

	if custBasicInfo != nil {
		custId = &custBasicInfo.ID
	}
	//// 没有unionid不处理SyncHistoryInfo
	//if custBasicInfo == nil {
	//	log.ZapLogger.Info("No CustBasicInfo ", zap.String("msg body", string(msgBody)))
	//
	//	return nil
	//}

	switch from {
	case "wx":
		a.bll.SyncScanRecord(nil, &msg.PaMsg, custId)
		action = &msg.PaMsg
	case "yz":
		// 同步订单信息
		if msg.YzMsg.Payload.Order != nil && msg.YzMsg.Payload.Order.Tid != "" {
			order, _ := a.bll.SyncEmOrder(nil, msg.YzMsg.Payload.Order, custId, "yz")
			// 补充订单归属信息
			a.bll.SyncOrderOwner(nil, order, custBasicInfo)
		}
		//同步退款单信息
		if msg.YzMsg.Payload.Refund != nil && msg.YzMsg.Payload.Refund.RefundID != "" {
			a.bll.SyncEmRight(nil, msg.YzMsg.Payload.Refund, "yz")
		}
		action = &msg.YzMsg
	case "wm":
		// 同步订单信息
		if msg.WmMsg.Payload.Order != nil && msg.WmMsg.Payload.Order.OrderNo != "" {
			order, _ := a.bll.SyncEmOrder(nil, msg.WmMsg.Payload.Order, custId, "wm")
			// 补充订单归属信息
			a.bll.SyncOrderOwner(nil, order, custBasicInfo)
		}
		//同步退款单信息
		if msg.WmMsg.Payload.Refund != nil && msg.WmMsg.Payload.Refund.ID != 0 {
			a.bll.SyncEmRight(nil, msg.WmMsg.Payload.Refund, "wm")
		}
		action = &msg.WmMsg
	case "mb":
		action = &msg.MbMsg
	case "dz":
		// a.bll.SyncDzRelationInfo(nil, &msg.DzReq, custBasicInfo.ID)
	case "qw":
		//a.bll.SyncDzInfo(nil, &dzMsg, custBasicInfo.ID)
		//action = &dzMsg
	}

	if custId != nil {
		a.bll.SyncAction(nil, action, custId, from)
	}

	return nil
}

func FormatMsg(msgBody []byte, from string) request.MsgMap {
	var msg request.MsgMap
	switch from {
	case "wx":
		var paMsg request.PaMsg
		json.Unmarshal(msgBody, &paMsg)
		msg.PaMsg = paMsg
	case "yz":
		var emMsg request.YzMsg
		json.Unmarshal(msgBody, &emMsg)
		msg.YzMsg = emMsg
	case "wm":
		var wmMsg request.WmMsg
		json.Unmarshal(msgBody, &wmMsg)
		msg.WmMsg = wmMsg
	case "mb":
		var mbMsg request.MbMsg
		json.Unmarshal(msgBody, &mbMsg)
		msg.MbMsg = mbMsg
	case "dz":
	case "qw":
	}

	return msg
}

func (a ReceiveMsgHandler) ReceiveQywechatMsgHandler(p broker.Event) error {
	log.ZapLogger.Info("Handler Received message through qywechat stream")

	m := p.Message()
	msgBody := string(m.Body)

	log.ZapLogger.Info("ReceiveQywechatMsgHandler", zap.String("msgBody", string(msgBody)))

	return nil
}
