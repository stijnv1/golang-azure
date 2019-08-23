package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-06-01/compute"
	"github.com/gorilla/mux"
	"github.com/stijnv1/golang-azure/internal/config"
	"github.com/stijnv1/golang-azure/internal/iam"
	"github.com/stijnv1/golang-azure/models"
)

var (
	ctx = context.Background()
)

func addLocalEnvAndParse() error {
	// parse env at top-level (also controls dotenv load)
	err := config.ParseEnvironment()
	if err != nil {
		return fmt.Errorf("failed to add top-level env: %v\n", err.Error())
	}
	return nil
}

func getVMClient() compute.VirtualMachinesClient {
	vmClient := compute.NewVirtualMachinesClient(config.SubscriptionID())
	a, _ := iam.GetResourceManagementAuthorizer()
	vmClient.Authorizer = a
	vmClient.AddToUserAgent(config.UserAgent())
	return vmClient
}

func GetAzureVMs(w http.ResponseWriter, r *http.Request) {
	vmClient := getVMClient()
	enableCors(&w)
	//vms, err := vmClient.ListAll(ctx)
	var azurevmlist models.AzureVMList

	for list, _ := vmClient.ListAllComplete(ctx); list.NotDone(); list.Next() {
		var azurevm models.AzureVM

		azurevm.Name = *list.Value().Name
		azurevm.VMID = *list.Value().VMID

		azurevmlist = append(azurevmlist, azurevm)
	}

	//if err != nil {
	//json.NewEncoder(w).Encode("error")
	//} else {
	json.NewEncoder(w).Encode(azurevmlist)
	//}
}

func GetAzureVMs_v2(w http.ResponseWriter, r *http.Request) {
	vmClient := getVMClient()
	vmList, err := vmClient.ListAllComplete(ctx)

	if err != nil {
		fmt.Fprint(w, "error occured: ", err.Error())
		return
	}

	json.NewEncoder(w).Encode((vmList.Response()))
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func main() {
	err := addLocalEnvAndParse()

	if err == nil {
		router := mux.NewRouter()
		router.HandleFunc("/getazurevms", GetAzureVMs).Methods("GET")
		router.HandleFunc("/getazurevms_v2", GetAzureVMs_v2).Methods("GET")
		log.Fatal(http.ListenAndServe(":8000", router))
	}

}
