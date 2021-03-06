package webpages

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"salsa.debian.org/autodeb-team/autodeb/internal/http/middleware"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/appctx"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/handlers/middleware/auth"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

//ArchiveUpgradesGetHandler returns a handler that renders the archive upgrades page
func ArchiveUpgradesGetHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request, user *models.User) {

		page := 0
		limit := 30
		if pageParam := r.URL.Query().Get("page"); pageParam != "" {
			page, _ = strconv.Atoi(pageParam)
		}

		archiveUpgrades, err := appCtx.JobsService().GetAllArchiveUpgradesPageLimit(page, limit)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			appCtx.RequestLogger().Error(r, err)
			return
		}

		data := struct {
			ArchiveUpgrades []*models.ArchiveUpgrade
			PreviousPage    int
			NextPage        int
		}{
			ArchiveUpgrades: archiveUpgrades,
			PreviousPage:    page - 1,
			NextPage:        page + 1,
		}

		renderWithBase(r, w, appCtx, user, "archive_upgrades.html", data)
	}

	handler := auth.MaybeWithUser(handlerFunc, appCtx)
	handler = middleware.HTMLHeaders(handler)

	return handler
}

// ArchiveUpgradeGetHandler returns a handler that renders the ArchiveUpgrade details page
func ArchiveUpgradeGetHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request, user *models.User) {

		vars := mux.Vars(r)
		archiveUpgradeID, err := strconv.Atoi(vars["archiveUpgradeID"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			appCtx.RequestLogger().Error(r, err)
			return
		}

		archiveUpgrade, err := appCtx.JobsService().GetArchiveUpgrade(uint(archiveUpgradeID))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			appCtx.RequestLogger().Error(r, err)
			return
		}
		if archiveUpgrade == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		jobs, err := appCtx.JobsService().GetAllJobsByArchiveUpgradeID(archiveUpgrade.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			appCtx.RequestLogger().Error(r, err)
			return
		}

		repositoryURL := fmt.Sprintf(
			"deb [trusted=yes] %s/%s unstable main",
			strings.Trim(
				appCtx.Config().Aptly.RepositoryBaseURL.String(),
				"/",
			),
			archiveUpgrade.RepositoryName(),
		)

		data := struct {
			ArchiveUpgrade *models.ArchiveUpgrade
			Jobs           []*models.Job
			RepositoryURL  string
		}{
			ArchiveUpgrade: archiveUpgrade,
			Jobs:           jobs,
			RepositoryURL:  repositoryURL,
		}

		renderWithBase(r, w, appCtx, user, "archive_upgrade.html", data)
	}

	handler := auth.MaybeWithUser(handlerFunc, appCtx)
	handler = middleware.HTMLHeaders(handler)

	return handler
}

// NewArchiveUpgradeGetHandler returns a handler that renders the page to create a new archive upgrade
func NewArchiveUpgradeGetHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request, user *models.User) {
		renderWithBase(r, w, appCtx, user, "archive_upgrade_new.html", nil)
	}

	handler := auth.WithUserOrRedirect(handlerFunc, appCtx)
	handler = middleware.HTMLHeaders(handler)

	return handler
}

// NewArchiveUpgradePostHandler returns a handle that creates a new archive upgrade
func NewArchiveUpgradePostHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request, user *models.User) {

		r.ParseForm()
		packageCount := r.Form.Get("package-count")

		packageCountInt, err := strconv.Atoi(packageCount)
		if err != nil {
			appCtx.Sessions().Flash(r, w, "danger", "invalid package count")
		} else if _, err := appCtx.JobsService().CreateArchiveUpgrade(user.ID, packageCountInt); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			appCtx.RequestLogger().Error(r, err)
		} else {
			appCtx.Sessions().Flash(r, w, "success", "Archive upgrade created successfully")
		}

		http.Redirect(w, r, "/new-archive-upgrade", http.StatusSeeOther)

	}

	handler := auth.WithUserOrRedirect(handlerFunc, appCtx)

	return handler
}
