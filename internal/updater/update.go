package updater

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/apex/log"
	"github.com/racoon-devel/raccoon-pirate/internal/config"
	"github.com/racoon-devel/raccoon-pirate/internal/db"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

type Updater struct {
	CurrentVersion string
	Storage        MetaInfoStorage

	executablePath string
}

func (u *Updater) doSelfUpdate() (bool, error) {
	v, err := ParseVersion(u.CurrentVersion)
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

func (u *Updater) AutoMigration(dbase db.Database, cfg config.Config) error {
	previousVersion, err := u.Storage.GetVersion()
	if err != nil {
		return fmt.Errorf("failed to load previous version: %s", err)
	}

	// just for developing
	if previousVersion == "" || previousVersion == "v0.0.0" || u.CurrentVersion == "v0.0.0" {
		return u.Storage.SetVersion(u.CurrentVersion)
	}
	if previousVersion != u.CurrentVersion {
		// performing storage migration
		log.Warnf("Run migration from previous version...")
		v, err := ParseVersion(previousVersion)
		if err != nil {
			return fmt.Errorf("parse version stored in the metadata failed: %w", err)
		}
		req := migrationRequest{
			major: v.Major,
			minor: v.Minor,
			dbase: dbase,
			cfg:   cfg,
		}
		return u.migrate(&req)
	}
	return nil
}

func (u *Updater) TryUpdate() (updated bool, err error) {
	if u.CurrentVersion == "v0.0.0" {
		return
	}

	u.executablePath, err = os.Executable()
	if err != nil {
		log.Warnf("Get current executable path failed: %s", err)
	}

	updated, err = u.doSelfUpdate()
	return
}

func (u *Updater) Restart() error {
	cmd := exec.Command(u.executablePath, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Start(); err != nil {
		return err
	}

	os.Exit(0)
	return nil
}
