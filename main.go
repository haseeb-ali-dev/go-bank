package main

import (
	"flag"
	"fmt"
	"log"
)

func seedAccount(store Storage, fname, lname, pw string) *Account {
	account, err := NewAccount(fname, lname, pw)
	if err != nil {
		log.Fatal(err)
	}

	if err := store.CreateAccount(account); err != nil {
		log.Fatal(err)
	}

	fmt.Println("new account ==> ", account.Number)

	return account
}

func seedData(s Storage) {
	seedAccount(s, "Ali", "Raza", "password@123")
}

func main() {
	seed := flag.Bool("seed", false, "Seed the database")
	flag.Parse()

	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	if *seed {
		fmt.Println("Seeding database...")
		seedData(store)
	}

	server := NewAPIServer(":3000", store)
	server.Run()
}
