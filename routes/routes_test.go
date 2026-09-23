package routes_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"survey-backend/internal/user"
	"survey-backend/routes"

	"github.com/gin-gonic/gin"
)

// buildRouter mirrors cmd/main.go's wiring but swaps the real auth middleware
// for one that injects a user of the given role, so the role gates can be
// exercised without a database or a token.
func buildRouter(role string) *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(gin.Recovery())

	protected := r.Group("/")
	protected.Use(func(c *gin.Context) {
		c.Set("currentUser", user.User{Role: role})
		c.Next()
	})

	routes.UserRoutes(protected)
	routes.RespondentRoutes(protected)
	routes.InterviewerRoutes(protected)
	routes.AdminRoutes(protected)

	return r
}

func status(t *testing.T, r *gin.Engine, method, path string) int {
	t.Helper()

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(method, path, nil))
	return w.Code
}

// TestRoutesRegister fails if gin panics on a path conflict, which it does when
// static and wildcard segments are registered incompatibly.
func TestRoutesRegister(t *testing.T) {
	r := buildRouter("interviewer")

	want := []string{
		"GET /me",
		"GET /regions/:country_id",
		"GET /interviewer/home",
		"GET /interviewer/completed",
		"GET /interviewer/survey/form",
		"POST /interviewer/survey/create",
		"GET /interviewer/survey/list",
		"GET /interviewer/survey/count",
		"GET /interviewer/survey/:id",
		"PUT /interviewer/survey/:id",
		"DELETE /interviewer/survey/:id",
		"POST /interviewer/survey/:id/publish",
		"POST /interviewer/survey/:id/pause",
		"GET /interviewer/survey/:id/analytics",
		"GET /interviewer/survey/:id/questions",
		"POST /interviewer/survey/:id/questions",
		"PUT /interviewer/survey/:id/questions/:question_id",
		"DELETE /interviewer/survey/:id/questions/:question_id",
		"GET /respondent/home",
		"GET /respondent/survey/list",
		"GET /respondent/survey/:id",
		"POST /respondent/survey/:id/answer",
		"GET /respondent/survey/:id/answer-count",
		"GET /respondent/saved",
		"GET /respondent/saved-ids",
		"POST /respondent/saved/:id",
		"DELETE /respondent/saved/:id",
		"GET /respondent/completed",
		"GET /respondent/completed/:id",
		"POST /admin/country/toggle",
		"POST /admin/points/grant",
	}

	registered := map[string]bool{}
	for _, route := range r.Routes() {
		registered[route.Method+" "+route.Path] = true
	}

	for _, route := range want {
		if !registered[route] {
			t.Errorf("route not registered: %s", route)
		}
	}
}

// TestRoleGatesApply checks each group rejects the wrong role.
func TestRoleGatesApply(t *testing.T) {
	cases := []struct {
		role, method, path string
	}{
		{"respondent", http.MethodGet, "/interviewer/home"},
		{"interviewer", http.MethodGet, "/respondent/home"},
		{"interviewer", http.MethodPost, "/admin/country/toggle"},
		{"respondent", http.MethodPost, "/admin/country/toggle"},
		{"interviewer", http.MethodPost, "/admin/points/grant"},
		{"respondent", http.MethodPost, "/admin/points/grant"},
	}

	for _, tc := range cases {
		got := status(t, buildRouter(tc.role), tc.method, tc.path)
		if got != http.StatusForbidden {
			t.Errorf("%s as %s: got %d, want 403", tc.path, tc.role, got)
		}
	}
}

// TestRoleMiddlewareDoesNotLeak guards the bug where every route file called
// Use on the same shared parent group: the role check was appended to the
// parent's chain, so groups registered afterwards required every earlier role
// too, making /interviewer/... unreachable for an interviewer.
//
// These requests reach their handler and then fail on the nil test database,
// which Recovery turns into a 500. A 403 here means a foreign role gate ran.
func TestRoleMiddlewareDoesNotLeak(t *testing.T) {
	cases := []struct {
		role, method, path string
	}{
		{"interviewer", http.MethodGet, "/interviewer/home"},
		{"interviewer", http.MethodGet, "/interviewer/survey/list"},
		{"respondent", http.MethodGet, "/respondent/home"},
		{"respondent", http.MethodGet, "/respondent/saved-ids"},
		{"respondent", http.MethodGet, "/me"},
	}

	for _, tc := range cases {
		got := status(t, buildRouter(tc.role), tc.method, tc.path)
		if got == http.StatusForbidden {
			t.Errorf("%s as %s: got 403, a role gate from another group ran", tc.path, tc.role)
		}
	}
}

// TestAdminReachesEveryGroup documents that RoleMiddleware lets admins through
// any role gate.
func TestAdminReachesEveryGroup(t *testing.T) {
	for _, path := range []string{"/interviewer/home", "/respondent/home"} {
		if got := status(t, buildRouter("admin"), http.MethodGet, path); got == http.StatusForbidden {
			t.Errorf("%s as admin: got 403, want the request to pass the role gate", path)
		}
	}
}
