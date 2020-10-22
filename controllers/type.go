package controllers

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/IrisAdminApi/libs"
	"github.com/snowlyg/IrisAdminApi/models"
	"github.com/snowlyg/IrisAdminApi/transformer"
	"github.com/snowlyg/IrisAdminApi/validates"
	gf "github.com/snowlyg/gotransformer"
)

/**
* @api {get} /admin/types/:id 根据id获取分类信息
* @apiName 根据id获取分类信息
* @apiGroup Types
* @apiVersion 1.0.0
* @apiDescription 根据id获取分类信息
* @apiSampleRequest /admin/types/:id
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
 */
func GetType(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	tt, err := models.GetTypeById(id)
	if err != nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(400, nil, err.Error()))
	}

	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(200, ttTransform(tt), "操作成功"))
}

/**
* @api {post} /admin/types/ 新建分类
* @apiName 新建分类
* @apiGroup Types
* @apiVersion 1.0.0
* @apiDescription 新建分类
* @apiSampleRequest /admin/types/
* @apiParam {string} name 分类名
* @apiParam {string} display_name
* @apiParam {string} description
* @apiParam {string} level
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiType null
 */
func CreateType(ctx iris.Context) {
	tt := new(models.Type)
	if err := ctx.ReadJSON(tt); err != nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(400, nil, err.Error()))
		return
	}
	err := validates.Validate.Struct(*tt)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs.Translate(validates.ValidateTrans) {
			if len(e) > 0 {
				ctx.StatusCode(iris.StatusOK)
				_, _ = ctx.JSON(ApiResource(400, nil, e))
				return
			}
		}
	}

	err = tt.CreateType()
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		_, _ = ctx.JSON(ApiResource(400, nil, fmt.Sprintf("Error create prem: %s", err.Error())))
		return
	}

	ctx.StatusCode(iris.StatusOK)
	if tt.ID == 0 {
		_, _ = ctx.JSON(ApiResource(400, tt, "操作失败"))
	} else {
		_, _ = ctx.JSON(ApiResource(200, ttTransform(tt), "操作成功"))
	}

}

/**
* @api {post} /admin/types/:id/update 更新分类
* @apiName 更新分类
* @apiGroup Types
* @apiVersion 1.0.0
* @apiDescription 更新分类
* @apiSampleRequest /admin/types/:id/update
* @apiParam {string} name 分类名
* @apiParam {string} display_name
* @apiParam {string} description
* @apiParam {string} level
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiType null
 */
func UpdateType(ctx iris.Context) {
	aul := new(models.Type)

	if err := ctx.ReadJSON(aul); err != nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(400, nil, err.Error()))
		return
	}
	err := validates.Validate.Struct(*aul)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs.Translate(validates.ValidateTrans) {
			if len(e) > 0 {
				ctx.StatusCode(iris.StatusOK)
				_, _ = ctx.JSON(ApiResource(400, nil, e))
				return
			}
		}
	}

	id, _ := ctx.Params().GetUint("id")
	aul.ID = id
	err = models.UpdateTypeById(id, aul)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		_, _ = ctx.JSON(ApiResource(400, nil, fmt.Sprintf("Error update type: %s", err.Error())))
		return
	}

	ctx.StatusCode(iris.StatusOK)
	if aul.ID == 0 {
		_, _ = ctx.JSON(ApiResource(400, nil, "操作失败"))
	} else {
		_, _ = ctx.JSON(ApiResource(200, ttTransform(aul), "操作成功"))
	}

}

/**
* @api {delete} /admin/types/:id/delete 删除分类
* @apiName 删除分类
* @apiGroup Types
* @apiVersion 1.0.0
* @apiDescription 删除分类
* @apiSampleRequest /admin/types/:id/delete
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiType null
 */
func DeleteType(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	err := models.DeleteTypeById(id)
	if err != nil {

		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(400, nil, err.Error()))
	}
	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(200, nil, "删除成功"))
}

/**
* @api {get} /tts 获取所有的分类
* @apiName 获取所有的分类
* @apiGroup Types
* @apiVersion 1.0.0
* @apiDescription 获取所有的分类
* @apiSampleRequest /tts
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiType null
 */
func GetAllTypes(ctx iris.Context) {
	offset := libs.ParseInt(ctx.URLParam("offset"), 1)
	limit := libs.ParseInt(ctx.URLParam("limit"), 20)
	name := ctx.FormValue("searchStr")
	orderBy := ctx.FormValue("orderBy")

	tts, err := models.GetAllTypes(name, orderBy, offset, limit)
	if err != nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(400, nil, err.Error()))
	}

	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(200, ttsTransform(tts), "操作成功"))
}

func ttsTransform(tts []*models.Type) []*transformer.Type {
	var rs []*transformer.Type
	for _, tt := range tts {
		r := ttTransform(tt)
		rs = append(rs, r)
	}
	return rs
}

func ttTransform(tt *models.Type) *transformer.Type {
	r := &transformer.Type{}
	g := gf.NewTransform(r, tt, time.RFC3339)
	_ = g.Transformer()
	return r
}