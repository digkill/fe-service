package api4

import (
	"im/model"
	"im/utils"
	"net/http"
	"strconv"
)

func (api *API) InitPromo() {
	api.BaseRoutes.Promos.Handle("/status", api.ApiHandler(updatePromosStatuses)).Methods("PUT")
	api.BaseRoutes.Promos.Handle("/{promo_id:[A-Za-z0-9]+}", api.ApiHandler(updatePromo)).Methods("PUT")

	api.BaseRoutes.Promos.Handle("", api.ApiHandler(getAllPromos)).Methods("GET")
	api.BaseRoutes.Promos.Handle("", api.ApiHandler(createPromo)).Methods("POST")

	api.BaseRoutes.Promos.Handle("/{promo_id:[A-Za-z0-9]+}", api.ApiHandler(getPromo)).Methods("GET")
	api.BaseRoutes.Promo.Handle("", api.ApiHandler(deletePromo)).Methods("DELETE")

	api.BaseRoutes.Promo.Handle("/status", api.ApiHandler(updatePromoStatus)).Methods("PUT")

}

func updatePromosStatuses(c *Context, w http.ResponseWriter, r *http.Request) {

	if c.Err != nil {
		return
	}

	status := model.PromoStatusFromJson(r.Body)
	if status == nil {
		c.SetInvalidParam("status")
		return
	}

	// The user being updated in the payload must be the same one as indicated in the URL.
	if len(status.PromoIds) == 0 {
		c.SetInvalidParam("promo_ids")
		return
	}

	//product, err := c.App.GetProduct(c.Params.ProductId)
	/*if err == nil && product.Status == model.STATUS_OUT_OF_OFFICE && status.Status != model.STATUS_OUT_OF_OFFICE {
		//c.App.DisableAutoResponder(c.Params.UserId, c.IsSystemAdmin())
	}*/

	//c.App.Srv.Go(func() {
	for _, promoId := range status.PromoIds {
		c.App.UpdatePromoStatus(promoId, status)
	}
	//})

	ReturnStatusOK(w)
}

func updatePromoStatus(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequirePromoId()
	if c.Err != nil {
		return
	}

	status := model.PromoStatusFromJson(r.Body)
	if status == nil {
		c.SetInvalidParam("status")
		return
	}

	// The user being updated in the payload must be the same one as indicated in the URL.
	if status.PromoId != c.Params.PromoId {
		c.SetInvalidParam("promo_id")
		return
	}

	//product, err := c.App.GetProduct(c.Params.ProductId)
	/*if err == nil && product.Status == model.STATUS_OUT_OF_OFFICE && status.Status != model.STATUS_OUT_OF_OFFICE {
		//c.App.DisableAutoResponder(c.Params.UserId, c.IsSystemAdmin())
	}*/

	if promo, err := c.App.UpdatePromoStatus(c.Params.PromoId, status); err != nil {
		c.Err = err
		return
	} else {
		w.Write([]byte(promo.ToJson()))
	}
}

