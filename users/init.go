package users

var (
	Service *service
)

func InitUserService() {
	Service = getService()
}
