package routes

import (
	"github.com/anti-raid/evil-befall/pkg/router"
	"github.com/anti-raid/evil-befall/pkg/routes/login"
	"github.com/anti-raid/evil-befall/pkg/routes/root"
)

func init() {
	router.AddRoute(&root.RootRoute{})
	router.AddRoute(&login.LoginRoute{})
}
