package jobrunner

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
	"salsa.debian.org/autodeb-team/autodeb/internal/ftpmasterapi"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"

	"pault.ag/go/debian/version"
)

func (jobRunner *JobRunner) execSetupArchiveBackport(
	ctx context.Context,
	job *models.Job,
	workingDirectory string,
	artifactsDirectory string,
	logFile io.Writer) error {

	if job.ParentType != models.JobParentTypeArchiveBackport {
		return errors.Errorf("unsupported parent type %s", job.ParentType)
	}

	ftpmasterapiClient := ftpmasterapi.NewClient(http.DefaultClient)

	// Get stable packages
	stablePackages := make(map[string]*ftpmasterapi.Source)
	pkgs, err := ftpmasterapiClient.GetSourcesInSuite("stable")
	if err != nil {
		return errors.WithMessage(err, "could not get unstable sources")
	}
	for _, pkg := range pkgs {
		stablePackages[pkg.Source] = pkg
	}

	// Get testing packages
	testingPackages := make(map[string]*ftpmasterapi.Source)
	pkgs, err = ftpmasterapiClient.GetSourcesInSuite("testing")
	if err != nil {
		return errors.WithMessage(err, "could not get testing sources")
	}
	for _, pkg := range pkgs {
		testingPackages[pkg.Source] = pkg
	}

	// Find backport candidates:
	// - stable < testing
	// - in testing but not in stable
	//
	// TODO: check dependencies and don't bother trying to build
	// packages where the dependencies are not satisfied in the
	// target suite.
	//
	var backportCandidates []string
	for _, testingPackage := range testingPackages {

		stablePackage, ok := stablePackages[testingPackage.Source]
		if !ok {
			backportCandidates = append(backportCandidates, testingPackage.Source)
			fmt.Fprintf(logFile, "Adding %s to backport candidates: not found in stable\n", testingPackage.Source)
			continue
		}

		testingVersion, err := version.Parse(testingPackage.Version)
		if err != nil {
			return errors.WithMessagef(err, "could not parse version %s", testingPackage.Version)
		}

		stableVersion, err := version.Parse(stablePackage.Version)
		if err != nil {
			return errors.WithMessagef(err, "could not parse version %s", stablePackage.Version)
		}

		if version.Compare(stableVersion, testingVersion) < 0 {
			backportCandidates = append(backportCandidates, testingPackage.Source)
			fmt.Fprintf(
				logFile,
				"Adding %s to backport candidates: %s (stable) < %s (testing)\n",
				testingPackage.Source,
				stableVersion.String(),
				testingVersion.String(),
			)
		}

	}

	// Create Upgrade jobs
	for _, backportCandidate := range backportCandidates {
		if _, err := jobRunner.apiClient.CreateJob(
			&models.Job{
				Type:       models.JobTypeBackport,
				ParentType: job.ParentType,
				ParentID:   job.ParentID,
				Input:      backportCandidate,
			},
		); err != nil {
			return errors.WithMessagef(err, "could not create backport job for package %s", backportCandidate)
		}
		fmt.Fprintf(
			logFile,
			"Created backport job for source package %s\n",
			backportCandidate,
		)
	}

	return nil
}
