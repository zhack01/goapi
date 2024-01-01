package api

type ServiceAPI struct {
}
type ResponseType interface {
	UserData
}
type Response[T ResponseType] struct {
	Status Status `json:"status"`
	Data   *T     `json:"data"`
}
type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
type UserData struct {
	Token      string `json:"token"`
	Expires_at string `json:"exp"`
}
