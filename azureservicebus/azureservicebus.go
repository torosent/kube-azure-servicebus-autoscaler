package azureservicebus

import (
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"

	"github.com/Azure/azure-sdk-for-go/arm/examples/helpers"
	"github.com/Azure/azure-sdk-for-go/arm/servicebus"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// NumMessages Number of messages in queue
func NumMessages(resourcegroup string, queuename string, namespace string) (int, error) {

	c := map[string]string{
		"AZURE_CLIENT_ID":       os.Getenv("AZURE_CLIENT_ID"),
		"AZURE_CLIENT_SECRET":   os.Getenv("AZURE_CLIENT_SECRET"),
		"AZURE_SUBSCRIPTION_ID": os.Getenv("AZURE_SUBSCRIPTION_ID"),
		"AZURE_TENANT_ID":       os.Getenv("AZURE_TENANT_ID")}

	if err := checkEnvVar(&c); err != nil {
		log.Fatalf("Error: %v", err)
		return 0, err
	}
	spt, err := helpers.NewServicePrincipalTokenFromCredentials(c, azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		log.Fatalf("Error: %v", err)
		return 0, err
	}

	client := servicebus.NewQueuesClient(c["AZURE_SUBSCRIPTION_ID"])
	client.Authorizer = autorest.NewBearerAuthorizer(spt)

	result, err := client.Get(resourcegroup, namespace, queuename)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to get number of messages in queue")
	}
	num := int(*result.MessageCount)
	return num, nil
}

func checkEnvVar(envVars *map[string]string) error {
	var missingVars []string
	for varName, value := range *envVars {
		if value == "" {
			missingVars = append(missingVars, varName)
		}
	}
	if len(missingVars) > 0 {
		return fmt.Errorf("Missing environment variables %v", missingVars)
	}
	return nil
}
