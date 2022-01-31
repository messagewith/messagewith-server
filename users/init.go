package users

import (
	"github.com/kamva/mgm/v3"
	database "messagewith-server/users/database"
)

var (
	Service                 *service
	collection              *mgm.Collection
	resetPasswordCollection *mgm.Collection
)

func InitService() {
	collection = database.GetDB().UseCollection()
	resetPasswordCollection = database.GetResetPasswordDB().UseCollection()
	Service = GetService(&Repository{}, &ResetPasswordRepository{})
}
