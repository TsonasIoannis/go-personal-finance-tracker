package docs

import _ "embed"

// SwaggerJSON stores the generated OpenAPI document served by the app.
//
//go:embed swagger.json
var SwaggerJSON []byte
