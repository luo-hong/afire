package business

func IsAdmin(character []string) bool {
	return len(character) > 0 && character[0] == "1"
}
