package services

import (
	"Kapi/models"
	"Kapi/repositories"
	"Kapi/utils"
	validators "Kapi/validator"
)

type IGroupService interface {
	CreateGroup(group *validators.CreateGroupValidator, member *models.Member) (gid int64, err error)
	UpdateGroup(uid int64,group *validators.UpdateGroupValidator) (gid int64, err error)
	GetAllGroupsByUid (uid int64) (groups []*models.Group, err error)
	GetGroupById(gid int64) (group *models.Group, err error)
	IsLeader(myId, gid int64) (ok bool, err error)

	AddGroupMember(uid int64, input *validators.AddMemberValidator) (mid int64, err error)
	DeleteMember(myId, gid, uid int64) (ok bool, err error)
	UpdateMember(myId int64, member *models.Member) (mid int64, err error)
	SelectAllMembers(gid int64) (userMembers []*models.UserMember, err error)
}

type GroupService struct {
	userRepository repositories.IUserManager
	groupRepository	repositories.IGroupManager
	memberRepository repositories.IMemberManager
}

func NewGroupService(userRepository repositories.IUserManager, groupRepository repositories.IGroupManager, memberRepository repositories.IMemberManager) IGroupService {
	return &GroupService{
		userRepository:userRepository,
		groupRepository: groupRepository,
		memberRepository: memberRepository,
	}
}

func (gs *GroupService) CreateGroup(group *validators.CreateGroupValidator, member *models.Member) (gid int64, err error) {
	return gs.groupRepository.InsertOne(group, member)
}

func (gs *GroupService) UpdateGroup(uid int64, group *validators.UpdateGroupValidator) (gid int64, err error) {
	isAllowed, err := gs.memberRepository.IsAllowed(uid, group.ID)
	if !isAllowed {
		if err != nil {
			err = utils.NewError(500, err.Error())
		}else{
			err = utils.NewError(400, "操作权限不够")
		}
		return
	}

	uid, err = gs.groupRepository.UpdateOne(group)
	if err != nil {
		err = utils.NewError(500, err.Error())
		return
	}
	return
}

func (gs *GroupService) GetAllGroupsByUid(uid int64) (groups []*models.Group, err error) {
	return gs.groupRepository.SelectAllByUid(uid)
}

func (gs *GroupService) GetGroupById(gid int64) (group *models.Group, err error) {
	group, err = gs.groupRepository.SelectOne(gid)
	if err != nil {
		if err.Error() == "没有查询到有效信息" {
			err = utils.NewError(400, err.Error())
		}else{
			err = utils.NewError(500, err.Error())
		}
		return
	}
	return
}

func (gs *GroupService) IsLeader(myId, gid int64) (ok bool, err error) {
	isAllowed, err := gs.memberRepository.IsAllowed(myId, gid)
	if !isAllowed {
		if err != nil {
			err = utils.NewError(500, err.Error())
		}else{
			err = utils.NewError(400, "操作权限不够")
		}
		return
	}
	return true, err
}

func (gs *GroupService) AddGroupMember(uid int64, input *validators.AddMemberValidator) (mid int64, err error) {
	isAllowed, err := gs.memberRepository.IsAllowed(uid, input.Gid)
	if !isAllowed {
		if err != nil {
			err = utils.NewError(500, err.Error())
		}else{
			err = utils.NewError(400, "操作权限不够")
		}
		return
	}
	user, err := gs.userRepository.SelectOne(input.Email)
	if err != nil {
		if err.Error() == "没有查询到有效信息" {
			err = utils.NewError(400, err.Error())
		} else {
			err = utils.NewError(500, err.Error())
		}
		return
	}

	member := &models.Member{
		Gid:     input.Gid,
		Uid:     user.ID,
		Role:    input.Role,
	}

	mid, err = gs.memberRepository.InsertOne(member)
	if err != nil {
		err = utils.NewError(500, err.Error())
		return
	}
	return
}

func (gs *GroupService) DeleteMember(myId, gid, uid int64) (ok bool, err error) {
	isAllowed, err := gs.memberRepository.IsAllowed(myId, gid)
	if !isAllowed {
		if err != nil {
			err = utils.NewError(500, err.Error())
		}else{
			err = utils.NewError(400, "操作权限不够")
		}
		return
	}

	ok, err = gs.memberRepository.DeleteOne(gid, uid)
	if err != nil {
		if err.Error() == "不存在的成员" {
			err = utils.NewError(400, err.Error())
		} else {
			err = utils.NewError(500, err.Error())
		}
		return
	}
	return
}

func (gs *GroupService) UpdateMember(myId int64, member *models.Member) (mid int64, err error) {
	isAllowed, err := gs.memberRepository.IsAllowed(myId, member.Gid)
	if !isAllowed {
		if err != nil {
			err = utils.NewError(500, err.Error())
		}else{
			err = utils.NewError(400, "操作权限不够")
		}
		return
	}

	mid, err = gs.memberRepository.UpdateOne(member)
	if err != nil {
		err = utils.NewError(500, err.Error())
		return
	}
	return
}

func (gs *GroupService) SelectAllMembers(gid int64) (userMembers []*models.UserMember, err error) {
	userMembers, err = gs.memberRepository.SelectAllByGid(gid)
	if err != nil {
		if err.Error() == "没有查询到有效信息" {
			err = utils.NewError(400, err.Error())
		} else {
			err = utils.NewError(500, err.Error())
		}
		return
	}
	return
}
