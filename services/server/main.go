package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	authv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var k8sClient *kubernetes.Clientset

func setup() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}
	k8sClient, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
}

func verifyToken(clientID string) (bool, error) {
	log.Printf("clientID : %s\n", clientID)
	ctx := context.TODO()
	tr := authv1.TokenReview{
		Spec: authv1.TokenReviewSpec{
			Token:     clientID,
			Audiences: []string{"server"},
		},
	}
	result, err := k8sClient.AuthenticationV1().TokenReviews().Create(ctx, &tr, metav1.CreateOptions{})
	if err != nil {
		return false, err
	}
	s, _ := json.MarshalIndent(result, "", "\t")
	log.Printf("%s\n", s)

	if result.Status.Authenticated {
		return true, nil
	}
	return false, nil

}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	// Read the value of the client identifier from the request header
	clientId := r.Header.Get("X-Client-Id")
	if len(clientId) == 0 {
		http.Error(w, "X-Client-Id not supplied", http.StatusUnauthorized)
		return
	}
	authenticated, err := verifyToken(clientId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !authenticated {
		http.Error(w, "Invalid token", http.StatusForbidden)
		return
	}
	io.WriteString(w, "Hello from data store. You have been authenticated")
}

func main() {
	setup()

	http.HandleFunc("/", handleIndex)
	http.ListenAndServe(":8081", nil)
}
