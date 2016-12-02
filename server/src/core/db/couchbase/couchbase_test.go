// to do fix
package couchbase

func testSync(db *cb.CbEngineEx) {
	key := "testSync"

	if added, err := db.InsertSync(key, []byte("testVal")); true {
		log.Println("insert:", added, err)
	}

	if val, _ := db.GetSync(key); true {
		log.Println(string(val))
	}

	if err := db.SetSync(key, []byte("testVal2")); true {
		log.Println("set:", err)
	}

	if val, _ := db.GetSync(key); true {
		log.Println(string(val))
	}
}

func testAsyn(db *cb.CbEngineEx) {
	key2 := "testAsyn"
	val2 := []byte("hello world")
	if op, _ := db.InsertAsyn(key2, val2, nil, true); true {
		res := <-op.ChRes
		log.Println("asyn insert:", res.Ok(), res.Error())
	}

	if op, _ := db.GetAsyn(key2, nil); true {
		res := <-op.ChRes
		log.Println("asyn get:", string(res.Res), res.Err)
	}
}

func testDBSet(data *cb.OpSetData, res *cb.ResSet) {
	log.Println("main asyn set:", res.Err)
}

func testDBGet(data *cb.OpGetData, res *cb.ResGet) {
	log.Println("main asyn get:", string(res.Res), res.Err)
}

func testAsynEx(db *cb.CbEngineEx) {
	key, val := "main test asyn", []byte("main proc")

	db.SetAsynExFunc(key, val, nil, testDBSet)
	db.GetAsynExFunc(key, nil, testDBGet)

	time.Sleep(time.Second)
	db.ProcAsynRes()
}

// 测试db
func testDB() {
	//db := cb.NewCbEngineEx()
	if err := db.Serve(&cb.ParamInit{"127.0.0.1:8091", "AOW", "111111", 100}); err != nil {
		log.Println(err)
		return
	}

	testSync(db.DB())
	testAsyn(db.DB())
	testAsynEx(db.DB())

	db.DB().Stop()
}
