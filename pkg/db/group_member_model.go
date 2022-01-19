package db

import (
	"errors"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/utils"
)

func (d *DataBase) GetGroupMemberInfoByGroupIDUserID(groupID, userID string) (*LocalGroupMember, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var groupMember LocalGroupMember
	return &groupMember, utils.Wrap(d.conn.Where("group_id = ? AND user_id = ?",
		groupID, userID).Error, "GetGroupMemberInfoByGroupIDUserID failed")
}

func (d *DataBase) GetAllGroupMemberList() ([]LocalGroupMember, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var groupMemberList []LocalGroupMember
	return groupMemberList, utils.Wrap(d.conn.Find(&groupMemberList).Error, "GetAllGroupMemberList failed")
}

func (d *DataBase) GetGroupSomeMemberInfo(groupID string, userIDList []string) ([]*LocalGroupMember, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var groupMemberList []LocalGroupMember
	err := d.conn.Where("group_id = ? And user_id IN ? ", groupID, userIDList).Find(&groupMemberList).Error
	var transfer []*LocalGroupMember
	for _, v := range groupMemberList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, "GetGroupMemberListByGroupID failed ")
}

func (d *DataBase) GetGroupMemberListByGroupID(groupID string) ([]*LocalGroupMember, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var groupMemberList []LocalGroupMember
	err := d.conn.Where("group_id = ? ", groupID).Find(&groupMemberList).Error
	var transfer []*LocalGroupMember
	for _, v := range groupMemberList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, "GetGroupMemberListByGroupID failed ")
}

func (d *DataBase) GetGroupOwnerAndAdminByGroupID(groupID string) ([]*LocalGroupMember, error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	var groupMemberList []LocalGroupMember
	err := d.conn.Where("group_id = ?  AND role_level > ?", groupID, constant.GroupOrdinaryUsers).Find(&groupMemberList).Error
	var transfer []*LocalGroupMember
	for _, v := range groupMemberList {
		v1 := v
		transfer = append(transfer, &v1)
	}
	return transfer, utils.Wrap(err, "GetGroupMemberListByGroupID failed ")
}

func (d *DataBase) GetGroupMemberUIDListByGroupID(groupID string) (result []string, err error) {
	d.mRWMutex.RLock()
	defer d.mRWMutex.RUnlock()
	err = d.conn.Where("group_id = ? ", groupID).Pluck("user_id", &result).Error
	return result, utils.Wrap(err, "GetGroupMemberListByGroupID failed ")
}

func (d *DataBase) InsertGroupMember(groupMember *LocalGroupMember) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return d.conn.Create(groupMember).Error
}

func (d *DataBase) DeleteGroupMember(groupID, userID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	groupMember := LocalGroupMember{GroupID: groupID, UserID: userID}
	return d.conn.Where("group_id=? and user_id=?", groupID, userID).Delete(&groupMember).Error
}

func (d *DataBase) UpdateGroupMember(groupMember *LocalGroupMember) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.Updates(groupMember)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return t.Error
}

func (d *DataBase) GetGroupMenberInfoIfOwnerOrAdmin() ([]*LocalGroupMember, error) {
	var ownerAndAdminList []*LocalGroupMember
	groupList, err := d.GetJoinedGroupList()
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	for _, v := range groupList {

		memberList, err := d.GetGroupOwnerAndAdminByGroupID(v.GroupID)
		if err != nil {
			return nil, utils.Wrap(err, "")
		}
		ownerAndAdminList = append(ownerAndAdminList, memberList...)
	}
	return ownerAndAdminList, nil
}
