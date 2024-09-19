package routes

import (
	"github.com/anti-raid/evil-befall/pkg/router"
	"github.com/anti-raid/evil-befall/pkg/routes/login"
)

func init() {
	router.AddRoute(&login.LoginRoute{})
}
