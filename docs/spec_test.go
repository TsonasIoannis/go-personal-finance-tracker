package docs

import (
	"encoding/json"
	"testing"
)

func TestSwaggerJSONSecurityShape(t *testing.T) {
	var document map[string]any
	if err := json.Unmarshal(SwaggerJSON, &document); err != nil {
		t.Fatalf("failed to parse embedded swagger json: %v", err)
	}

	security, ok := document["security"].([]any)
	if !ok || len(security) != 1 {
		t.Fatalf("expected one global security requirement, got %v", document["security"])
	}

	paths, ok := document["paths"].(map[string]any)
	if !ok {
		t.Fatal("expected paths object in swagger document")
	}

	assertOperationHasAnonymousOverride(t, paths, "/api/v1/login", "post")
	assertOperationHasAnonymousOverride(t, paths, "/api/v1/register", "post")
	assertOperationHasAnonymousOverride(t, paths, "/health", "get")
	assertOperationHasAnonymousOverride(t, paths, "/ready", "get")
	assertOperationHasBearerSecurity(t, paths, "/api/v1/transactions", "get")
	assertOperationHasBearerSecurity(t, paths, "/api/v1/budgets", "get")
}

func assertOperationHasAnonymousOverride(t *testing.T, paths map[string]any, route string, method string) {
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
		t.Fatalf("expected anonymous override for %s %s, got %v", method, route, operation["security"])
	}

	requirement, ok := security[0].(map[string]any)
	if !ok {
		t.Fatalf("expected security requirement object for %s %s", method, route)
	}

	if len(requirement) != 0 {
		t.Fatalf("expected empty security requirement object for %s %s, got %v", method, route, requirement)
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
