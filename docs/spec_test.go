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

	if _, ok := document["security"]; ok {
		t.Fatal("expected no global security requirement in swagger document")
	}

	paths, ok := document["paths"].(map[string]any)
	if !ok {
		t.Fatal("expected paths object in swagger document")
	}

	assertOperationHasNoSecurity(t, paths, "/api/v1/login", "post")
	assertOperationHasNoSecurity(t, paths, "/api/v1/register", "post")
	assertOperationHasNoSecurity(t, paths, "/health", "get")
	assertOperationHasNoSecurity(t, paths, "/ready", "get")
	assertOperationHasBearerSecurity(t, paths, "/api/v1/transactions", "get")
	assertOperationHasBearerSecurity(t, paths, "/api/v1/budgets", "get")
}

func assertOperationHasNoSecurity(t *testing.T, paths map[string]any, route string, method string) {
	t.Helper()

	pathItem, ok := paths[route].(map[string]any)
	if !ok {
		t.Fatalf("expected path item for %s", route)
	}

	operation, ok := pathItem[method].(map[string]any)
	if !ok {
		t.Fatalf("expected %s operation for %s", method, route)
	}

	if _, ok := operation["security"]; ok {
		t.Fatalf("expected no explicit security for %s %s", method, route)
	}
}

func assertOperationHasBearerSecurity(t *testing.T, paths map[string]any, route string, method string) {
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
	if !ok || len(security) != 1 {
		t.Fatalf("expected bearer security for %s %s, got %v", method, route, operation["security"])
	}

	requirement, ok := security[0].(map[string]any)
	if !ok {
		t.Fatalf("expected security requirement object for %s %s", method, route)
	}

	scopes, ok := requirement["BearerAuth"].([]any)
	if !ok || len(scopes) != 0 {
		t.Fatalf("expected BearerAuth with no scopes for %s %s, got %v", method, route, requirement)
	}
}
