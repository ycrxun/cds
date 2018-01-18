package main

import (
	"fmt"

	"github.com/ovh/cds/sdk"
	"github.com/ovh/cds/sdk/vcs"
	"github.com/ovh/cds/sdk/vcs/git"
)

func extractVCSInformations(params []sdk.Parameter) (string, *git.AuthOpts, error) {
<<<<<<< HEAD
	var gitURL string
	var auth *git.AuthOpts

	// Get connection type
	connetionType := sdk.ParameterFind(&params, "git.connection.type")
	if connetionType == nil || (connetionType.Value != "ssh" && connetionType.Value != "https") {
		return gitURL, nil, fmt.Errorf("git connection type is not set. nothing to perform")
=======
	var gitUrl string
	var auth *git.AuthOpts

	// Get connection type
	connetionType := sdk.ParameterFind(params, "git.connection.type")
	if connetionType == nil || (connetionType.Value != "ssh" && connetionType.Value != "https") {
		return gitUrl, nil, fmt.Errorf("git connection type is not set. nothing to perform")
>>>>>>> master
	}

	switch connetionType.Value {
	case "ssh":
<<<<<<< HEAD
		keyName := sdk.ParameterFind(&params, "git.ssh.key")
		if keyName == nil || keyName.Value == "" {
			return gitURL, nil, fmt.Errorf("git ssh key is not set. nothing to perform")
		}

		privateKey := sdk.ParameterFind(&params, "cds.key."+keyName.Value+".priv")
		if privateKey == nil || privateKey.Value == "" {
			return gitURL, nil, fmt.Errorf("ssh key not found. Nothing to perform")
		}
		if err := vcs.SetupSSHKey(nil, keysDirectory, privateKey); err != nil {
			return gitURL, nil, fmt.Errorf("unable to setup ssh key. %s", err)
		}
		key, errK := vcs.GetSSHKey(params, keysDirectory, privateKey)
		if errK != nil && errK != sdk.ErrKeyNotFound {
			return gitURL, nil, fmt.Errorf("unable to setup ssh key. %s", errK)
=======
		keyName := sdk.ParameterFind(params, "git.ssh.key")
		if keyName == nil || keyName.Value == "" {
			return gitUrl, nil, fmt.Errorf("git ssh key is not set. nothing to perform")
		}

		privateKey := sdk.ParameterFind(params, "cds.key."+keyName.Value+".priv")
		if privateKey == nil || privateKey.Value == "" {
			return gitUrl, nil, fmt.Errorf("ssh key not found. Nothing to perform")
		}
		if err := vcs.SetupSSHKey(nil, keysDirectory, privateKey); err != nil {
			return gitUrl, nil, fmt.Errorf("unable to setup ssh key. %s", err)
		}
		key, errK := vcs.GetSSHKey(params, keysDirectory, privateKey)
		if errK != nil && errK != sdk.ErrKeyNotFound {
			return gitUrl, nil, fmt.Errorf("unable to setup ssh key. %s", errK)
>>>>>>> master
		}
		if key != nil {
			if auth == nil {
				auth = new(git.AuthOpts)
			}
			auth.PrivateKey = *key
		}

<<<<<<< HEAD
		url := sdk.ParameterFind(&params, "git.url")
		if url == nil || url.Value == "" {
			return gitURL, nil, fmt.Errorf("SSH Url (git.url) not found. Nothing to perform")
		}
		gitURL = url.Value
	case "https":
		user := sdk.ParameterFind(&params, "git.http.user")
		password := sdk.ParameterFind(&params, "git.http.password")
=======
		url := sdk.ParameterFind(params, "git.url")
		if url == nil || url.Value == "" {
			return gitUrl, nil, fmt.Errorf("SSH Url (git.url) not found. Nothing to perform")
		}
		gitUrl = url.Value
	case "https":
		user := sdk.ParameterFind(params, "git.http.user")
		password := sdk.ParameterFind(params, "git.http.password")
>>>>>>> master

		if user != nil || password != nil {
			auth = new(git.AuthOpts)
			if user != nil {
				auth.Username = user.Value
			}
			if password != nil {
				auth.Password = password.Value
			}
		}

<<<<<<< HEAD
		url := sdk.ParameterFind(&params, "git.http_url")
		if url == nil || url.Value == "" {
			return gitURL, nil, fmt.Errorf("SSH Url (git.http_url) not found. Nothing to perform")
		}
		gitURL = url.Value
	}

	pgpKeyName := sdk.ParameterFind(&params, "git.pgp.key")
=======
		url := sdk.ParameterFind(params, "git.http_url")
		if url == nil || url.Value == "" {
			return gitUrl, nil, fmt.Errorf("SSH Url (git.http_url) not found. Nothing to perform")
		}
		gitUrl = url.Value
	}

	pgpKeyName := sdk.ParameterFind(params, "git.pgp.key")
>>>>>>> master
	if pgpKeyName != nil && pgpKeyName.Value != "" {
		auth.SignKey = vcs.PGPKey{
			Name: pgpKeyName.Value,
		}
<<<<<<< HEAD
		pgpPrivate := sdk.ParameterFind(&params, "cds.key."+pgpKeyName.Value+".priv")
		if pgpPrivate != nil && pgpPrivate.Value != "" {
			auth.SignKey.Private = pgpPrivate.Value
		}
		pgpPublic := sdk.ParameterFind(&params, "cds.key."+pgpKeyName.Value+".pub")
		if pgpPublic != nil && pgpPublic.Value != "" {
			auth.SignKey.Public = pgpPublic.Value
		}
		pgpID := sdk.ParameterFind(&params, "cds.key."+pgpKeyName.Value+".id")
=======
		pgpPrivate := sdk.ParameterFind(params, "cds.key."+pgpKeyName.Value+".priv")
		if pgpPrivate != nil && pgpPrivate.Value != "" {
			auth.SignKey.Private = pgpPrivate.Value
		}
		pgpPublic := sdk.ParameterFind(params, "cds.key."+pgpKeyName.Value+".pub")
		if pgpPublic != nil && pgpPublic.Value != "" {
			auth.SignKey.Public = pgpPublic.Value
		}
		pgpID := sdk.ParameterFind(params, "cds.key."+pgpKeyName.Value+".id")
>>>>>>> master
		if pgpID != nil && pgpID.Value != "" {
			auth.SignKey.ID = pgpID.Value
		}

	}
<<<<<<< HEAD
	return gitURL, auth, nil
=======
	return gitUrl, auth, nil
>>>>>>> master
}
