package utils

import (
	"smartgas-payment/internal/models"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

type TrackErrorOpts struct {
	Extras      map[string]any
	Context     map[string]map[string]any
	Tags        map[string]string
	Customer    *models.Customer
	Admin       *models.User
	Application *models.AuthorizedApplication
	Level       sentry.Level
}

func TrackError(c *gin.Context, err error, opts *TrackErrorOpts) {
	if hub := sentrygin.GetHubFromContext(c); hub != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			if opts != nil {
				scope.SetExtras(opts.Extras)
				scope.SetTags(opts.Tags)
				scope.SetContexts(opts.Context)

				if opts.Level != "" {
					scope.SetLevel(opts.Level)
				} else {
					scope.SetLevel(sentry.LevelFatal)
				}
				if opts.Customer != nil {
					scope.SetUser(sentry.User{
						ID:    opts.Customer.ID.String(),
						Email: opts.Customer.Email,
					})
				}
				if opts.Admin != nil {
					scope.SetUser(sentry.User{
						ID:    opts.Admin.ID.String(),
						Email: opts.Admin.Email,
					})
				}

				if opts.Application != nil {
					scope.SetUser(sentry.User{
						ID:   opts.Application.ID.String(),
						Name: opts.Application.ApplicationName,
					})
				}
			}
			hub.CaptureException(err)
		})
	}
}
