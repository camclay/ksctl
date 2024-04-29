package controllers

import (
	"github.com/ksctl/ksctl/pkg/helpers/consts"
	"github.com/ksctl/ksctl/pkg/resources"
)

// Controller TODO: use the ksctlClient as the struct to whom the function act as a method
type Controller interface {
	CreateManagedCluster(*resources.KsctlClient) error
	DeleteManagedCluster(*resources.KsctlClient) error

	SwitchCluster(*resources.KsctlClient) (*string, error)
	Applications(*resources.KsctlClient, consts.KsctlOperation) error
	GetCluster(*resources.KsctlClient) error

	Credentials(*resources.KsctlClient) error

	CreateHACluster(*resources.KsctlClient) error
	DeleteHACluster(*resources.KsctlClient) error

	AddWorkerPlaneNode(*resources.KsctlClient) error
	DelWorkerPlaneNode(*resources.KsctlClient) error
}
