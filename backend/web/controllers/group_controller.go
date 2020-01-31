package controllers

import (
	"Kapi/models"
	"Kapi/services"
	"Kapi/utils"
	validators "Kapi/validator"
	"encoding/json"
	"github.com/kataras/iris"
	"gopkg.in/validator.v2"
)

type GroupController struct {
	Ctx iris.Context
	GroupService services.IGroupService
}

func (gc *GroupController) PostCreate() {
	uid := gc.Ctx.Values().Get("uid").(int64)
	member := &models.Member{
		Uid:     uid,
		Role:    "leader",
	}

	createGroupValidator := new(validators.CreateGroupValidator)
	dec := utils.NewDecoder(&utils.DecoderOptions{TagName:"kapi"})
	if err := dec.Decode(gc.Ctx.FormValues(),createGroupValidator); err != nil {
		utils.ErrorWithCode(err,"/group/create",400,gc.Ctx)
		return
	}
	if err := validator.Validate(createGroupValidator); err != nil {
		utils.ErrorWithCode(err,"/group/create",400,gc.Ctx)
		return
	}

	_, err := gc.GroupService.CreateGroup(createGroupValidator,member)
	if err != nil {
		utils.ErrorWithCode(err,"/group/create",500,gc.Ctx)
		return
	}
	gc.Ctx.JSON(utils.MakeDefaultRes(1,"分组创建成功",nil))
}
/* 创建小组 */
func (gc *GroupController) PostEdit() {
	updateGroupValidator := new(validators.UpdateGroupValidator)
	dec := utils.NewDecoder(&utils.DecoderOptions{TagName:"kapi"})
	if err := dec.Decode(gc.Ctx.FormValues(),updateGroupValidator); err != nil {
		utils.ErrorWithCode(err,"/group/edit",400,gc.Ctx)
		return
	}
	if err := validator.Validate(updateGroupValidator); err != nil {
		utils.ErrorWithCode(err,"/group/edit",400,gc.Ctx)
		return
	}

	_, err := gc.GroupService.UpdateGroup(gc.Ctx.Values().Get("uid").(int64), updateGroupValidator)
	if err != nil {
		utils.ErrorWithError(err,"/group/edit",gc.Ctx)
		return
	}
	gc.Ctx.JSON(utils.MakeDefaultRes(1,"更新分组信息成功！",nil))
}
/* 更改小组信息 */
func (gc *GroupController) PostGetmine() {
	uid := gc.Ctx.Values().Get("uid").(int64)
	groups, err := gc.GroupService.GetAllGroupsByUid(uid)
	if err != nil {
		utils.ErrorWithCode(err,"/group/getmine",500,gc.Ctx)
		return
	}
	data, err := json.Marshal(groups)
	if err != nil {
		utils.ErrorWithCode(err,"/group/getmine",500,gc.Ctx)
		return
	}
	gc.Ctx.JSON(utils.MakeDefaultRes(1,"查询成功", string(data)))
}
/* 获取我参加的所有小组 */
func (gc *GroupController) PostGetBy(id int64) {
	group, err := gc.GroupService.GetGroupById(id)
	if err != nil {
		utils.ErrorWithError(err,"/group/get/?",gc.Ctx)
		return
	}
	data, err := json.Marshal(group)
	if err != nil {
		utils.ErrorWithCode(err,"/group/get/?",500,gc.Ctx)
		return
	}
	gc.Ctx.JSON(utils.MakeDefaultRes(1,"查询成功", string(data)))
}
/* 获取小组信息 */
func (gc *GroupController) PostIsleaderBy(gid int64) {
	ok, _ := gc.GroupService.IsLeader(gc.Ctx.Values().Get("uid").(int64), gid)
	data := make(map[string]bool)
	data["isleader"] = ok
	gc.Ctx.JSON(utils.MakeDefaultRes(1,"身份查询成功！", data))
}
/* 是否为小组组长(决定前端显示界面) */

func (gc *GroupController) PostMemberJoin() {
	addMemberValidator := new(validators.AddMemberValidator)
	dec := utils.NewDecoder(&utils.DecoderOptions{TagName:"kapi"})
	if err := dec.Decode(gc.Ctx.FormValues(),addMemberValidator); err != nil {
		utils.ErrorWithCode(err,"/group/join",400,gc.Ctx)
		return
	}
	if err := validator.Validate(addMemberValidator); err != nil {
		utils.ErrorWithCode(err,"/group/join",400,gc.Ctx)
		return
	}

	_, err := gc.GroupService.AddGroupMember(gc.Ctx.Values().Get("uid").(int64),addMemberValidator)
	if err != nil {
		utils.ErrorWithError(err,"/group/join",gc.Ctx)
		return
	}
	gc.Ctx.JSON(utils.MakeDefaultRes(1,"添加组员成功！",nil))
}
/* 添加成员 */
func (gc *GroupController) PostMemberDelete() {
	delMemberValidator := new(validators.DelMemberValidator)
	dec := utils.NewDecoder(&utils.DecoderOptions{TagName:"kapi"})
	if err := dec.Decode(gc.Ctx.FormValues(),delMemberValidator); err != nil {
		utils.ErrorWithCode(err,"/group/delete",400,gc.Ctx)
		return
	}
	if err := validator.Validate(delMemberValidator); err != nil {
		utils.ErrorWithCode(err,"/group/delete",400,gc.Ctx)
		return
	}
	ok, err := gc.GroupService.DeleteMember(gc.Ctx.Values().Get("uid").(int64), delMemberValidator.Gid, delMemberValidator.Uid)
	if !ok || err != nil {
		utils.ErrorWithError(err,"/group/delete",gc.Ctx)
		return
	}
	gc.Ctx.JSON(utils.MakeDefaultRes(1,"删除成员成功！",nil))
}
/* 删除小组成员 */
func (gc *GroupController) PostMemberUpdate() {
	updateMembervalidator := new(validators.UpdateMemberValidator)
	dec := utils.NewDecoder(&utils.DecoderOptions{TagName:"kapi"})
	if err := dec.Decode(gc.Ctx.FormValues(),updateMembervalidator); err != nil {
		utils.ErrorWithCode(err,"/group/update",400,gc.Ctx)
		return
	}
	if err := validator.Validate(updateMembervalidator); err != nil {
		utils.ErrorWithCode(err,"/group/update",400,gc.Ctx)
		return
	}

	member := &models.Member{
		ID:      updateMembervalidator.ID,
		Gid:     updateMembervalidator.Gid,
		Uid:     updateMembervalidator.Uid,
		Role:    updateMembervalidator.Role,
	}

	_, err := gc.GroupService.UpdateMember(gc.Ctx.Values().Get("uid").(int64), member)
	if err != nil {
		utils.ErrorWithError(err, "/group/update", gc.Ctx)
	}
	gc.Ctx.JSON(utils.MakeDefaultRes(1,"更新组员信息成功！",nil))
}
/* 更新小组成员 */
func (gc *GroupController) PostMemberQueryBy(gid int64) {
	userMembers, err := gc.GroupService.SelectAllMembers(gid)
	if err != nil {
		utils.ErrorWithError(err,"/group/member/guery",gc.Ctx)
		return
	}
	data, err := json.Marshal(userMembers)
	if err != nil {
		utils.ErrorWithCode(err,"/group/member/guery",500,gc.Ctx)
		return
	}
	gc.Ctx.JSON(utils.MakeDefaultRes(1,"查询组员成功！",string(data)))
}
/* 查询小组所有组员 */