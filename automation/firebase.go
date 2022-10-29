package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

func getPackages(client firestore.Client, func_ctx context.Context) []AutomationPackage {
	// get all documents from the wpa-packages collection
	iter := client.Collection("wpa-packages").Documents(func_ctx)

	// iterate over the documents and add them to the packages slice
	var packages []AutomationPackage
	for {
		doc, err := iter.Next()
		if err != nil {
			log.Fatalf("error getting next document: %v\n", err)
		}

		var pkg AutomationPackage
		doc.DataTo(&pkg)
		packages = append(packages, pkg)
	}

	return packages
}

func updatePkgInfo(fieldName, fieldValue, packageIdentifier string, client *firestore.Client, func_ctx context.Context) {
	_, err := client.Doc("wpa-packages/"+packageIdentifier).Update(
		func_ctx,
		[]firestore.Update{
			{
				Path:  fieldName,
				Value: fieldValue,
			},
		},
	)
	if err != nil {
		log.Fatalf("error updating previous version: %v\n", err)
	}
}

func getFirestoreClient(func_ctx context.Context) (firestore.Client, error) {
	creds := os.Getenv("FB_CREDS")
	sa := option.WithCredentialsJSON([]byte(creds))
	app, err := firebase.NewApp(func_ctx, nil, sa)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	client, err := app.Firestore(func_ctx)
	if err != nil {
		return *client, fmt.Errorf("error initializing firestore client: %v\n", err)
	}

	return *client, nil
}
