package updater

import (
	"fmt"

	"github.com/apex/log"
	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

type Updater struct {
	CurrentVersion string
	Storage        MetaInfoStorage
}

func (u Updater) doSelfUpdate() (bool, error) {
	v, err := semver.Parse(u.CurrentVersion[1:])
	if err != nil {
		return false, fmt.Errorf("parse current version failed: %w", err)
	}
	latest, err := selfupdate.UpdateSelf(v, "racoon-devel/RaccoonPirate")
	if err != nil {
		return false, err
	}

	if latest.Version.Equals(v) {
		log.Info("Nothing to update")
		return false, nil
	}
	log.Infof("Successfully updated to %s", latest.Version)
	return true, nil
}

func (u Updater) updateStorage() error {
	previousVersion, err := u.Storage.GetVersion()
	if err != nil {
		return fmt.Errorf("failed to load previous version: %s", err)
	}
	if previousVersion == "" || previousVersion == "0.0.0" {
		return u.Storage.SetVersion(u.CurrentVersion)
	}
	if previousVersion != u.CurrentVersion {
		// perfoming storage update
		return u.Storage.SetVersion(u.CurrentVersion)
	}
	return nil
}

func (u Updater) TryUpdate() (updated bool, err error) {
	if u.CurrentVersion == "0.0.0" {
		return
	}

	if err := u.updateStorage(); err != nil {
		log.Warnf("Update database failed: %s", err)
	}

	updated, err = u.doSelfUpdate()
	return
}
