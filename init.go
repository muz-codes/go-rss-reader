package go_rss_reader

import "log"

func init() {
	_, err := InitZapLogger()
	if err != nil {
		log.Panic(err)
	}
	defer logger.Sync()
}
