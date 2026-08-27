package app

import (
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/api"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/admin"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/contabilita"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/login"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/magazzino"
	style "github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/styles"
)

type paginaAttiva int

const (
	LOGIN paginaAttiva = iota
	MAGMENU
	CONTMENU
	ADMINMENU
)

type userRole int

const (
	MAGAZZINIERE userRole = iota
	CONTABILE
	ADMIN
)

type mainLayout struct {
	apiClient *api.APIClient
	UserRole  userRole

	currentPage paginaAttiva
	WindowSize  style.WindowSize

	Login login.Login

	MagMenu   magazzino.MagLayout
	ContMenu  contabilita.ContLayout
	AdminMenu admin.AdminLayout
}

func New(client *api.APIClient) mainLayout {
	loginModel := login.New()
	loginModel.Client = client

	return mainLayout{
		currentPage: LOGIN,
		WindowSize:  style.NewWindowSize(0, 0),

		Login:     loginModel,
		apiClient: client,
	}
}
