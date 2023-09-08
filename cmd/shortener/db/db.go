package db

var db map[string]string = make(map[string]string)

func GetDB(key string) string {
	return db[key]
}

func SetDB(key, val string) {
	db[key] = val
}
