package routes

import (
	"github.com/anti-raid/evil-befall/pkg/router"
	"github.com/anti-raid/evil-befall/pkg/routes/login"
	"github.com/anti-raid/evil-befall/pkg/routes/showstate"
)

func init() {
	router.AddRoute(&login.LoginRoute{})
	router.AddRoute(&showstate.ShowStateRoute{})
}
