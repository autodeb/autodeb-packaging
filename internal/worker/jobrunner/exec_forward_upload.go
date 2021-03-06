package jobrunner

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
	"salsa.debian.org/autodeb-team/autodeb/internal/exec"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

func (jobRunner *JobRunner) execForwardUpload(
	ctx context.Context,
	job *models.Job,
	workingDirectory string,
	artifactsDirectory string,
	logFile io.Writer) error {

	if job.ParentType != models.JobParentTypeUpload {
		return errors.Errorf("unsupported parent type %s", job.ParentType)
	}

	// Retrieve the corresponding Upload
	upload, err := jobRunner.apiClient.GetUpload(job.ParentID)
	if err != nil {
		return errors.WithMessage(err, "could not retrieve corresponding upload")
	}

	// Download the upload
	changesURL := jobRunner.apiClient.GetUploadChangesURL(job.ParentID)
	if err := exec.RunCtxDirStdoutStderr(
		ctx, workingDirectory, logFile, logFile,
		"dget", "--allow-unauthenticated", changesURL.String(),
	); err != nil {
		return errors.WithMessage(err, "dget failed")
	}

	// Rename the changes file
	changesFileName := path.Base(changesURL.EscapedPath())
	newChangesFileName := fmt.Sprintf("%s_%s_.source.changes", upload.Source, upload.Version)
	newChangesPath := filepath.Join(workingDirectory, newChangesFileName)
	if err := os.Rename(
		filepath.Join(workingDirectory, changesFileName),
		newChangesPath,
	); err != nil {
		return errors.WithMessagef(err, "could not rename changes file to %s", newChangesFileName)
	}

	// Run dput
	if err := exec.RunCtxDirStdoutStderr(
		ctx, workingDirectory, logFile, logFile,
		"dput", "--unchecked", newChangesPath,
	); err != nil {
		return errors.WithMessage(err, "dput failed")
	}

	return nil
}
