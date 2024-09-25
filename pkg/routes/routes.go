package routes

import (
	"github.com/anti-raid/evil-befall/pkg/router"
	"github.com/anti-raid/evil-befall/pkg/routes/apiexec_exec"
	"github.com/anti-raid/evil-befall/pkg/routes/apiexec_ls"
	"github.com/anti-raid/evil-befall/pkg/routes/choose_guild"
	"github.com/anti-raid/evil-befall/pkg/routes/login"
	"github.com/anti-raid/evil-befall/pkg/routes/showstate"
)

func init() {
	router.AddRoute(&apiexec_ls.ApiExecLsRoute{})
	router.AddRoute(&apiexec_exec.ApiExecExecRoute{})
	router.AddRoute(&choose_guild.ChooseGuildRoute{})
	router.AddRoute(&login.LoginRoute{})
	router.AddRoute(&showstate.ShowStateRoute{})
}
