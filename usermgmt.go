package main

import (
	"models"
	"utils/crypto/bcrypt"
	"time"
	"fmt"
)

const (
	MAX_NUM_SESSIONS = 100
	SESSION_TIMEOUT = 300
)

type UserData struct {
	userName        string
	sessionId       uint32
	sessionTimer   *time.Timer
}

func (mgr *ConfigMgr)CreateDefaultUser() (status bool) {
	var found bool
	var user models.UserConfig
	defaultPassword := []byte("admin123")
	rows, err := mgr.dbHdl.Query("select * from UserConfig where UserName=?", "admin")
	if err != nil {
		logger.Println("ERROR: Error in reaing UserConfig table ", err)
		return false
	}
	defer rows.Close()
	for rows.Next() {
		if found {
			logger.Println("ERROR: more than  one admin present in UserConfig table ", err)
			return false
		}
		err = rows.Scan(&user.UserName, &user.Password, &user.Description, &user.Previledge)
		if err == nil {
			found = true
		}
	}
	if found == false {
		logger.Println("Creating default user")
		hashedPassword, err := bcrypt.GenerateFromPassword(defaultPassword, bcrypt.DefaultCost)
		user.UserName = "admin"
		user.Password = string(hashedPassword)
		user.Description = "administrator"
		user.Previledge = "w"
		if err != nil {
			logger.Println("Failed to encrypt password for ", user.UserName)
		}
		_, _ = user.StoreObjectInDb(mgr.dbHdl)
	}
	return true
}

func (mgr *ConfigMgr)ReadConfiguredUsersFromDb() (status bool) {
	var userConfig models.UserConfig
	var userData UserData
	dbCmd := "select * from UserConfig"
	rows, err := mgr.dbHdl.Query(dbCmd)
	if err != nil {
		fmt.Println(fmt.Sprintf("DB method Query failed for 'UserConfig' with error UserConfig", dbCmd, err))
		return false
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&userConfig.UserName, &userConfig.Password, &userConfig.Description, &userConfig.Previledge); err != nil {
			fmt.Println("Db Scan failed when interating over UserConfig")
		}
		userData.userName = userConfig.UserName
		userData.sessionId = 0
		mgr.users = append(mgr.users, userData)
	}
	return true
}

func (mgr *ConfigMgr)CreateUser(userName string) (status bool) {
	var userData UserData
	_, _, found := mgr.GetUserByUserName(userName)
	if found {
		logger.Printf("User %s already exists\n", userName)
		return false
	}
	userData.userName = userName
	userData.sessionId = 0
	mgr.users = append(mgr.users, userData)
	return true
}

func (mgr *ConfigMgr)DeleteUser(userName string) (status bool) {
	var userData UserData
	_, idx, found := mgr.GetUserByUserName(userName)
	if found == false {
		logger.Printf("User %s does not exists\n", userName)
		return false
	}
	userData.userName = userName
	userData.sessionId = 0
	mgr.users = append(mgr.users[:idx], mgr.users[idx+1:]...)
	return true
}

func (mgr *ConfigMgr)StartUserSessionHandler() (status bool) {
	logger.Println("Starting SessionHandler thread")
	for {
		sessionId := <-mgr.sessionChan
		user, _, found := mgr.GetUserBySessionId(sessionId)
		if found {
			go user.StartSessionTimer()
		}
	}
	return true
}

func (mgr *ConfigMgr)GetUserBySessionId(sessionId uint32) (user UserData, idx int, found bool) {
	for i, user := range mgr.users {
		if user.sessionId == sessionId {
			return user, i, true
		}
	}
	return user, 0, false
}

func (mgr *ConfigMgr)GetUserByUserName(userName string) (user UserData, idx int, found bool) {
	for i, user := range mgr.users {
		if user.userName == userName {
			return user, i, true
		}
	}
	return user, 0, false
}

func (user UserData)StartSessionTimer() (err error) {
	user.sessionTimer = time.NewTimer(time.Second * SESSION_TIMEOUT)
	<-user.sessionTimer.C
	logger.Printf("Session timeout for user %s session %d\n", user.userName, user.sessionId)
	LogoutUser(user.userName, user.sessionId)
	return nil
}

func LoginUser(userName, password string) (sessionId uint32, status bool) {
	var found bool
	var user models.UserConfig
	rows, err := gMgr.dbHdl.Query("select * from UserConfig where UserName=?", userName)
	if err != nil {
		logger.Println("ERROR: Error in reaing UserConfig table ", err)
		return 0, false
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&user.UserName, &user.Password, &user.Description, &user.Previledge)
		if err == nil {
			found = true
			break
		}
	}
	if found {
		// Comparing the password with the hash
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			logger.Println("Password didn't match for ", userName, err)
			return 0, false
		} else {
			gMgr.sessionId += 1
			userData, _, found := gMgr.GetUserByUserName(userName)
			if found {
				userData.sessionId = gMgr.sessionId
				fmt.Printf("Password matched for %s: sessionId is %d\n", userData.userName, userData.sessionId)
				gMgr.sessionChan <-userData.sessionId
				return userData.sessionId, true
			} else {
				logger.Println("Didn't find user in configmgr's users table")
			}
		}
	}
	return 0, false
}

func LogoutUser(userName string, sessionId uint32) (status bool) {
	user, _, found := gMgr.GetUserByUserName(userName)
	if found {
		if user.sessionId != sessionId {
			logger.Println("Logout: Failed due to session handle mismatch - ", sessionId, user.sessionId)
			return false
		}
		user.sessionTimer.Stop()
		user.sessionId = 0
	}
	return true
}

func AuthenticateSessionId(sessionId uint32) (status bool) {
	user, _, found := gMgr.GetUserBySessionId(sessionId)
	if found {
		user.sessionTimer.Reset(time.Second * SESSION_TIMEOUT)
		return true
	}
	return false
}
