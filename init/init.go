package init

import (
	"Gomeisa/internal/migrations"
	"Gomeisa/internal/usession"
	"Gomeisa/pkg/utils"
	"encoding/gob"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"math/rand"
	"os"
	"time"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("No .env file found", err)
	}

	file, err := os.OpenFile("gomeisa.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", err)
	}
	utils.Logger = log.New(file, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)

	err = utils.ConnectDB()
	if err != nil {
		log.Fatalln("Failed to open connect DB.\n", err)
	}

	migrations.Up()

	rand.Seed(time.Now().UnixNano())
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)

	usession.Store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	usession.Store.Options = &sessions.Options{
		MaxAge:   60 * 15,
		HttpOnly: true,
	}

	gob.Register(usession.UserSession{})
}