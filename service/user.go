package service

import (
	"fmt"
	"log"
	"errors"

	"portal/database"
	"portal/common"
	"portal/model"

	"golang.org/x/crypto/bcrypt"
)
// 用户登录
// 验证用户信息,生成token
func Signin(email, passwd string) (int, string) {
	var (
		id int
		name string
		password string
		status int
		checkStatus int
	)
	// 根据email查询用户
	rows, err := database.ConnDB().Query(`SELECT 
		id, name, password, status, check_status 
		FROM portal_users 
		WHERE email = ? 
		AND deleted_at = ?`, email, common.DeletedAt)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	// 遍历结果rows
	for rows.Next() {
		err := rows.Scan(&id, &name, &password, &status, &checkStatus)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return id, name
}
/**
 * 用户注册
 * @params 
 *   User struct (post body)
 * @return
 *   code int (0: 成功)
 *   msg  error
 */
func Signup(User common.SignupForm) (int, error)  {
	// 验证邮箱,手机号是否已注册
	if flag, _ := database.FindOneUser(User, "(`email` = ? OR `mobile` = ?) AND `deleted_at` = ?"); flag {
		return 100010, errors.New("该邮箱或是手机号已注册")
	}
	// 密码加密
	hash, _ := bcrypt.GenerateFromPassword([]byte(User.Password), bcrypt.DefaultCost)
	User.Password = string(hash)
	// insert user data to user table
	err := database.AddUser(User)
	if (err != nil) {
		return 1, err
	}
	return 0, nil
}
/**
 * 查询用户列表
 * @params
 *   query: struct (query参数)
 * @reutrn 
 *   resutl [] *model.User
 *   msg    error
 */
func QueryUserList(query common.UserQueryBody) ([]*model.User, error) {
	var (
		where string = "`u`.`status` = ?"
		values []string
	)
	// 用户状态, status default value 1
	if query.Where.Status != "" {
		values = append(values, query.Where.Status)
	} else {
		values = append(values, "1")
	}
	// 邮箱
	if query.Where.Email != "" {
		// where += " AND `u`.`email` LIKE '%'||?||'%'"
		where += " AND `u`.`email` LIKE ?"
		values = append(values, query.Where.Email)
	}
	// 审核状态
	if query.Where.CheckStatus != "" {
		where += " AND `u`.`check_status` = ?"
		values = append(values, query.Where.CheckStatus)
	}
	if query.Limit == 0 {
		query.Limit = 10
	}
	// select offset and limit
	where += " LIMIT ?, ?"
	// slice不能直接传递给interface slice
	params := make([]interface{}, len(values)+2)
	for i, v := range values {
		params[i] = v
	}
	// 加入分页
	params[len(values)] = query.Offset
	params[len(values) + 1] = query.Limit
	// Select user table
	res, err := database.FindAllUser(where, params...)
	if err != nil {
		return nil, err
	}
	return res, nil
}
/**
 * 用户状态变更, status -1 注销, 0 禁用, 1 启用
 * @params 
 *   id:     用户id
 *   status: 状态值
 *   remark: 操作描述
 */
func UpdateUserStatus(id string, status int, remark string) (int, error) {
	err := database.UpdateUserStatus(id, status, remark)
	// return
	if err != nil {
		return 1, err
	}
	return 0, nil
}

func Test()  {
	fmt.Println("yes man")
}