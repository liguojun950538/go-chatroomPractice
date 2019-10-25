package model

import (
	"chatroom/common"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/garyburd/redigo/redis"
)

var (
	TheUserDao *UserDao
)

// 数据接入层, 控制层通过数据接入层和数据库打交道

type UserDao struct {
	pool *redis.Pool
}

func NewUserDao(pool *redis.Pool) (ud *UserDao) {
	ud = &UserDao{
		pool: pool,
	}
	return
}

// 通过id, 从redis获得一个user实例
func (this *UserDao) getUser(conn redis.Conn, id int) (user *common.User, err error) {
	result, err := redis.String(conn.Do("HGet", "users", fmt.Sprintf("%d", id)))
	if err != nil {
		err = ErrUserNotExist
		user = nil
		return
	}

	user = &common.User{}
	err = json.Unmarshal([]byte(result), user)
	if err != nil {
		user = nil
		err = errors.New("json解析失败, 可能是注册时添加信息的格式不对, 请检查")
		return
	}
	err = nil
	return
}

// 根据用户名和密码在数据库验证是否正确
func (this *UserDao) ConfirmLogin(id int, passwd string) (user *common.User, err error) {
	conn := this.pool.Get()
	defer conn.Close()

	user, err = this.getUser(conn, id)
	if err != nil {
		return nil, err
	}

	if passwd != user.Passwd {
		err = ErrInvalidPasswd
		return nil, err
	}
	err = nil
	return
}

// 完成注册
func (this *UserDao) Register(user *common.User) (err error) {
	conn := this.pool.Get()
	defer conn.Close()

	// 如果用户已经存在, 则为无效的注册
	_, err = this.getUser(conn, user.UserID)
	if err == nil {
		err = ErrUserExist
		return
	}

	userByte, err := json.Marshal(user)
	if err != nil {
		return
	}

	_, err = conn.Do("HSet", "users", fmt.Sprintf("%d", user.UserID), string(userByte))
	if err != nil {
		err = errors.New("注册到数据库失败, 请检查内存, 或者数据库服务端是否打开")
		return
	}
	return nil
}
