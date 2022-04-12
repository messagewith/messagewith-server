package chats

import (
	database "messagewith-server/chats/database"

	"github.com/kamva/mgm/v3"
)

var (
	Service    *service
	collection *mgm.Collection
)

func InitService() {
	Service = getService(&Repository{})
	collection = database.GetDB().UseCollection()
}
