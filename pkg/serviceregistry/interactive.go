package serviceregistry

import (
	"context"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/redhat-developer/app-services-cli/pkg/connection"
	"github.com/redhat-developer/app-services-cli/pkg/logging"
	srsmgmtv1 "github.com/redhat-developer/app-services-sdk-go/registrymgmt/apiv1/client"
)

const (
	queryLimit = 1000
)

func InteractiveSelect(ctx context.Context, connection connection.Connection, logger logging.Logger) (*srsmgmtv1.Registry, error) {
	api := connection.API()

	response, _, err := api.ServiceRegistryMgmt().GetRegistries(ctx).Size(queryLimit).Execute()
	if err != nil {
		return nil, fmt.Errorf("unable to list Service Registry instances: %w", err)
	}

	if response.Size == 0 {
		logger.Info("No Service Registry instances were found.")
		return nil, nil
	}

	regisries := []string{}
	for index := 0; index < len(response.Items); index++ {
		regisries = append(regisries, *response.Items[index].Name)
	}

	prompt := &survey.Select{
		Message:  "Select Service Registry instance to connect:",
		Options:  regisries,
		PageSize: 10,
	}

	var selectedRegistryIndex int
	err = survey.AskOne(prompt, &selectedRegistryIndex)
	if err != nil {
		return nil, err
	}

	selectedRegistry := response.Items[selectedRegistryIndex]

	return &selectedRegistry, nil
}
