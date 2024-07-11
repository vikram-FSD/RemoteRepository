package language

func Message(lang string, key string) string {
	return message[lang][key]
}
