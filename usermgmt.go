package main

import (
	"fmt"
	"models"
	"time"
	"utils/crypto/bcrypt"
)

const (
	MAX_NUM_SESSIONS = 10
	SESSION_TIMEOUT  = 300
)

type UserData struct {
	userName      string
	sessionId     uint64
	sessionTimer  *time.Timer
	lastLoginTime time.Time
	lastLoginIp   string
	numAPICalled  uint32
}

func (mgr *ConfigMgr) CreateDefaultUser() (status bool) {
	var found bool
	var user models.User
	defaultPassword := []byte("admin123")
	rows, err := mgr.dbHdl.Query("select * from User where UserName=?", "admin")
	if err != nil {
		logger.Println("ERROR: Error in reaing User table ", err)
		return false
	}
	defer rows.Close()
	for rows.Next() {
		if found {
			logger.Println("ERROR: more than  one admin present in User table ", err)
			return false
		}
		err = rows.Scan(&user.UserName, &user.Password, &user.Description, &user.Privilege)
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
		user.Privilege = "w"
		if err != nil {
			logger.Println("Failed to encrypt password for ", user.UserName)
		}
		if ok := mgr.CreateUser(user.UserName); ok == false {
			logger.Println("Failed to create default user")
		}
		_, _ = user.StoreObjectInDb(mgr.dbHdl)
	}
	return true
}

func (mgr *ConfigMgr) ReadConfiguredUsersFromDb() (status bool) {
	var userConfig models.User
	var userData UserData
	dbCmd := "select * from User"
	rows, err := mgr.dbHdl.Query(dbCmd)
	if err != nil {
		fmt.Println(fmt.Sprintf("DB method Query failed for 'User' with error User", dbCmd, err))
		return false
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&userConfig.UserName, &userConfig.Password, &userConfig.Description, &userConfig.Privilege); err != nil {
			fmt.Println("Db Scan failed when interating over User")
		}
		userData.userName = userConfig.UserName
		userData.sessionId = 0
		mgr.users = append(mgr.users, userData)
	}
	return true
}

func (mgr *ConfigMgr) CreateUser(userName string) (status bool) {
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

func (mgr *ConfigMgr) DeleteUser(userName string) (status bool) {
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

func (mgr *ConfigMgr) StartUserSessionHandler() (status bool) {
	logger.Println("Starting SessionHandler thread")
	mgr.sessionChan = make(chan uint64, MAX_NUM_SESSIONS)
	for {
		select {
		case sessionId := <-mgr.sessionChan:
			logger.Println("SessionHandler - received sessionid ", sessionId)
			_, idx, found := mgr.GetUserBySessionId(sessionId)
			if found {
				logger.Println("SessionHandler - starting session timer for ", sessionId)
				gMgr.users[idx].sessionTimer = time.NewTimer(time.Second * SESSION_TIMEOUT)
				go gMgr.users[idx].WaitOnSessionTimer()
			}
		}
	}
	return true
}

func (mgr *ConfigMgr) GetUserBySessionId(sessionId uint64) (user UserData, idx int, found bool) {
	for i, user := range mgr.users {
		if user.sessionId == sessionId {
			return user, i, true
		}
	}
	return user, 0, false
}

func (mgr *ConfigMgr) GetUserByUserName(userName string) (user UserData, idx int, found bool) {
	for i, user := range mgr.users {
		if user.userName == userName {
			return user, i, true
		}
	}
	return user, 0, false
}

func (user UserData) WaitOnSessionTimer() (err error) {
	<-user.sessionTimer.C
	logger.Printf("Session timeout for user %s session %d\n", user.userName, user.sessionId)
	LogoutUser(user.userName, user.sessionId)
	return nil
}

func LoginUser(userName, password string) (sessionId uint64, status bool) {
	var found bool
	var user models.User
	rows, err := gMgr.dbHdl.Query("select * from User where UserName=?", userName)
	if err != nil {
		logger.Println("ERROR: Error in reaing User table ", err)
		return 0, false
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&user.UserName, &user.Password, &user.Description, &user.Privilege)
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
			userData, idx, found := gMgr.GetUserByUserName(userName)
			if found {
				gMgr.sessionId += 1
				userData.sessionId = gMgr.sessionId
				gMgr.users[idx].sessionId = gMgr.sessionId
				fmt.Printf("Password matched for %s: sessionId is %d\n", userData.userName, userData.sessionId)
				gMgr.sessionChan <- userData.sessionId
				fmt.Printf("SessionId is %d is sent over chan\n", userData.sessionId)
				return userData.sessionId, true
			} else {
				logger.Println("Didn't find user in configmgr's users table")
			}
		}
	}
	return 0, false
}

func LogoutUser(userName string, sessionId uint64) (status bool) {
	user, idx, found := gMgr.GetUserByUserName(userName)
	if found {
		if user.sessionId != sessionId {
			logger.Println("Logout: Failed due to session handle mismatch - ", sessionId, user.sessionId)
			return false
		}
		gMgr.users[idx].sessionTimer.Stop()
		gMgr.users[idx].sessionId = 0
	}
	return true
}

func AuthenticateSessionId(sessionId uint64) (status bool) {
	_, idx, found := gMgr.GetUserBySessionId(sessionId)
	if found {
		gMgr.users[idx].sessionTimer.Reset(time.Second * SESSION_TIMEOUT)
		return true
	}
	return false
}
