package sessions

import (
	"github.com/kamva/mgm/v3"
	database "messagewith-server/sessions/database"
)

var (
	Service    *service
	collection *mgm.Collection
)

func InitService() {
	collection = database.GetDB().UseCollection()
	Service = getService(&Repository{})
}
