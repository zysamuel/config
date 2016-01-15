package main

import (
	"models"
	"utils/crypto/bcrypt"
	//"fmt"
)

func (mgr *ConfigMgr)CreateDefaultUser() bool {
	var found bool
	var user models.UserConfig
	defaultPassword := []byte("admin123")
	rows, err := mgr.dbHdl.Query("select * from UserConfig where UserName=?", "admin")
	if err != nil {
		logger.Println("ERROR: Error in reaing UserConfig table ", err)
		return false
	}
	for rows.Next() {
		if found {
			logger.Println("ERROR: more than  one admin present in UserConfig table ", err)
			return false
		}
		err = rows.Scan(&user.UserName, &user.Password, &user.Description, &user.Previledge)
		if err == nil {
			found = true
			logger.Println("Found admin user: ", user)
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
			logger.Println("Failed to encrypt password for user ", user)
		}
		// Comparing the password with the hash
		//err = bcrypt.CompareHashAndPassword(user.Password, defaultPassword)
		//if err != nil {
		//	fmt.Println("Password didn't match ", err)
		//} else {
		//	fmt.Println("Password matched ")
		//}
		_, _ = user.StoreObjectInDb(mgr.dbHdl)
	}
	return true
}
