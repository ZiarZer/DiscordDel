package discord

type User struct {
	Id            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	GlobalName    string `json:"global_name"`
	Avatar        string `json:"avatar"`
}

type Guild struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}