func getAllPromos(c *Context, w http.ResponseWriter, r *http.Request) {
	//c.RequireUserId()
	if c.Err != nil {
		return
	}

	appId := c.Params.AppId
	if len(appId) == 0 {
		if user, _ := c.App.GetUser(c.App.Session.UserId); user != nil {
			appId = user.AppId
		}
	}

	promoGetOptions := &model.PromoGetOptions{
		AppId:      appId,
		CategoryId: c.Params.CategoryId,
		OfficeId:   c.Params.OfficeId,
		Status:     c.Params.Status,
	}

	if active := r.URL.Query().Get("active"); active != "" {
		promoGetOptions.Active = &c.Params.Active
	}

	if len(promoGetOptions.AppId) == 0 {
		promoGetOptions.AppId = c.App.Session.AppId
	}

	if utils.StringInSlice(c.App.Session.Roles, []string{model.CHANNEL_USER_ROLE_ID, ""}) {
		promoGetOptions.Status = model.PROMO_STATUS_ACCEPTED
		promoGetOptions.Active = model.NewBool(true)
		promoGetOptions.Mobile = true
	}

	afterPromo := r.URL.Query().Get("after")
	beforePromo := r.URL.Query().Get("before")
	sinceString := r.URL.Query().Get("since")

	var since int64
	var parseError error

	if len(sinceString) > 0 {
		since, parseError = strconv.ParseInt(sinceString, 10, 64)
		if parseError != nil {
			c.SetInvalidParam("since")
			return
		}
	}

	/*	if !c.App.SessionHasPermissionToChannel(c.Session, c.Params.ChannelId, model.PERMISSION_READ_CHANNEL) {
		c.SetPermissionError(model.PERMISSION_READ_CHANNEL)
		return
	}*/

	var list *model.PromoList
	var err *model.AppError
	//etag := ""

	if since > 0 {
		list, err = c.App.GetAllPromosSince(since, promoGetOptions)
	} else if len(afterPromo) > 0 {

		list, err = c.App.GetAllPromosAfterPromo(afterPromo, c.Params.Page, c.Params.PerPage, promoGetOptions)
	} else if len(beforePromo) > 0 {

		list, err = c.App.GetAllPromosBeforePromo(beforePromo, c.Params.Page, c.Params.PerPage, promoGetOptions)
	} else {
		list, err = c.App.GetAllPromosPage(c.Params.Page, c.Params.PerPage, promoGetOptions)
	}

	if err != nil {
		c.Err = err
		return
	}

	/*	if len(etag) > 0 {
		w.Header().Set(model.HEADER_ETAG_SERVER, etag)
	}*/

	w.Write([]byte(list.ToJson()))
}

func getPromo(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequirePromoId()
	if c.Err != nil {
		return
	}

	promo, err := c.App.GetPromo(c.Params.PromoId)

	if err != nil {
		c.Err = err
		return
	}

	w.Write([]byte(promo.ToJson()))

}

func updatePromo(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequirePromoId()
	if c.Err != nil {
		return
	}

	promo := model.PromoFromJson(r.Body)

	if promo == nil {
		c.SetInvalidParam("promo")
		return
	}

	// The promo being updated in the payload must be the same one as indicated in the URL.
	if promo.Id != c.Params.PromoId {
		c.SetInvalidParam("id")
		return
	}

	promo.Id = c.Params.PromoId

	rpromo, err := c.App.UpdatePromo(promo, false)
	if err != nil {
		c.Err = err
		return
	}

	w.Write([]byte(rpromo.ToJson()))
}

func createPromo(c *Context, w http.ResponseWriter, r *http.Request) {

	promo := model.PromoFromJson(r.Body)
	if promo.AppId == "" {
		promo.AppId = c.App.Session.AppId
	}

	if promo == nil {
		c.SetInvalidParam("promo")
		return
	}

	result, err := c.App.CreatePromo(promo)
	if err != nil {
		c.Err = err
		return
	}
	w.Write([]byte(result.ToJson()))
}

func deletePromo(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequirePromoId()
	if c.Err != nil {
		return
	}

	_, err := c.App.GetPromo(c.Params.PromoId)
	if err != nil {
		c.SetPermissionError(model.PERMISSION_DELETE_POST)
		return
	}

	/*if c.App.Session.UserId == promo.UserId {
		if !c.App.SessionHasPermissionToChannel(c.App.Session, promo.ChannelId, model.PERMISSION_DELETE_POST) {
			c.SetPermissionError(model.PERMISSION_DELETE_POST)
			return
		}
	} else {
		if !c.App.SessionHasPermissionToChannel(c.App.Session, promo.ChannelId, model.PERMISSION_DELETE_OTHERS_POSTS) {
			c.SetPermissionError(model.PERMISSION_DELETE_OTHERS_POSTS)
			return
		}
	}*/

	if _, err := c.App.DeletePromo(c.Params.PromoId, c.App.Session.UserId); err != nil {
		c.Err = err
		return
	}

	ReturnStatusOK(w)
}
