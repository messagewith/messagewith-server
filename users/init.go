package users

var (
	Service *service
)

func InitService() {
	Service = getService()
}
