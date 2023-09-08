package db

func ShortToLong(shortLink string) (string, bool) {
	var link string
	for key, value := range db {
		if value == shortLink {
			link = key
			break
		}
	}
	if link == "" {
		return "", true
	}

	return link, false
}
