package docs

import (
	"encoding/json"
	"testing"
)

func TestSwaggerJSONPublicOperationsOverrideGlobalSecurity(t *testing.T) {
	var document map[string]any
	if err := json.Unmarshal(SwaggerJSON, &document); err != nil {
		t.Fatalf("failed to parse embedded swagger json: %v", err)
	}

	paths, ok := document["paths"].(map[string]any)
	if !ok {
		t.Fatal("expected paths object in swagger document")
	}

	assertOperationHasEmptySecurity(t, paths, "/api/v1/login", "post")
	assertOperationHasEmptySecurity(t, paths, "/api/v1/register", "post")
	assertOperationHasEmptySecurity(t, paths, "/health", "get")
	assertOperationHasEmptySecurity(t, paths, "/ready", "get")
}

func assertOperationHasEmptySecurity(t *testing.T, paths map[string]any, route string, method string) {
	t.Helper()

	pathItem, ok := paths[route].(map[string]any)
	if !ok {
		t.Fatalf("expected path item for %s", route)
	}

	operation, ok := pathItem[method].(map[string]any)
	if !ok {
		t.Fatalf("expected %s operation for %s", method, route)
	}

	security, ok := operation["security"].([]any)
	if !ok {
		t.Fatalf("expected explicit security override for %s %s", method, route)
	}

	if len(security) != 0 {
		t.Fatalf("expected empty security override for %s %s, got %v", method, route, security)
	}
}
