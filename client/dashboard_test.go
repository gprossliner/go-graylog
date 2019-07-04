package client_test

import (
	"context"
	"testing"

	"github.com/suzuki-shunsuke/go-graylog/testutil"
)

func TestCreateDashboard(t *testing.T) {
	ctx := context.Background()
	server, client, err := testutil.GetServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	if server != nil {
		defer server.Close()
	}

	// nil check
	if _, err := client.CreateDashboard(ctx, nil); err == nil {
		t.Fatal("dashboard is nil")
	}
	// success
	dashboard := testutil.Dashboard()
	if _, err := client.CreateDashboard(ctx, dashboard); err != nil {
		t.Fatal(err)
	}
	// clean
	defer client.DeleteDashboard(ctx, dashboard.ID)
}

func TestDeleteDashboard(t *testing.T) {
	ctx := context.Background()
	server, client, err := testutil.GetServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	if server != nil {
		defer server.Close()
	}

	// id required
	if _, err := client.DeleteDashboard(ctx, ""); err == nil {
		t.Fatal("id is required")
	}
	// invalid id
	if _, err := client.DeleteDashboard(ctx, "h"); err == nil {
		t.Fatal(`no dashboard with id "h" is found`)
	}
}

func TestGetDashboard(t *testing.T) {
	ctx := context.Background()
	server, client, err := testutil.GetServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	if server != nil {
		defer server.Close()
	}

	dashboard := testutil.Dashboard()
	if _, err := client.CreateDashboard(ctx, dashboard); err != nil {
		t.Fatal(err)
	}
	defer client.DeleteDashboard(ctx, dashboard.ID)

	r, _, err := client.GetDashboard(ctx, dashboard.ID)
	if err != nil {
		t.Fatal(err)
	}
	if r == nil {
		t.Fatal("dashboard is nil")
	}
	if r.ID != dashboard.ID {
		t.Fatalf(`dashboard.ID = "%s", wanted "%s"`, r.ID, dashboard.ID)
	}
	if _, _, err := client.GetDashboard(ctx, ""); err == nil {
		t.Fatal("id is required")
	}
	if _, _, err := client.GetDashboard(ctx, "h"); err == nil {
		t.Fatal("dashboard should not be found")
	}
}

func TestGetDashboards(t *testing.T) {
	ctx := context.Background()
	server, client, err := testutil.GetServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	if server != nil {
		defer server.Close()
	}

	if _, _, _, err := client.GetDashboards(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateDashboard(t *testing.T) {
	ctx := context.Background()
	server, client, err := testutil.GetServerAndClient()
	if err != nil {
		t.Fatal(err)
	}
	if server != nil {
		defer server.Close()
	}
	dashboard := testutil.Dashboard()
	if _, err := client.CreateDashboard(ctx, dashboard); err != nil {
		t.Fatal(err)
	}
	// clean
	defer client.DeleteDashboard(ctx, dashboard.ID)

	dashboard.Description = "changed!"
	if _, err := client.UpdateDashboard(ctx, dashboard); err != nil {
		t.Fatal(err)
	}
	dashboard.ID = ""
	if _, err := client.UpdateDashboard(ctx, dashboard); err == nil {
		t.Fatal("id is required")
	}
	dashboard.ID = "h"
	if _, err := client.UpdateDashboard(ctx, dashboard); err == nil {
		t.Fatal(`no dashboard whose id is "h"`)
	}
	if _, err := client.UpdateDashboard(ctx, nil); err == nil {
		t.Fatal("dashboard is nil")
	}
}
