// Role and resource
// 提供角色与资源的权限操作
package service

import (
	"portal/model"
	"portal/database"
)
// Grant permission
func Grant(arg model.RolePrivilege) error {
	return database.BindRoleRes(arg)
}
// View role permission
func GetRolePermisson(roleId int) (*model.ResBelongRole, error) {
	return database.GetRolePermmison(roleId)
}