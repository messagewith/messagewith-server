package auth

var (
	Service *service
)

func InitService() {
	Service = &service{}
}
