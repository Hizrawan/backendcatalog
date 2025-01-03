package reqdata

import (
	"context"

	"github.com/xinchuantw/hoki-tabloid-backend/internal/app"
)

// Context stores information about the current app-related request context for easier access instead of
// having to retrieve them from the request context.Context directly.
type Context struct {
	App         *app.Registry
	Auth        AuthInformation
	HTTPContext context.Context
}
