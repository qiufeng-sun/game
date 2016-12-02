package db

import (
	"log"
    "database/sql"
	
	_ "github.com/go-sql-driver/mysql"
	
	"conf"
)

// mysql错误码
const (
	E_MySql_Success			uint16		= 0					// ok
	E_MySql_Duplicate		uint16		= 1062				// Duplicate entry '%s' for key %d
)

var sqlDB *sql.DB

func Start(dbCfg *conf.DBConf) (err error) {
	log.Println("mysql start!")
	
	sqlDB, err = sql.Open("mysql", GetDSN(dbCfg))
	if nil == err {
		err = sqlDB.Ping()
	}
	
	return err
}

func Close() {
	sqlDB.Close()
}

// user:password@tcp(localhost:5555)/dbname?charset=utf8
func GetDSN(dbCfg *conf.DBConf) string {
	return dbCfg.User + ":" + dbCfg.Pwd + "@tcp(" + 
			dbCfg.IP + ":" + dbCfg.Port + ")/" + dbCfg.DB + "?charset=utf8"
}


/*
    //插入数据
    stmt, err := db.Prepare("INSERT userinfo SET username=?,departname=?,created=?")
    checkErr(err)

    res, err := stmt.Exec("astaxie", "研发部门", "2012-12-09")
    checkErr(err)

    id, err := res.LastInsertId()
    checkErr(err)

    log.Println(id)
    //更新数据
    stmt, err = db.Prepare("update userinfo set username=? where uid=?")
    checkErr(err)

    res, err = stmt.Exec("astaxieupdate", id)
    checkErr(err)

    affect, err := res.RowsAffected()
    checkErr(err)

    log.Println(affect)

    //查询数据
    rows, err := db.Query("SELECT * FROM userinfo")
    checkErr(err)

    for rows.Next() {
        var uid int
        var username string
        var department string
        var created string
        err = rows.Scan(&uid, &username, &department, &created)
        checkErr(err)
        log.Println(uid)
        log.Println(username)
        log.Println(department)
        log.Println(created)
    }

    //删除数据
    stmt, err = db.Prepare("delete from userinfo where uid=?")
    checkErr(err)

    res, err = stmt.Exec(id)
    checkErr(err)

    affect, err = res.RowsAffected()
    checkErr(err)

    log.Println(affect)

    db.Close()
*/