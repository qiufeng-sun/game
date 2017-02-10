// 游戏数据库
package db

//
import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"

	"util/dbs/database"
	"util/dbs/redis"
	"util/logs"
)

////////////////////////////////////////////////////////////////////////////
//
func Init(path string) {
	// redis
	InitRedis(path)

	// mysql
	InitMysql(path)
}

//
func HealthCheck() error {
	// redis
	if e := HealthCheckRedis(); e != nil {
		return e
	}

	// mysql
	return HealthCheckMysql()
}

////////////////////////////////////////////////////////////////////////////
//
var g_cachePools *redis.RedisPools

func InitRedis(path string) {
	fileName := path + "redis.json"
	redis.InitByFile(fileName)

	g_cachePools = redis.GetRedisPools("cache")

	logs.Info("init redis ok!")
}

//
func HealthCheckRedis() error {
	return redis.HealthCheck()
}

//
func getCacheConn(uid string) *redis.RedisConn {
	return g_cachePools.GetConn()
}

////////////////////////////////////////////////////////////////////////////
//
func InitMysql(path string) {
	fileName := path + "mysql.json"
	database.InitByFile(fileName)

	logs.Info("init mysql ok!")
}

//
func HealthCheckMysql() error {
	return database.HealthCheck()
}

//
func getPlayerOrm() orm.Ormer {
	return database.GetDefOrm()
}
