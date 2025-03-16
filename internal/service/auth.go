package service

import (
	"EduCommentSync/internal/auth"
	"fmt"
	"net/http"
)

func (s *Service) authHandler(w http.ResponseWriter, r *http.Request) {
	authURL := auth.GetAuthURL()
	http.Redirect(w, r, authURL, http.StatusFound)
}

func (s *Service) oauth2CallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	client, err := auth.GetClient(code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get client: %v", err), http.StatusInternalServerError)
		return
	}

	s.clientMutex.Lock()
	s.clientSync = client
	s.clientMutex.Unlock()

	fmt.Fprintln(w, "Authentication successful! You can now use the application.")
}
