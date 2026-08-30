package api

import (
	"net/http"
	"net/url"
)

// UserCreateInput rappresenta il payload per la creazione di un utente.
type UserCreateInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Ruolo    string `json:"ruolo"`
}

// UserUpdateInput rappresenta il payload per la modifica di un utente esistente.
type UserUpdateInput struct {
	Username         string `json:"username,omitempty"`
	Ruolo            string `json:"ruolo,omitempty"`
	Bloccato         *bool  `json:"bloccato,omitempty"`
	TentativiFalliti *int   `json:"tentativi_falliti,omitempty"`
	Password         string `json:"password,omitempty"`
}

// ToggleUserBlock modifica lo stato di blocco di un utente nel backend. Se sbloccato, azzera i tentativi falliti.
func (c *APIClient) ToggleUserBlock(username string, bloccato bool) (User, error) {
	tentativi := 0
	return c.UpdateUser(username, UserUpdateInput{
		Bloccato:         &bloccato,
		TentativiFalliti: &tentativi,
	})
}

// GetUsersByUsername recupera dal backend l'utente identificato dallo username.
// Effettua una richiesta GET a /api/v1/users/{username} e trasforma i dati
// server in un oggetto User. Restituisce un errore se la richiesta fallisce.
func (c *APIClient) GetUsersByUsername(username string) (User, error) {
	var dto serverUserDTO
	if err := c.doJSON(
		http.MethodGet,
		"/api/v1/users/"+url.PathEscape(username),
		nil,
		&dto,
	); err != nil {
		return User{}, err
	}

	return User{
		ID:               dto.ID,
		Username:         dto.Username,
		Ruolo:            dto.Ruolo,
		TentativiFalliti: dto.TentativiFalliti,
		Bloccato:         dto.Bloccato,
		CreatedAt:        dto.CreatedAt,
	}, nil
}

// CreateUser crea un nuovo utente nel backend.
// Accetta i dati di creazione via UserCreateInput e effettua una richiesta POST a /api/v1/users.
// Restituisce l'utente appena creato oppure un errore se l'operazione fallisce.
func (c *APIClient) CreateUser(input UserCreateInput) (User, error) {
	var dto serverUserDTO
	if err := c.doJSON(http.MethodPost, "/api/v1/users", input, &dto); err != nil {
		return User{}, err
	}
	return User{
		ID:               dto.ID,
		Username:         dto.Username,
		Ruolo:            dto.Ruolo,
		TentativiFalliti: dto.TentativiFalliti,
		Bloccato:         dto.Bloccato,
		CreatedAt:        dto.CreatedAt,
	}, nil
}

// UpdateUser modifica username e ruolo di un utente esistente nel backend.
// Accetta lo username dell'utente da modificare e i nuovi valori via UserUpdateInput.
// Effettua una richiesta PUT a /api/v1/users/{username} e restituisce l'utente aggiornato oppure un errore.
func (c *APIClient) UpdateUser(username string, input UserUpdateInput) (User, error) {
	var dto serverUserDTO
	if err := c.doJSON(http.MethodPut, "/api/v1/users/"+url.PathEscape(username), input, &dto); err != nil {
		return User{}, err
	}
	return User{
		ID:               dto.ID,
		Username:         dto.Username,
		Ruolo:            dto.Ruolo,
		TentativiFalliti: dto.TentativiFalliti,
		Bloccato:         dto.Bloccato,
		CreatedAt:        dto.CreatedAt,
	}, nil
}

// DeleteUser elimina un utente dal backend.
// Effettua una richiesta DELETE a /api/v1/users/{username} e restituisce un errore se l'operazione fallisce.
func (c *APIClient) DeleteUser(username string) error {
	return c.doJSON(http.MethodDelete, "/api/v1/users/"+url.PathEscape(username), nil, nil)
}
