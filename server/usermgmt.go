//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//	 Unless required by applicable law or agreed to in writing, software
//	 distributed under the License is distributed on an "AS IS" BASIS,
//	 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	 See the License for the specific language governing permissions and
//	 limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __  
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  | 
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  | 
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   | 
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  | 
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__| 
//                                                                                                           

package server

import (
	"fmt"
	//"models"
	"time"
	//"utils/crypto/bcrypt"
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
	/*
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
				_ = user.StoreObjectInDb(mgr.dbHdl)
			}*/
	return true
}

func (mgr *ConfigMgr) ReadConfiguredUsersFromDb() (status bool) {
	/*
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
	*/
	return true
}

func (mgr *ConfigMgr) CreateUser(userName string) (status bool) {
	var userData UserData
	_, _, found := mgr.GetUserByUserName(userName)
	if found {
		gConfigMgr.logger.Err(fmt.Sprintln("User %s already exists\n", userName))
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
		gConfigMgr.logger.Err(fmt.Sprintln("User %s does not exists\n", userName))
		return false
	}
	userData.userName = userName
	userData.sessionId = 0
	mgr.users = append(mgr.users[:idx], mgr.users[idx+1:]...)
	return true
}

func (mgr *ConfigMgr) StartUserSessionHandler() (status bool) {
	gConfigMgr.logger.Debug("Starting SessionHandler thread")
	mgr.sessionChan = make(chan uint64, MAX_NUM_SESSIONS)
	for {
		select {
		case sessionId := <-mgr.sessionChan:
			gConfigMgr.logger.Info(fmt.Sprintln("SessionHandler - received sessionid ", sessionId))
			_, idx, found := mgr.GetUserBySessionId(sessionId)
			if found {
				gConfigMgr.logger.Debug(fmt.Sprintln("SessionHandler - starting session timer for ", sessionId))
				gConfigMgr.users[idx].sessionTimer = time.NewTimer(time.Second * SESSION_TIMEOUT)
				go gConfigMgr.users[idx].WaitOnSessionTimer()
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
	gConfigMgr.logger.Debug(fmt.Sprintln("Session timeout for user %s session %d\n", user.userName, user.sessionId))
	LogoutUser(user.userName, user.sessionId)
	return nil
}

func LoginUser(userName, password string) (sessionId uint64, status bool) {
	/*
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
		}*/
	return 0, false
}

func LogoutUser(userName string, sessionId uint64) (status bool) {
	user, idx, found := gConfigMgr.GetUserByUserName(userName)
	if found {
		if user.sessionId != sessionId {
			gConfigMgr.logger.Err(fmt.Sprintln("Logout: Failed due to session handle mismatch - ", sessionId, user.sessionId))
			return false
		}
		gConfigMgr.users[idx].sessionTimer.Stop()
		gConfigMgr.users[idx].sessionId = 0
	}
	return true
}

func AuthenticateSessionId(sessionId uint64) (status bool) {
	_, idx, found := gConfigMgr.GetUserBySessionId(sessionId)
	if found {
		gConfigMgr.users[idx].sessionTimer.Reset(time.Second * SESSION_TIMEOUT)
		return true
	}
	return false
}
