package app

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/vps2/cisco-switches-crawler/internal/infrastructure/cisco"
	"github.com/vps2/cisco-switches-crawler/internal/usecase"
	"github.com/vps2/cisco-switches-crawler/pkg/ip"
	"github.com/vps2/cisco-switches-crawler/pkg/telnet"
	"golang.org/x/term"
)

var (
	rootDevIP string
	verbose   bool
	include   string
	pretty    bool
)

var (
	user     string
	password string
)

func Run() {
	flag.StringVar(&rootDevIP, "address", "", "ip address of the switch")
	flag.StringVar(&user, "user", "", "the name of the user to access the switches")
	flag.StringVar(&password, "password", "", "the user's password. If not specified, the application will ask for a password")
	flag.BoolVar(&verbose, "verbose", false, "show verbose")
	flag.StringVar(&include, "include", "", "ip addresses (separated by commas) included in the selection. Example: [192.168.1.1,192.168.1.0/24]")
	flag.BoolVar(&pretty, "pretty", false, "beautiful print of the result")
	flag.Parse()

	if rootDevIP == "" {
		log.Fatal("IP address of the switch is empty")
	}
	if ok := checkIP(rootDevIP); !ok {
		log.Fatal("IP address of the switch is incorrect")
	}

	if user == "" {
		log.Fatal("The user name for accessing the switches is not set")
	}
	if password == "" {
		fmt.Print("password: ")
		var err error
		if password, err = readUserPassword(); err != nil {
			log.Fatal("Error receiving the user's password")
		}

		fmt.Println()

		if password == "" {
			log.Fatal("An empty password is not allowed")
		}
	}

	ipFilter := ip.NewFilter(ip.AllowAnyIfEmpty(true))
	if include != "" {
		includeIPs := strings.Split(include, ",")
		for _, ip := range includeIPs {
			if err := ipFilter.Add(strings.TrimSpace(ip)); err != nil {
				log.Fatal("Include parameter has an incorrect value of ip addresses or incorrect format")
			}
		}
	}

	//--------------------------------------------------------------------------------------------------------------------

	telnet := telnet.New()
	var networkBuilder *usecase.NetworkBuilder

	if verbose {
		client := cisco.NewClient(telnet, cisco.WithVerbose())
		networkBuilder = usecase.NewNetworkBuilder(client, usecase.WithShowOutput(), usecase.WithIPFiltering(ipFilter))
	} else {
		client := cisco.NewClient(telnet)
		networkBuilder = usecase.NewNetworkBuilder(client, usecase.WithIPFiltering(ipFilter))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt)

		select {
		case <-done:
			cancel()
		case <-ctx.Done():
			//exit
		}
	}()

	networkBuilder.Build(ctx, rootDevIP, user, password)
	fmt.Println()
	if pretty {
		fmt.Println(string(networkBuilder.ToPrettyJSON()))
	} else {
		fmt.Println(string(networkBuilder.ToJSON()))
	}
}

func checkIP(ip string) bool {
	if ip := net.ParseIP(ip); ip != nil {
		return true
	}

	return false
}

func readUserPassword() (string, error) {
	pwd, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", fmt.Errorf("get password error: %w", err)
	}

	return string(pwd), nil
}
