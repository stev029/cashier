package autoload

import "github.com/stev029/cashier/etc/database"

func init() {
	database.InitRedis()
	database.DBConnect()
}
