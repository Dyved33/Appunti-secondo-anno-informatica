package app

import (
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/admin"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/contabilita"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/login"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/pages/magazzino"
	style "github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/styles"
	"github.com/tommasogarofani/ingsw_25-26_group-11/client/ui/tabnav"

	tea "charm.land/bubbletea/v2"
)

func (f mainLayout) Init() tea.Cmd {
	// Inizializza anche il sottomodello di login
	return f.Login.Init()
}

func (f mainLayout) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyPressMsg:
		// Intercetta la shortcut solo se l'utente è attualmente autenticato
		if f.currentPage != LOGIN && (msg.String() == "ctrl+l") {
			// 1. Torna alla schermata di Login
			f.currentPage = LOGIN

			// 2. Resetta il sotto-modello di Login per pulire i campi username/password
			f.Login = login.New()
			f.Login.Client = f.apiClient

			// 3. (Opzionale) Reset della sessione/token se la struct APIClient lo prevede
			// f.apiClient.ClearToken()

			// 4. Riavvia Init() del Login per rimettere il focus sul primo campo di input
			return f, f.Login.Init()
		}

	case tea.WindowSizeMsg:
		f.WindowSize.ChangeSize(msg.Width, msg.Height)
		f.WindowSize.Border.SetSpacing(2, 2, 2, 2)

		// Propaga il resize a tutte le sezioni top-level per mantenere i layout sincronizzati.
		loginUpdated, _ := f.Login.Update(msg)
		f.Login = loginUpdated.(login.Login)
		f.MagMenu.DimUpdate(msg)
		f.ContMenu.DimUpdate(msg)
		f.AdminMenu.DimUpdate(msg)

		// Inoltra il resize anche alla sottopagina attiva: senza questo passaggio
		// i modelli di pagina (es. scorrimento liste ordini) non vengono mai a
		// conoscenza dell'altezza reale del terminale.
		switch f.currentPage {
		case MAGMENU:
			updatedMag, _ := f.MagMenu.Update(msg)
			f.MagMenu = updatedMag.(magazzino.MagLayout)
		case CONTMENU:
			updatedCont, _ := f.ContMenu.Update(msg)
			f.ContMenu = updatedCont.(contabilita.ContLayout)
		case ADMINMENU:
			updatedAdmin, _ := f.AdminMenu.Update(msg)
			f.AdminMenu = updatedAdmin.(admin.AdminLayout)
		}
		return f, nil

	// Login completato con successo
	// Viene cambiata la schermata attiva in base al ruolo dell'utente
	case login.LoginSubmitMsg:
		f.UserRole = parseUserRole(f.apiClient.Role())

		// I layout nascono con le dimensioni a zero: vanno inizializzate
		// SUBITO con quelle correnti del terminale (arrivate prima del login),
		// altrimenti le hitbox dei tab cliccabili vengono misurate su un
		// terminale fittizio 80x24 e i click cadono nel vuoto fino al prossimo resize.
		dims := tea.WindowSizeMsg{
			Width:  style.ReadWidth(f.WindowSize),
			Height: style.ReadHeight(f.WindowSize),
		}

		switch f.UserRole {
		case MAGAZZINIERE:
			f.currentPage = MAGMENU
			f.MagMenu = magazzino.New(f.apiClient)
			f.MagMenu.DimUpdate(dims)
			return f, f.MagMenu.Init()

		case CONTABILE:
			f.currentPage = CONTMENU
			f.ContMenu = contabilita.New(f.apiClient)
			f.ContMenu.DimUpdate(dims)
			return f, f.ContMenu.Init()
		case ADMIN:
			f.currentPage = ADMINMENU
			f.AdminMenu = admin.New(f.apiClient)
			f.AdminMenu.DimUpdate(dims)
			return f, f.AdminMenu.Init()
		}
		return f, nil
	}

	// Aggiorna la schermata attiva in base al messaggio ricevuto
	switch f.currentPage {
	case LOGIN:
		nuovoLogin, cmd := f.Login.Update(msg)
		f.Login = nuovoLogin.(login.Login)
		return f, cmd

	case MAGMENU:
		nuovoMag, cmd := f.MagMenu.Update(msg)
		f.MagMenu = nuovoMag.(magazzino.MagLayout)
		return f, cmd

	case CONTMENU:
		nuovoCont, cmd := f.ContMenu.Update(msg)
		f.ContMenu = nuovoCont.(contabilita.ContLayout)
		return f, cmd

	case ADMINMENU:
		nuovoAdmin, cmd := f.AdminMenu.Update(msg)
		f.AdminMenu = nuovoAdmin.(admin.AdminLayout)
		return f, cmd
	}

	return f, nil
}

func (f mainLayout) View() tea.View {
	// Contenuto pagina attiva
	var contentBlock string
	var cursor *tea.Cursor
	switch f.currentPage {
	case LOGIN:
		v := f.Login.View()
		contentBlock = v.Content
		cursor = v.Cursor
	case MAGMENU:
		v := f.MagMenu.View()
		contentBlock = v.Content
		cursor = v.Cursor
	case CONTMENU:
		v := f.ContMenu.View()
		contentBlock = v.Content
		cursor = v.Cursor
	case ADMINMENU:
		v := f.AdminMenu.View()
		contentBlock = v.Content
		cursor = v.Cursor
	}

	// Layout finale: la composizione è definita una sola volta in tabnav
	// (la stessa usata per calcolare le aree cliccabili dei tab).
	w := style.ReadWidth(f.WindowSize)
	h := style.ReadHeight(f.WindowSize)
	finalPage := tabnav.PlacePage(w, h, contentBlock)

	altscreen := tea.NewView(finalPage)
	altscreen.AltScreen = true
	altscreen.Cursor = cursor
	// Abilita gli eventi mouse (necessari per i tab cliccabili dei layout)
	altscreen.MouseMode = tea.MouseModeCellMotion

	return altscreen
}

func parseUserRole(role string) userRole {
	switch role {
	case "magazziniere":
		return MAGAZZINIERE
	case "contabile":
		return CONTABILE
	case "admin":
		return ADMIN
	default:
		return MAGAZZINIERE
	}
}
