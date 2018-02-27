package view

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"simpledex/model"
	"simpledex/utils"
	"time"
)

// Start function
func Start(m *model.Model, listener net.Listener) {
	server := &http.Server{
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 16}

	http.Handle("/signup/", signupHandler(m))
	http.Handle("/confirmation/", confirmationHandler(m))


	go server.Serve(listener)
}

type signup struct {
	Email    string
	Password string
	ReferredBy string
}


type successResponse struct {
	Email    string
	Is_confirmed bool
	Created_at time.Time
}


func confirmationHandler(m *model.Model) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		keys, ok := r.URL.Query()["activationHash"]
		if !ok || len(keys) < 1 {
			log.Println("Url Param 'activationHash' is missing")
			http.NotFound(w, r)
			return
		} else {
			user, err := m.GetUserByActivationLink(keys[0])
			if err != nil || user == nil  {
				http.NotFound(w, r)
				return
			} else {
				if user != nil {
					_, err = m.ActivateUser(user.Id)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					} else {
						w.Header().Set("Content-Type", "application/json")
						jsonData := []byte(`{"result": "User activated"}`)
						w.Write(jsonData)
						return
						}
					}

			}
		}

	})
}


func signupHandler(m *model.Model) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request(): %v", r)

		if r.Method == "POST" {
			decoder := json.NewDecoder(r.Body)

			var sup signup
			err := decoder.Decode(&sup)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				jsonData := []byte(`{"result": "Not valid data"}`)
				w.Write(jsonData)
				return
			}
			validationErr := utils.ValidatEmail(sup.Email)
			if validationErr != nil {
				w.Header().Set("Content-Type", "application/json")
				jsonData := []byte(`{"result": "Not valid data"}`)
				w.Write(jsonData)
				return
			}
			user, err := m.OneUser(sup.Email)
			if err !=nil {
				w.Header().Set("Content-Type", "application/json")
				jsonData := []byte(`{"result": "Not valid data"}`)
				w.Write(jsonData)
				return
			}
			if len(user) == 0 {
				_, err := m.CreateUser(sup.Email, sup.Password, sup.ReferredBy, false)

				if err != nil {
					panic(err)
				} else{
					createdUser, err := m.OneUser(sup.Email)
					if err != nil {
						fmt.Println("Error:", err.Error())
					} else {
						fmt.Println(createdUser)
						response := successResponse{createdUser[0].Email,
													createdUser[0].Is_confirmed,
													createdUser[0 ].Created_at}
						jsonData, err := json.Marshal(response)
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}
						w.Header().Set("Content-Type", "application/json")
						w.Write(jsonData)
						}
				}

			} else{
				w.Header().Set("Content-Type", "application/json")
				jsonData := []byte(`{"result": "User exists"}`)
				w.Write(jsonData)
			}
		} else {
			fmt.Fprintf(w, "Only POST method is supported.")
		}
	})
}
