package bll

import (
	"newline.com/newline/src/common/log"
	"newline.com/newline/src/common/utils"
	"newline.com/newline/src/srv/customer/handler/request"
	"newline.com/newline/src/srv/customer/model"
	//"newline.com/newline/src/srv/customer/utils"
	"context"
	"errors"
	"github.com/imdario/mergo"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"strings"
	"time"
)

//根据unionid获取客户基本信息
func (a *CustBll) FindCustomByUnionId(ctx context.Context, unionId string) (*model.CustBasicInfo, error) {
	custBasic, err := a.custBasicInfoDao.GetByUnionId(ctx, unionId)
	return custBasic, err
}

func (a *CustBll) UpsertCustomerByMsg(ctx context.Context, msgMap *request.MsgMap, from string) *model.CustBasicInfo {
	var cust interface{}
	ifNew := true

	switch from {
	case "wx":
		pa := msgMap.PaMsg.Payload.Customer
		//ifNew = msgMap.PaMsg.Payload.IfNewCustomer
		if pa == nil || pa.Openid == "" {
			return nil
		}
		cust = *pa
	case "yz":
		em := msgMap.YzMsg.Payload.Customer
		if em == nil || em.YzOpenID == "" {
			return nil
		}
		cust = *em
		//ifNew = msgMap.YzMsg.Payload.IfNewCustomer
	case "mb":
		mb := msgMap.MbMsg.Payload.Customer
		if mb == nil {
			return nil
		}
		cust = *mb
		//ifNew = msgMap.MbMsg.Payload.IfNewCustomer
	case "wm":
		em := msgMap.WmMsg.Payload.Customer
		if em == nil || em.Wid == "" {
			return nil
		}
		cust = *em
		//ifNew = msgMap.WmMsg.Payload.IfNewCustomer
		//case "dz":
		//	dzCustInfo := a.ConvertDzReqToDzInfo(&msgMap.DzReq)
		//	baseInfo, err = a.MergeDzInfo(ctx, dzCustInfo, msgMap.DzReq)
	}

	baseInfo := a.UpsertCustomer(ctx, cust, from, ifNew)

	return baseInfo
}

