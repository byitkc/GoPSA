package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/byitkc/gopsa"
	"github.com/joho/godotenv"
)

func main() {
	config := initEverything()

	companyInfo, err := gopsa.GetCompanyInfo(http.DefaultClient, config.Parameters.Protocol, config.Parameters.Host, config.Parameters.CompanyID)
	if err != nil {
		slog.Error("failed to retrieve ConnectWise PSA company information")
	}
	fmt.Printf("%+v", companyInfo)
	config.Parameters.APIURLBase = fmt.Sprintf("%s://%s/%s", config.Parameters.Protocol, companyInfo.SiteURL, companyInfo.Codebase)

	fmt.Printf("%+v\n", companyInfo)
}

func initEverything() gopsa.Config {
	envFilePath := ".env"
	creds, connParams := initEnv(envFilePath)

	config := gopsa.Config{
		Credentials: creds,
		Parameters:  connParams,
	}

	return config
}

func initEnv(envFilePath string) (gopsa.Credentials, gopsa.ConnectionParameters) {
	godotenv.Load(envFilePath)

	creds := gopsa.Credentials{
		PublicKey:  stringEnvMustLoad("CW_PSA_PUBLIC_KEY"),
		PrivateKey: stringEnvMustLoad("CW_PSA_PRIVATE_KEY"),
		ClientID:   stringEnvMustLoad("CW_PSA_CLIENT_ID"),
	}

	protocol := stringEnvMustLoad("CW_PSA_PROTOCOL")
	if protocol != "http" && protocol != "https" {
		slog.Error("invalid protocol specified", "protocol", protocol)
		os.Exit(1)
	}
	companyID := stringEnvMustLoad("CW_PSA_COMPANY_ID")
	host := stringEnvMustLoad("CW_PSA_HOST")

	params := gopsa.ConnectionParameters{
		Protocol:  protocol,
		CompanyID: companyID,
		Host:      host,
		URLBase:   fmt.Sprintf("%s://%s", protocol, host),
	}

	return creds, params
}

func stringEnvMustLoad(variable string) string {
	s := os.Getenv(variable)
	if s == "" {
		slog.Error("required environment variable could not be loaded", "variable", variable)
	}
	return s
}
