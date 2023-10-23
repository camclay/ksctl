package civo

import (
	"github.com/kubesimplify/ksctl/pkg/resources"
	"github.com/kubesimplify/ksctl/pkg/utils"
	. "github.com/kubesimplify/ksctl/pkg/utils/consts"
)

// DelSSHKeyPair implements resources.CloudFactory.
func (obj *CivoProvider) DelSSHKeyPair(storage resources.StorageFactory) error {
	if len(civoCloudState.SSHID) == 0 {
		storage.Logger().Success("[skip] ssh keypair already deleted")
		return nil
	}

	_, err := obj.client.DeleteSSHKey(civoCloudState.SSHID)
	if err != nil {
		return err
	}
	path := generatePath(UtilClusterPath, clusterType, clusterDirName, STATE_FILE_NAME)

	storage.Logger().Success("[civo] ssh keypair deleted", civoCloudState.SSHID)

	civoCloudState.SSHID = ""
	civoCloudState.SSHPrivateKeyLoc = ""
	civoCloudState.SSHUser = ""

	return saveStateHelper(storage, path)
}

// CreateUploadSSHKeyPair implements resources.CloudFactory.
func (obj *CivoProvider) CreateUploadSSHKeyPair(storage resources.StorageFactory) error {
	name := obj.metadata.resName
	obj.mxName.Unlock()

	if len(civoCloudState.SSHID) != 0 {
		storage.Logger().Success("[skip] ssh keypair already uploaded")
		return nil
	}

	keyPairToUpload, err := utils.CreateSSHKeyPair(storage, CloudCivo, clusterDirName)
	if err != nil {
		return err
	}
	if err := obj.uploadSSH(storage, name, keyPairToUpload); err != nil {
		return err
	}
	storage.Logger().Success("[civo] ssh keypair created and uploaded", name)
	return nil
}

func (obj *CivoProvider) uploadSSH(storage resources.StorageFactory, resName, pubKey string) error {
	sshResp, err := obj.client.NewSSHKey(resName, pubKey)
	if err != nil {
		return err
	}

	civoCloudState.SSHID = sshResp.ID
	civoCloudState.SSHUser = "root"
	civoCloudState.SSHPrivateKeyLoc = utils.GetPath(UtilSSHPath, CloudCivo, clusterType, clusterDirName)

	path := generatePath(UtilClusterPath, clusterType, clusterDirName, STATE_FILE_NAME)

	return saveStateHelper(storage, path)
}