func (a *CustBll) UpsertCustomer(ctx context.Context, cust interface{}, from string, ifUpdate bool) *model.CustBasicInfo {
	var baseInfo *model.CustBasicInfo
	var tags []*model.CustTagInfo
	var err error

	switch from {
	case "wx":
		var paCustInfo *model.CustPaInfo
		paInfo := cust.(request.PaInfo)
		paCustInfo, baseInfo, tags = a.GetPaAndBaseInfo(&paInfo)
		//baseInfo, _, err = a.UpsertBaseInfo(ctx, baseInfo, from)
		baseInfo, _, err = a.UpsertBaseInfo(ctx, baseInfo, from, ifUpdate)
		existedPaInfo, err := a.custPaInfoDao.GetByOpenId(ctx, paInfo.Openid)

		if baseInfo != nil {
			paCustInfo.CustId = &baseInfo.ID
		}

		// 合并微信用户信息
		if gorm.IsRecordNotFoundError(err) {
			a.custPaInfoDao.Create(ctx, paCustInfo)
		} else {
			if existedPaInfo == nil {
				return nil
			}
			if ifUpdate {
				a.MergePaInfo(ctx, paCustInfo, existedPaInfo)
			}
		}

	case "yz":
		var emCustInfo *model.CustEmInfo
		var cards []*model.CustCardInfo
		emCustInfo, baseInfo, cards, tags = a.StructEmInfo(cust, "yz")
		baseInfo, _, err = a.UpsertBaseInfo(ctx, baseInfo, from, ifUpdate)
		existedEmInfo, err := a.custEmInfoDao.GetByOuterId(ctx, emCustInfo.OuterId)

		if baseInfo != nil {
			emCustInfo.CustId = &baseInfo.ID
		}

		// 合并微信用户信息
		if gorm.IsRecordNotFoundError(err) {
			a.custEmInfoDao.Create(ctx, emCustInfo)
		} else {
			if ifUpdate {
				if existedEmInfo == nil {
					return nil
				}
				a.MergeEmInfo(ctx, emCustInfo, existedEmInfo)
				// 如果识别出了用户身份需要更新相关表的cust_id
				if existedEmInfo.CustId != emCustInfo.CustId && emCustInfo.CustId != nil {
					a.custEmInfoDao.UpdateEmInfoCustId(ctx, emCustInfo.OuterId, emCustInfo.CustId)
				}
			}
		}

		if len(cards) > 0 {
			a.SyncCard(ctx, cards, baseInfo)
		}

	case "wm":
		var emCustInfo *model.CustEmInfo
		var cards []*model.CustCardInfo
		emCustInfo, baseInfo, cards, tags = a.StructEmInfo(cust, "wm")
		baseInfo, _, err = a.UpsertBaseInfo(ctx, baseInfo, from, ifUpdate)
		existedEmInfo, err := a.custEmInfoDao.GetByOuterId(ctx, emCustInfo.OuterId)

		if baseInfo != nil {
			emCustInfo.CustId = &baseInfo.ID
		}

		// 合并微信用户信息
		if gorm.IsRecordNotFoundError(err) {
			a.custEmInfoDao.Create(ctx, emCustInfo)
		} else {
			if ifUpdate {
				if existedEmInfo == nil {
					return nil
				}
				a.MergeEmInfo(ctx, emCustInfo, existedEmInfo)
				// 如果识别出了用户身份需要更新相关表的cust_id
				if existedEmInfo.CustId != emCustInfo.CustId && emCustInfo.CustId != nil {
					a.custEmInfoDao.UpdateEmInfoCustId(ctx, emCustInfo.OuterId, emCustInfo.CustId)
				}
			}
		}

		if len(cards) > 0 {
			a.SyncCard(ctx, cards, baseInfo)
		}
		//if baseInfo, _, err = a.UpsertBaseInfo(ctx, baseInfo, from, ifUpdate); err == nil {
		//	existedEmInfo, err := a.custEmInfoDao.GetByCustId(ctx, baseInfo.ID)
		//	// 合并微信用户信息
		//	if !gorm.IsRecordNotFoundError(err) {
		//		if ifUpdate {
		//			a.MergeEmInfo(ctx, emCustInfo, existedEmInfo)
		//		}
		//	} else {
		//		emCustInfo.CustId = baseInfo.ID
		//		a.custEmInfoDao.Create(ctx, emCustInfo)
		//	}
		//
		//	if len(cards) > 0 {
		//		a.SyncCard(ctx, cards, baseInfo.ID)
		//	}
		//}
	case "mb":
		var mbCustInfo *model.CustMbInfo
		mbInfo := cust.(request.MbInfo)
		mbCustInfo, baseInfo = a.GetMbAndBaseInfo(&mbInfo)
		if baseInfo, _, err = a.UpsertBaseInfo(ctx, baseInfo, from, ifUpdate); err == nil {
			existedMbInfo, err := a.custMbInfoDao.GetByCustId(ctx, baseInfo.ID)
			// 合并会员信息
			if gorm.IsRecordNotFoundError(err) {
				mbCustInfo.CustId = baseInfo.ID
				a.custMbInfoDao.Create(ctx, mbCustInfo)
			} else {
				if ifUpdate {
					if existedMbInfo == nil {
						return nil
					}
					a.MergeMbInfo(ctx, mbCustInfo, existedMbInfo)
				}
			}
		}
	case "dz":
		var dzCustInfo *model.CustDzInfo
		var dzRelationInfos []*model.CustDzRelationInfo

		//获取客户基本数据
		dzCustInfo, baseInfo, dzRelationInfos, tags = a.StructDzInfo(cust)

		if baseInfo != nil {
			baseInfo, _, err = a.UpsertBaseInfo(ctx, baseInfo, from, ifUpdate)
			if err != nil {
				log.GetLogger().Error("UpsertCustomer", zap.Error(err))
			}
			existedDzInfo, err := a.custDzInfoDao.GetByCustId(ctx, baseInfo.ID)
			//合并用户信息
			if !gorm.IsRecordNotFoundError(err) {
				if ifUpdate {
					a.MergeDzInfo(ctx, dzCustInfo, existedDzInfo)
				}
			} else {
				dzCustInfo.CustId = baseInfo.ID
				//dzCustInfo.UnionId = baseInfo.UnionId
				_, err := a.custDzInfoDao.Create(ctx, dzCustInfo)
				if err != nil {
					log.GetLogger().Error("UpsertCustomer", zap.Error(err))
				}
			}
			if len(dzRelationInfos) > 0 {
				err := a.SyncDzRelationInfo(ctx, dzRelationInfos, baseInfo.ID)
				if err != nil {
					log.GetLogger().Error("SyncDzRelationInfo", zap.Error(err))
				}
			}
			dzReq := cust.(request.DzInfo)
			//belongDate, err := utils.ParseStrToDateTime2(dzReq.LastDzBelongDate)
			//if belongDate == nil {
			//	belongDate, err = utils.ParseStrToDateTime2(dzReq.DzBelongDate)
			//}
			belongDate, err := utils.ParseStrToDateTime2(dzReq.DzBelongDate)
			if err != nil {
				log.GetLogger().Error("ParseStrToDateTime2", zap.Error(err))
			}
			err = a.AddOwnerHistory(dzCustInfo, baseInfo, belongDate, dzReq.MatchType)
			if err != nil {
				log.GetLogger().Error("AddOwnerHistory", zap.Error(err))
			}
			err = a.UpdateBelongOrder(baseInfo, belongDate)
			if err != nil {
				log.GetLogger().Error("AddOwnerHistory", zap.Error(err))
			}
		}

	}

	// 创建更新用户标签
	if len(tags) > 0 {
		a.SyncTag(ctx, tags, baseInfo)
	}

	return baseInfo
}

