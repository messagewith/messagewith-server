package chats

import (
	"github.com/kamva/mgm/v3"
	database "messagewith-server/chats/database"
)

var (
	Service    *service
	collection *mgm.Collection
)

func InitService() {
	Service = getService(&Repository{})
	collection = database.GetDB().UseCollection()
}
