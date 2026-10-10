package alldebrid_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	alldebrid "github.com/pierre-emmanuelJ/open-alldebrid/sdks/go"
)

func Example() {
	ctx := context.Background()
	client := alldebrid.New(alldebrid.DemoPremiumKey, alldebrid.WithUserAgent("my-app"))

	user, err := client.User(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user.Username, user.IsPremium)
}

func ExampleClient_WaitForPin() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	client := alldebrid.New("", alldebrid.WithUserAgent("my-app"))

	pin, err := client.GetPin(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Open", pin.UserURL, "to authorize this app")
	key, err := client.WaitForPin(ctx, pin, 0)
	if err != nil {
		log.Fatal(err)
	}
	client = client.WithAPIKey(key)
}

func ExampleClient_Unlock() {
	ctx := context.Background()
	client := alldebrid.New(os.Getenv("ALLDEBRID_API_KEY"))

	u, err := client.Unlock(ctx, "https://1fichier.com/?abcdefghij", "")
	switch {
	case errors.Is(err, alldebrid.CodeLinkDown):
		fmt.Println("the file is gone")
		return
	case err != nil:
		log.Fatal(err)
	case u.Delayed != "":
		link, err := client.WaitDelayedLink(ctx, u.Delayed.String(), 0)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(link)
	default:
		fmt.Println(u.Link)
	}
}

func ExampleClient_LiveMagnets() {
	ctx := context.Background()
	client := alldebrid.New(os.Getenv("ALLDEBRID_API_KEY"))

	live := client.LiveMagnets()
	for range time.Tick(5 * time.Second) {
		magnets, err := live.Poll(ctx)
		if err != nil {
			log.Fatal(err)
		}
		for _, m := range magnets {
			fmt.Printf("%s %s %d%%\n", m.Filename, m.Status, percent(m))
		}
	}
}

func percent(m alldebrid.Magnet) int64 {
	if m.Size == 0 {
		return 0
	}
	return m.Downloaded * 100 / m.Size
}
