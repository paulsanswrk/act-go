package web

import "ACT_GO/db"

func auth(username string, pwd string) bool {
	user := db.GetUser(username, pwd)
	return user != nil
}
