// main.go
package main

import (
	"1nce.com/testing-task/client"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)




var clients []client.Client

func returnAllClients(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllClients")
	err := json.NewEncoder(w).Encode(clients)
	if err != nil {
		http.Error(w, "Failed to return clients", http.StatusInternalServerError)
	}
}

func returnSingleClient(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnSingleClient")
	id, err := getIdFromRequest(r)
	if err != nil {
		http.Error(w, "Provided id parameter is not integer", http.StatusBadRequest)
		return
	}

	for _, client := range clients {
		if *client.Id == id {
			error := json.NewEncoder(w).Encode(client)
			if error != nil {
				http.Error(w, "Failed to return client", http.StatusInternalServerError)
			}
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func createNewClient(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createNewClient")

	client, err := generateClientFromRequestBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for _, currClient := range clients {
		// with * we are then comparing pointer values and not pointer addresses
		if *currClient.Id == *client.Id {
			http.Error(w, fmt.Sprintf("Client with same ID %d already exists", *client.Id), http.StatusForbidden)
			return
		}
	}
	err = client.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// update our global Articles array to include
	// our new Article
	clients = append(clients, client)

	w.WriteHeader(http.StatusCreated)
}



func deleteClient(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: deleteClient")
	id, err := getIdFromRequest(r)
	if err != nil {
		http.Error(w, "Provided id parameter is not integer", http.StatusNotFound)
		return
	}
	for index, client := range clients {
		if *client.Id == id {
			clients = append(clients[:index], clients[index+1:]...)
			return
		}
	}
	http.Error(w, fmt.Sprintf("Client with ID %d does not exist", id), http.StatusNotFound)

}

func findClientById(id int) *client.Client {
	for _, currClient := range clients {
		if *currClient.Id == id {
			return &currClient
		}
	}
	return nil
}

func patchClient(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: patchClient")

	id, err := getIdFromRequest(r)
	if err != nil {
		http.Error(w, "Provided id parameter is not integer", http.StatusBadRequest)
		return
	}
	foundClient := findClientById(id)
	if foundClient == nil {
		http.Error(w, fmt.Sprintf("Client with ID %d does not exist", id), http.StatusNotFound)
		return
	}
	client, err := generateClientFromRequestBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// try to merge client instances
	foundClient.Merge(client)
}

func getIdFromRequest(r *http.Request) (id int, err error) {
	vars := mux.Vars(r)
	idString := vars["id"]
	id, err = strconv.Atoi(idString)
	return
}

func generateClientFromRequestBody(body io.ReadCloser) (client client.Client, err error) {
	reqBody, err := ioutil.ReadAll(body)
	if err != nil {
		return
	}
	decoder := json.NewDecoder(bytes.NewReader(reqBody))
	// strict validation to prevent unknown fields
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&client); err != nil {

	}
	return
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	//myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/clients", returnAllClients)
	myRouter.HandleFunc("/client", createNewClient).Methods("POST")
	myRouter.HandleFunc("/client/{id}", deleteClient).Methods("DELETE")
	myRouter.HandleFunc("/client/{id}", patchClient).Methods("PATCH")
	myRouter.HandleFunc("/client/{id}", returnSingleClient)
	log.Fatal(http.ListenAndServe(":8888", myRouter))
}

func main() {
	clients = []client.Client{
	}
	handleRequests()
}