//修改CustOrderOwner表中的owner account信息
func (a *CustBll) UpdateBelongOrder(basicInfo *model.CustBasicInfo, belongTime *time.Time) error {

	custOrderOwnerList, err := a.custOrderOwnerDao.GetListByCustId(nil, basicInfo.ID)
	if err != nil {
		return err
	}

	if custOrderOwnerList != nil && len(*custOrderOwnerList) > 0 {
		for _, orderOwner := range *custOrderOwnerList {
			//订单日期>业绩归属日期,并且订单没有被指定业绩归属人,并且业绩没有被锁定
			if belongTime != nil && orderOwner.OrderTime.After(*belongTime) && !orderOwner.IsLocked {
				orderOwner.OwnerFrom = basicInfo.OwnerFrom
				orderOwner.OwnerAccountName = basicInfo.OwnerAccountName
				orderOwner.OwnerAccountId = basicInfo.OwnerAccountId
				_, err := a.custOrderOwnerDao.Update(nil, orderOwner.ID, &orderOwner)
				if err != nil {
					log.GetLogger().Error("UpdateBelongOrder", zap.Error(err))
					return err
				}
			}

		}
	}
	return nil
}

func (a *CustBll) AddOwnerHistory(dzInfo *model.CustDzInfo, basicInfo *model.CustBasicInfo, time2 *time.Time, matchType int32) error {
	effectiveDate := time.Now()

	maxDate := effectiveDate.AddDate(10, 0, 0)
	list, err := a.CustOwnerHistoryDao.GetListByCustId(nil, dzInfo.CustId, "")
	if !gorm.IsRecordNotFoundError(err) && len(*list) > 0 {
		if (*list)[0].OwnerAccountId != basicInfo.OwnerAccountId {
			//最新的业绩归属人和当前人不同，需要修改业绩有效期限
			item := (*list)[0]
			item.EffectiveEndTime = &effectiveDate
			_, err := a.CustOwnerHistoryDao.Update(nil, &item)
			if err != nil {
				log.GetLogger().Error("AddOwnerHistory", zap.Error(err))
			}
			if time2 != nil {
				effectiveDate = *time2
			}
			if matchType == 0 {
				matchType = 1 //不传值，默认1，系统自动匹配
			}
			//增加新的owner历史记录
			ownserHistory := model.CustOwnerHistory{
				CustId:             dzInfo.CustId,
				OwnerFrom:          basicInfo.OwnerFrom,
				OwnerAccountId:     basicInfo.OwnerAccountId,
				OwnerAccountName:   basicInfo.OwnerAccountName,
				EffectiveStartTime: &effectiveDate,
				EffectiveEndTime:   &maxDate,
				MatchType:          matchType,
			}
			_, err1 := a.CustOwnerHistoryDao.Create(nil, &ownserHistory)
			if err1 != nil {
				log.GetLogger().Error("AddOwnerHistory", zap.Error(err1))
			}
		}
	} else {
		if time2 != nil {
			effectiveDate = *time2
		}
		if matchType == 0 {
			matchType = 1 //不传值，默认1，系统自动匹配
		}
		//增加新的owner历史记录
		ownserHistory := model.CustOwnerHistory{
			CustId:             dzInfo.CustId,
			OwnerFrom:          basicInfo.OwnerFrom,
			OwnerAccountId:     basicInfo.OwnerAccountId,
			OwnerAccountName:   basicInfo.OwnerAccountName,
			EffectiveStartTime: &effectiveDate,
			EffectiveEndTime:   &maxDate,
			MatchType:          matchType,
		}
		_, err1 := a.CustOwnerHistoryDao.Create(nil, &ownserHistory)
		if err1 != nil {
			log.GetLogger().Error("AddOwnerHistory", zap.Error(err1))
		}
	}

	return err
}

