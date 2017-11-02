package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/ferux/mvcPractice/db"

	"github.com/ferux/mvcPractice/model/client"
	"github.com/ferux/mvcPractice/model/item"
	"github.com/ferux/mvcPractice/model/order"

	"github.com/ferux/mvcpractice/model/cart"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

//Config for http server
type Config struct {
	ListenIP   string
	ListenPort string
	ExternalIP string
}

var dbConn *sqlx.DB

//Run a http server
func Run(c Config, dbc *sqlx.DB) {
	var err error
	dbConf := db.New(dbc)
	if err := dbConf.Init(); err != nil {
		log.Println("Got an error while creating DB scheme: ", err)
	}
	dbConn = dbc
	listenString := fmt.Sprintf("%s:%s", c.ListenIP, c.ListenPort)
	listener, err := net.Listen("tcp", listenString)
	if err != nil {
		log.Fatalf("Can't start listening to %s.\nReason:%v", listenString, err)
	}
	r := mux.NewRouter()
	//View data
	appendCreate(r)
	appendDelete(r)
	appendEdit(r)
	appendView(r)
	log.Printf("Launching server at: %s", listenString)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("../view/build/"))))
	headersOk := handlers.AllowedHeaders([]string{"Origin", "Content-Type", "X-Auth-Token", "Accept", "X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{c.ExternalIP})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})
	log.Fatal(http.Serve(listener, handlers.CORS(headersOk, originsOk, methodsOk)(r)))
}

func appendView(r *mux.Router) {
	//Cart
	r.HandleFunc("/api/carts", handleCartsView).Methods("GET")
	r.HandleFunc("/api/cart", handleCartsView).Methods("GET")
	r.HandleFunc("/api/cart/{id}", handleCartView).Methods("GET")
	//Client
	r.HandleFunc("/api/clients", handleClientsView).Methods("GET")
	r.HandleFunc("/api/client", handleClientsView).Methods("GET")
	r.HandleFunc("/api/client/{id}", handleClientView).Methods("GET")
	//Item
	r.HandleFunc("/api/items", handleItemsView).Methods("GET")
	r.HandleFunc("/api/item", handleItemsView).Methods("GET")
	r.HandleFunc("/api/item/{id}", handleItemView).Methods("GET")
	//Order
	r.HandleFunc("/api/orders", handleOrdersView).Methods("GET")
	r.HandleFunc("/api/order", handleOrdersView).Methods("GET")
	r.HandleFunc("/api/order/{id}", handleOrderView).Methods("GET")
	//ItemTypes
	r.HandleFunc("/api/itemtypes", handleItemTypesView).Methods("GET")
}

func appendEdit(r *mux.Router) {
	r.HandleFunc("/api/cart/{id}", handleCartEdit).Methods("POST")
	r.HandleFunc("/api/client/{id}", handleClientEdit).Methods("POST")
	r.HandleFunc("/api/item/{id}", handleItemEdit).Methods("POST")
	r.HandleFunc("/api/order/{id}", handleOrderEdit).Methods("POST")

	r.HandleFunc("/api/clients", handleClientsEdit).Methods("POST")
	// r.HandleFunc("/api/carts", handleCartsEdit).Methods("POST")
	r.HandleFunc("/api/items", handleItemsEdit).Methods("POST")
	r.HandleFunc("/api/orders", handleOrdersEdit).Methods("POST")
}

func appendDelete(r *mux.Router) {
	r.HandleFunc("/api/cart/{id}", handleCartDelete).Methods("DELETE")
	r.HandleFunc("/api/client/{id}", handleClientDelete).Methods("DELETE")
	r.HandleFunc("/api/item/{id}", handleItemDelete).Methods("DELETE")
	r.HandleFunc("/api/order/{id}", handleOrderDelete).Methods("DELETE")

	r.HandleFunc("/api/clients", handleClientsDelete).Methods("DELETE")
	// r.HandleFunc("/api/carts/{id}", handleCartsDelete).Methods("DELETE")
	r.HandleFunc("/api/items", handleItemsDelete).Methods("DELETE")
	r.HandleFunc("/api/orders", handleOrdersDelete).Methods("DELETE")
}

func appendCreate(r *mux.Router) {

	r.HandleFunc("/api/cart", handleCartCreate).Methods("PUT")
	r.HandleFunc("/api/client", handleClientCreate).Methods("PUT")
	r.HandleFunc("/api/item", handleItemCreate).Methods("PUT")
	r.HandleFunc("/api/order", handleOrderCreate).Methods("PUT")

	r.HandleFunc("/api/clients", handleClientsCreate).Methods("PUT")
	// r.HandleFunc("/api/carts", handleCartsCreate).Methods("PUT")
	r.HandleFunc("/api/items", handleItemsCreate).Methods("PUT")
	r.HandleFunc("/api/orders", handleOrdersCreate).Methods("PUT")
}

//Views
func handleCartView(w http.ResponseWriter, r *http.Request) {
	log.Println("[Get] CartView")
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Something wrong with the ID.\n"+err.Error(), http.StatusBadRequest)
		return
	}
	thisRow, err := cart.ReadUUID(dbConn, id)
	if err != nil {
		http.Error(w, "Item's ID is wrong.\n"+err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Conten-type", "application/json")
	json.NewEncoder(w).Encode(thisRow)
}
func handleCartsView(w http.ResponseWriter, r *http.Request) {
	rows, err := cart.Read(dbConn)
	if err != nil {
		http.Error(w, "Something wrong with this.\n"+err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Conten-type", "application/json")
	json.NewEncoder(w).Encode(rows)
}

func handleClientView(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Something wrong with the ID.\n"+err.Error(), http.StatusBadRequest)
		return
	}
	thisRow, err := client.ReadUUID(dbConn, id)
	if err != nil {
		http.Error(w, "Item's ID is wrong.\n"+err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Conten-type", "application/json")
	json.NewEncoder(w).Encode(thisRow)
}
func handleClientsView(w http.ResponseWriter, r *http.Request) {
	rows, err := client.Read(dbConn)
	if err != nil {
		http.Error(w, "Something wrong with this.\n"+err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Conten-type", "application/json")
	json.NewEncoder(w).Encode(rows)
}

func handleItemView(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Something wrong with the ID.\n"+err.Error(), http.StatusBadRequest)
		return
	}
	thisRow, err := item.ReadUUID(dbConn, id)
	if err != nil {
		http.Error(w, "Item's ID is wrong.\n"+err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Conten-type", "application/json")
	json.NewEncoder(w).Encode(thisRow)
}
func handleItemsView(w http.ResponseWriter, r *http.Request) {
	rows, err := item.Read(dbConn)
	if err != nil {
		http.Error(w, "Something wrong with this.\n"+err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Conten-type", "application/json")
	json.NewEncoder(w).Encode(rows)
}

func handleOrderView(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Something wrong with the ID.\n"+err.Error(), http.StatusBadRequest)
		return
	}
	thisRow, err := order.ReadUUID(dbConn, id)
	if err != nil {
		http.Error(w, "Item's ID is wrong.\n"+err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(thisRow)
}
func handleOrdersView(w http.ResponseWriter, r *http.Request) {
	rows, err := order.Read(dbConn)
	if err != nil {
		http.Error(w, "Something wrong with this.\n"+err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Conten-type", "application/json")
	json.NewEncoder(w).Encode(rows)
}

func handleItemTypesView(w http.ResponseWriter, r *http.Request) {
	log.Println("[ItemTypes] View")
	types, err := item.GetItemTypes(dbConn)
	if err != nil {
		http.Error(w, "Can't get types: "+err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Conten-type", "application/json")
	json.NewEncoder(w).Encode(types)
}

//Edits
func handleCartEdit(w http.ResponseWriter, r *http.Request) {
	log.Println("[POST] CartEdit")
	vars := mux.Vars(r)
	thisObject := cart.Cart{}
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Something wrong with the ID.\n"+err.Error(), http.StatusBadRequest)
		return
	}
	thisObject.UUID = id
	if err := json.NewDecoder(r.Body).Decode(&thisObject); err != nil {
		http.Error(w, "Can't parse request body.\n"+err.Error(), http.StatusBadRequest)
		return
	}

	if err := cart.Update(dbConn, thisObject); err != nil {
		http.Error(w, "Can't update data.\n"+err.Error(), http.StatusBadRequest)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	http.Redirect(w, r, fmt.Sprintf("%s", r.URL.Path), http.StatusFound)
}

func handleClientEdit(w http.ResponseWriter, r *http.Request) {
	log.Println("[POST] ClientEdit")
	vars := mux.Vars(r)
	thisObject := client.Client{}
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Something wrong with the ID.\n"+err.Error(), http.StatusBadRequest)
		return
	}
	thisObject.UUID = id
	if err := json.NewDecoder(r.Body).Decode(&thisObject); err != nil {
		http.Error(w, "Can't parse request body.\n"+err.Error(), http.StatusBadRequest)
		return
	}

	if err := client.Update(dbConn, thisObject); err != nil {
		http.Error(w, "Can't update data.\n"+err.Error(), http.StatusBadRequest)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	http.Redirect(w, r, fmt.Sprintf("%s", r.URL.Path), http.StatusFound)
}

func handleClientsEdit(w http.ResponseWriter, r *http.Request) {
	log.Println("[POST] ClientsEdit")
	var this []client.Client
	if err := json.NewDecoder(r.Body).Decode(&this); err != nil {
		http.Error(w, "Can't parse request body.\n"+err.Error(), http.StatusBadRequest)
		return
	}
	if err := client.UpdateBundle(dbConn, this); err != nil {
		http.Error(w, "Can't update data.\n"+err.Error(), http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func handleItemEdit(w http.ResponseWriter, r *http.Request) {
	log.Println("[POST] ItemEdit")
	vars := mux.Vars(r)
	thisObject := item.Item{}
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Something wrong with the ID.\n"+err.Error(), http.StatusBadRequest)
		return
	}
	thisObject.UUID = id
	if err := json.NewDecoder(r.Body).Decode(&thisObject); err != nil {
		http.Error(w, "Can't parse request body.\n"+err.Error(), http.StatusBadRequest)
		return
	}

	if err := item.Update(dbConn, thisObject); err != nil {
		http.Error(w, "Can't update data.\n"+err.Error(), http.StatusBadRequest)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}
func handleItemsEdit(w http.ResponseWriter, r *http.Request) {
	log.Println("[POST] ItemsEdit")
	var this []item.Item
	if err := json.NewDecoder(r.Body).Decode(&this); err != nil {
		http.Error(w, "Can't parse request body.\n"+err.Error(), http.StatusBadRequest)
		return
	}
	if err := item.UpdateBundle(dbConn, this); err != nil {
		http.Error(w, "Can't update data.\n"+err.Error(), http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func handleOrderEdit(w http.ResponseWriter, r *http.Request) {
	log.Println("[POST] OrderEdit")
	vars := mux.Vars(r)
	thisObject := order.Order{}
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Something wrong with the ID.\n"+err.Error(), http.StatusBadRequest)
		return
	}
	thisObject.UUID = id
	if err := json.NewDecoder(r.Body).Decode(&thisObject); err != nil {
		http.Error(w, "Can't parse request body.\n"+err.Error(), http.StatusBadRequest)
		return
	}

	if err := order.Update(dbConn, thisObject); err != nil {
		http.Error(w, "Can't update data.\n"+err.Error(), http.StatusBadRequest)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}
func handleOrdersEdit(w http.ResponseWriter, r *http.Request) {
	log.Println("[POST] Orders")
	var this []order.Order
	if err := json.NewDecoder(r.Body).Decode(&this); err != nil {
		http.Error(w, "Can't parse request body.\n"+err.Error(), http.StatusBadRequest)
		return
	}
	if err := order.UpdateBundle(dbConn, this); err != nil {
		http.Error(w, "Can't update data.\n"+err.Error(), http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

//Delete
func handleCartDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Something wrong with the ID.\n"+err.Error(), http.StatusBadRequest)
		return
	}
	err = cart.Delete(dbConn, id)
	if err != nil {
		http.Error(w, "Something went wrong.\n"+err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}
func handleClientDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Something wrong with the ID.\n"+err.Error(), http.StatusBadRequest)
		return
	}
	err = client.Delete(dbConn, id)
	if err != nil {
		http.Error(w, "Something went wrong.\n"+err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func handleClientsDelete(w http.ResponseWriter, r *http.Request) {
	log.Println("[DELETE] Clients")
	var ids []uuid.UUID

	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		http.Error(w, "Can't parse request body.\n"+err.Error(), http.StatusBadRequest)
		return
	}

	err := client.DeleteBundle(dbConn, ids)
	if err != nil {
		http.Error(w, "Something went wrong.\n"+err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Add("status", "200")
	fmt.Fprintf(w, "OK")
}

func handleItemDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Something wrong with the ID.\n"+err.Error(), http.StatusBadRequest)
		return
	}
	err = item.Delete(dbConn, id)
	if err != nil {
		http.Error(w, "Something went wrong.\n"+err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}
func handleItemsDelete(w http.ResponseWriter, r *http.Request) {
	log.Println("[DELETE] Items")
	var ids []uuid.UUID

	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		http.Error(w, "Can't parse request body.\n"+err.Error(), http.StatusBadRequest)
		return
	}

	err := item.DeleteBundle(dbConn, ids)
	if err != nil {
		http.Error(w, "Something went wrong.\n"+err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Add("status", "200")
	fmt.Fprintf(w, "OK")
}

func handleOrderDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Something wrong with the ID.\n"+err.Error(), http.StatusBadRequest)
		return
	}
	err = order.Delete(dbConn, id)
	if err != nil {
		http.Error(w, "Something went wrong.\n"+err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")

}
func handleOrdersDelete(w http.ResponseWriter, r *http.Request) {
	log.Println("[DELETE] Orders")
	var ids []uuid.UUID

	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		http.Error(w, "Can't parse request body.\n"+err.Error(), http.StatusBadRequest)
		return
	}

	err := order.DeleteBundle(dbConn, ids)
	if err != nil {
		http.Error(w, "Something went wrong.\n"+err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Add("status", "200")
	fmt.Fprintf(w, "OK")
}

//Create
func handleCartCreate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var thisObject cart.Cart
	if err := json.NewDecoder(r.Body).Decode(&thisObject); err != nil {
		http.Error(w, "Can't parse form.\nReason:"+err.Error(), http.StatusBadRequest)
		return
	}
	id, err := cart.Create(dbConn, thisObject)
	if err != nil {
		http.Error(w, "Can't create object.\nReason:"+err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Cart has been created. UUID: %s", id.String())
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func handleClientCreate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var thisObject client.Client
	if err := json.NewDecoder(r.Body).Decode(&thisObject); err != nil {
		http.Error(w, "Can't parse form.\nReason:"+err.Error(), http.StatusBadRequest)
		return
	}
	id, err := client.Create(dbConn, thisObject)
	if err != nil {
		http.Error(w, "Can't create object.\nReason:"+err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Client has been created. UUID: %s", id.String())
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}
func handleClientsCreate(w http.ResponseWriter, r *http.Request) {
	log.Println("[CREATE] Clients")
	var requestData []client.Client

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Can't parse request body.\n"+err.Error(), http.StatusBadRequest)
		return
	}

	err := client.CreateBundle(dbConn, requestData)
	if err != nil {
		http.Error(w, "Something went wrong.\n"+err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Add("status", "200")
	fmt.Fprintf(w, "OK")
}

func handleItemCreate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var thisObject item.Item
	if err := json.NewDecoder(r.Body).Decode(&thisObject); err != nil {
		http.Error(w, "Can't parse form.\nReason:"+err.Error(), http.StatusBadRequest)
		return
	}
	id, err := item.Create(dbConn, thisObject)
	if err != nil {
		http.Error(w, "Can't create object.\nReason: "+err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Item has been created. UUID: %s", id.String())
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}
func handleItemsCreate(w http.ResponseWriter, r *http.Request) {
	log.Println("[CREATE] Items")
	var requestData []item.Item

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Can't parse request body.\n"+err.Error(), http.StatusBadRequest)
		return
	}

	err := item.CreateBundle(dbConn, requestData)
	if err != nil {
		http.Error(w, "Something went wrong.\n"+err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Add("status", "200")
	fmt.Fprintf(w, "OK")
}

func handleOrderCreate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var thisObject order.Order
	if err := json.NewDecoder(r.Body).Decode(&thisObject); err != nil {
		http.Error(w, "Can't parse form.\nReason:"+err.Error(), http.StatusBadRequest)
		return
	}
	id, err := order.Create(dbConn, thisObject)
	if err != nil {
		http.Error(w, "Can't create object.\nReason:"+err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Order has been created. UUID: %s", id.String())
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}
func handleOrdersCreate(w http.ResponseWriter, r *http.Request) {
	log.Println("[CREATE] Orders")
	var requestData []order.Order

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Can't parse request body.\n"+err.Error(), http.StatusBadRequest)
		return
	}

	err := order.CreateBundle(dbConn, requestData)
	if err != nil {
		http.Error(w, "Something went wrong.\n"+err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Add("status", "200")
	fmt.Fprintf(w, "OK")
}
