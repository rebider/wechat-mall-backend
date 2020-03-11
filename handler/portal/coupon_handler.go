package portal

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"wechat-mall-backend/defs"
	"wechat-mall-backend/errs"
)

// 查询优惠券列表
func (h *Handler) GetCouponList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page, _ := strconv.Atoi(vars["page"])
	size, _ := strconv.Atoi(vars["size"])
	userId, _ := strconv.Atoi(vars["userId"])

	couponList, total := h.service.CouponService.GetCouponList(page, size, 1)
	voList := []defs.PortalCouponVO{}
	for _, v := range *couponList {
		couponLog := h.service.CouponService.QueryCouponLog(userId, v.Id)
		status := 0
		if couponLog.Id != 0 {
			status = 1
		}
		couponVO := defs.PortalCouponVO{}
		couponVO.Id = v.Id
		couponVO.Title = v.Title
		couponVO.FullMoney = v.FullMoney
		couponVO.Minus = v.Minus
		couponVO.Rate = v.Rate
		couponVO.Type = v.Type
		couponVO.StartTime = v.StartTime
		couponVO.EndTime = v.EndTime
		couponVO.Description = v.Description
		couponVO.Status = status
		voList = append(voList, couponVO)
	}
	resp := make(map[string]interface{})
	resp["list"] = voList
	resp["total"] = total
	defs.SendNormalResponse(w, resp)
}

// 领取优惠券
func (h *Handler) TakeCoupon(w http.ResponseWriter, r *http.Request) {
	req := defs.TakeCouponReq{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		panic(errs.ErrorParameterValidate)
	}
	coupon := h.service.CouponService.GetCouponById(req.CouponId)
	if coupon.Id == 0 {
		panic(errs.ErrorCoupon)
	}
	couponLog := h.service.CouponService.QueryCouponLog(req.UserId, req.CouponId)
	if couponLog.Id == 0 {
		panic(errs.NewErrorCoupon("请勿重复领取！"))
	}
	h.service.CouponService.RecordCouponLog(req.UserId, req.CouponId)
	defs.SendNormalResponse(w, "ok")
}

// 查询用户领取的优惠券
// status: 0-未使用 1-已使用 2-已过期
func (h *Handler) GetUserCouponList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, _ := strconv.Atoi(vars["userId"])
	status, _ := strconv.Atoi(vars["status"])
	page, _ := strconv.Atoi(vars["page"])
	size, _ := strconv.Atoi(vars["size"])

	voList, total := h.service.CouponService.QueryUserCoupon(userId, status, page, size)
	resp := make(map[string]interface{})
	resp["list"] = voList
	resp["total"] = total
	defs.SendNormalResponse(w, resp)
}

// 删除领取的优惠券
func (h *Handler) DoDeleteCouponLog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, _ := strconv.Atoi(vars["userId"])
	couponLogId, _ := strconv.Atoi(vars["id"])
	couponLog := h.service.CouponService.QueryCouponLog(userId, couponLogId)
	if couponLog.Id == 0 {
		panic(errs.NewErrorCoupon("未知的记录！"))
	}
	h.service.CouponService.DoDeleteCouponLog(couponLogId)
	defs.SendNormalResponse(w, "ok")
}
