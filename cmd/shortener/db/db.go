package db

var db map[string]string = make(map[string]string)

func GetDb(key string) string {
	return db[key]
}

func SetDb(key, val string) {
	db[key] = val
}
