package civo

import (
	"time"

	"github.com/kubesimplify/ksctl/pkg/resources"
	. "github.com/kubesimplify/ksctl/pkg/utils/consts"
)

// NewNetwork implements resources.CloudFactory.
func (obj *CivoProvider) NewNetwork(storage resources.StorageFactory) error {
	name := obj.metadata.resName
	obj.mxName.Unlock()

	log.Debug("Printing", "Name", name)

	// check if the networkID already exist
	if len(civoCloudState.NetworkIDs.NetworkID) != 0 {
		log.Print("skipped network creation found", "networkID", civoCloudState.NetworkIDs.NetworkID)
		return nil
	}

	res, err := obj.client.CreateNetwork(name)
	if err != nil {
		return log.NewError(err.Error())
	}
	civoCloudState.NetworkIDs.NetworkID = res.ID
	log.Debug("Printing", "networkID", res.ID)
	log.Success("Created network", "name", name)

	// NOTE: as network creation marks first resource we should create the directoy
	// when its success

	if err := storage.Path(generatePath(UtilClusterPath, clusterType, clusterDirName)).
		Permission(FILE_PERM_CLUSTER_DIR).CreateDir(); err != nil {
		return log.NewError(err.Error())
	}

	path := generatePath(UtilClusterPath, clusterType, clusterDirName, STATE_FILE_NAME)

	return log.NewError(saveStateHelper(storage, path).Error())
}

// DelNetwork implements resources.CloudFactory.
func (obj *CivoProvider) DelNetwork(storage resources.StorageFactory) error {

	if len(civoCloudState.NetworkIDs.NetworkID) == 0 {
		log.Print("skipped network already deleted")
	} else {

		currRetryCounter := KsctlCounterConsts(0)
		for currRetryCounter < CounterMaxWatchRetryCount {
			var err error
			_, err = obj.client.DeleteNetwork(civoCloudState.NetworkIDs.NetworkID)
			if err != nil {
				currRetryCounter++
				log.Warn("RETRYING", err)
			} else {
				break
			}
			time.Sleep(5 * time.Second)
		}
		if currRetryCounter == CounterMaxWatchRetryCount {
			return log.NewError("failed to delete network timeout")
		}

		civoCloudState.NetworkIDs.NetworkID = ""
		if err := saveStateHelper(storage, generatePath(UtilClusterPath, clusterType, clusterDirName, STATE_FILE_NAME)); err != nil {
			return log.NewError(err.Error())
		}
		log.Success("Deleted network", "networkID", civoCloudState.NetworkIDs.NetworkID)
	}
	path := generatePath(UtilClusterPath, clusterType, clusterDirName)

	return log.NewError(storage.Path(path).DeleteDir().Error())
}