func (a *CustBll) UpsertBaseInfo(ctx context.Context, info *model.CustBasicInfo, from string, ifUpdate bool) (*model.CustBasicInfo, bool, error) {
	if info.UnionId == "" {
		return nil, false, errors.New("Unionid is required")
	}

	var res *model.CustBasicInfo
	ifNew := false
	existedInfo, err := a.custBasicInfoDao.GetByUnionId(ctx, info.UnionId)
	fromUpper := strings.ToUpper(from)
	if gorm.IsRecordNotFoundError(err) {
		info.CreatedBy = fromUpper
		info.UpdatedBy = fromUpper
		a.custBasicInfoDao.Create(ctx, info)
		res = info
		ifNew = true
	} else if ifUpdate {
		info.CreatedBy = existedInfo.CreatedBy
		// 以新传入数据为主
		if from == "wx" || from == "mb" || existedInfo.UpdatedBy == "" || existedInfo.UpdatedBy == fromUpper {
			mergo.Merge(info, existedInfo)
			res = info
		} else {
			mergo.Merge(existedInfo, info)
			res = existedInfo
		}
		res.UpdatedBy = fromUpper

		//switch {
		//case from == "wx" || existedInfo.UpdatedBy == strings.ToUpper(from):
		//	mergo.Merge(existedInfo, info)
		//	res = existedInfo
		//case "yz":
		//	mergo.Merge(existedInfo, info)
		//	res = existedInfo
		//case "wm":
		//	mergo.Merge(existedInfo, info)
		//	res = existedInfo
		//case "mb":
		//	mergo.Merge(info, existedInfo)
		//	res = existedInfo
		//case "dz":
		//	mergo.Merge(existedInfo, info)
		//	res = existedInfo
		//}

		a.custBasicInfoDao.Save(ctx, res)
	} else {
		res = existedInfo
	}

	//res.ID = res.Model.ID
	return res, ifNew, nil
}

func (a *CustBll) SyncHistoryInfos(ctx context.Context, tasks string, lastOid string, pageSize int32, startDate string, endDate string) {
	taskList := strings.Split(tasks, ",")
	log.ZapLogger.Info("SyncHistoryInfos", zap.String("tasks", tasks))

	for _, v := range taskList {
		switch v {
		case "yz_customer":
			utils.PaggingLoop(pageSize, lastOid, startDate, endDate, a.BatchSyncYzCusts)
		case "yz_order_right":
			utils.PaggingLoop(pageSize, lastOid, startDate, endDate, a.BatchSyncYzRights)
			// 更新退款单和用户的绑定关系
			a.custOrderRightDao.UpdateOrderRightCustId(ctx)
		case "yz_order":
			utils.PaggingLoop(pageSize, lastOid, startDate, endDate, a.BatchSyncYzOrders)
			// 更新订单和用户的绑定关系
			a.custOrderRecordDao.UpdateOrderCustId(ctx)
		case "wm_customer":
			utils.PaggingLoop(pageSize, lastOid, startDate, endDate, a.BatchSyncWmCusts)
		case "wm_order_right":
			utils.PaggingLoop(pageSize, lastOid, startDate, endDate, a.BatchSyncWmRights)
			// 更新退款单和用户的绑定关系
			a.custOrderRightDao.UpdateOrderRightCustId(ctx)
		case "wm_order":
			utils.PaggingLoop(pageSize, lastOid, startDate, endDate, a.BatchSyncWmOrders)
			// 更新订单和用户的绑定关系
			a.custOrderRecordDao.UpdateOrderCustId(ctx)
		case "yz_customer_order":
			utils.PaggingLoop(pageSize, lastOid, startDate, endDate, a.BatchSyncYzCustomerAndOrders)
		case "wm_customer_order":
			utils.PaggingLoop(pageSize, lastOid, startDate, endDate, a.BatchSyncWmCustomerAndOrders)
		case "wx_customer":
			utils.PaggingLoop(pageSize, lastOid, startDate, endDate, a.BatchSyncPaInfos)
		case "mb_customer":
			utils.PaggingLoop(pageSize, lastOid, startDate, endDate, a.BatchSyncMbInfos)
		case "dz_customer":
			utils.PaggingLoop(pageSize, lastOid, startDate, endDate, a.BatchSyncDzInfos)
			//更新订单业绩归属
			a.custOrderOwnerDao.SyncOrderOwnerHistory(ctx)
		}
	}
}

func (a *CustBll) GetCustBasicInfo(custId uint32, unionId string) (*model.CustBasicInfo, error) {
	var custBasicInfo *model.CustBasicInfo
	var err error
	if custId != 0 {
		custBasicInfo, err = a.custBasicInfoDao.GetByCustId(nil, custId)
		if err != nil {
			log.GetLogger().Error("GetByCustId", zap.Error(err))
		}
	} else if unionId != "" {
		custBasicInfo, err = a.custBasicInfoDao.GetByUnionId(nil, unionId)
		if err != nil {
			log.GetLogger().Error("GetByUnionId", zap.Error(err))
		}
	}

	return custBasicInfo, err
}
func (a *CustBll) GetCustChannelInfo(id uint32, unionid string) (*model.CustChannelInfo, error) {
	return a.custBasicInfoDao.GetCustChannelInfo(id, unionid)
}
