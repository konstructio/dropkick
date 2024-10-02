package sdk

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/konstructio/dropkick/internal/civo/sdk/testutils"
)

func Test_GetAll(t *testing.T) {
	t.Run("fetch paginated resource", func(t *testing.T) {
		ctx := context.TODO()
		perPage := 2

		// Matrix of paginated responses
		responses := []Instance{
			{ID: "1", Name: "test-instance-1"},
			{ID: "2", Name: "test-instance-2"},
			{ID: "3", Name: "test-instance-3"},
			{ID: "4", Name: "test-instance-4"},
			{ID: "5", Name: "test-instance-5"},
		}

		c := &testutils.MockCivo{
			FnDo: func(ctx context.Context, location, method string, output interface{}, params map[string]string) error {
				// ensure context is being passed
				testutils.AssertEqual(t, ctx, context.TODO())

				// ensure the appropriate endpoint is being called
				testutils.AssertEqual(t, location, Instance{}.GetAPIEndpoint())

				// ensure the region is being passed
				testutils.AssertEqual(t, params["region"], "lon1")

				// ensure page is being passed
				page, found := params["page"]
				if !found {
					t.Fatalf("expected page to be set, got %v", page)
				}

				// convert the page to an integer
				pageInt, err := strconv.Atoi(page)
				if err != nil {
					t.Fatalf("expected page to be an integer, got %v", page)
				}

				// decide what to return when the page is requested
				res, pageInt, perPage, totalPages := testutils.GetResultsForPage(t, responses, pageInt, perPage)

				// Return the response
				*output.(*PaginatedResponse[Instance]) = PaginatedResponse[Instance]{
					Page:    pageInt,
					PerPage: perPage,
					Pages:   totalPages,
					Items:   res,
				}

				return nil
			},

			// mock region for all requests
			FnGetRegion: func() string { return "lon1" },
		}

		instances, err := getAll[Instance](ctx, c)
		testutils.AssertNoError(t, err)

		// Ensure all instances are returned
		testutils.AssertEqualf(t, len(instances), 5, "expected 5 instances, got %d", len(instances))

		// Ensure all instances are returned in the correct order
		for i, instance := range instances {
			testutils.AssertEqualf(t, instance, responses[i], "expected instance to be %v, got %v", responses[i], instance)
		}
	})

	t.Run("fetch paginated resource with no items", func(t *testing.T) {
		ctx := context.TODO()

		c := &testutils.MockCivo{
			FnDo: func(ctx context.Context, location, method string, output interface{}, params map[string]string) error {
				// ensure context is being passed
				testutils.AssertEqual(t, ctx, context.TODO())

				// ensure the appropriate endpoint is being called
				testutils.AssertEqual(t, location, Instance{}.GetAPIEndpoint())

				// ensure the region is being passed
				testutils.AssertEqual(t, params["region"], "lon1")

				// ensure page is being passed
				page, found := params["page"]
				if !found {
					t.Fatalf("expected page to be set, got %v", page)
				}

				// convert the page to an integer
				pageInt, err := strconv.Atoi(page)
				if err != nil {
					t.Fatalf("expected page to be an integer, got %v", page)
				}

				// Check that you got the first page
				testutils.AssertEqualf(t, pageInt, 1, "expected page to be 1, got %v", pageInt)

				// Civo responds with the first page if you get
				// an overflown page number
				resp := PaginatedResponse[Instance]{
					Page:    1,
					PerPage: 100,
					Pages:   1,
					Items:   []Instance{},
				}
				*output.(*PaginatedResponse[Instance]) = resp

				return nil
			},

			// mock region for all requests
			FnGetRegion: func() string { return "lon1" },
		}

		instances, err := getAll[Instance](ctx, c)
		testutils.AssertNoError(t, err)

		// Ensure no instances are returned
		testutils.AssertEqualf(t, len(instances), 0, "expected 0 instances, got %d", len(instances))
	})

	t.Run("fetch paginated resource with failure on first page", func(t *testing.T) {
		ctx := context.TODO()
		expectedError := errors.New("failure")

		c := &testutils.MockCivo{
			FnDo: func(ctx context.Context, location, method string, output interface{}, params map[string]string) error {
				// ensure context is being passed
				testutils.AssertEqual(t, ctx, context.TODO())

				// ensure the appropriate endpoint is being called
				testutils.AssertEqual(t, location, Instance{}.GetAPIEndpoint())

				// ensure the region is being passed
				testutils.AssertEqual(t, params["region"], "lon1")

				// ensure page is being passed
				page, found := params["page"]
				if !found {
					t.Fatalf("expected page to be set, got %v", page)
				}

				// convert the page to an integer
				pageInt, err := strconv.Atoi(page)
				if err != nil {
					t.Fatalf("expected page to be an integer, got %v", page)
				}

				// Check that you got the first page
				testutils.AssertEqualf(t, pageInt, 1, "expected page to be 1, got %v", pageInt)

				// Return an error on the first page
				return expectedError
			},

			// mock region for all requests
			FnGetRegion: func() string { return "lon1" },
		}

		// Ensure the error is the expected error
		_, err := getAll[Instance](ctx, c)
		testutils.AssertErrorEqual(t, err, expectedError)
	})

	t.Run("fetch single paged resource", func(t *testing.T) {
		ctx := context.TODO()

		response := []Firewall{
			{ID: "1", Name: "test-firewall-1"},
			{ID: "2", Name: "test-firewall-2"},
			{ID: "3", Name: "test-firewall-3"},
			{ID: "4", Name: "test-firewall-4"},
			{ID: "5", Name: "test-firewall-5"},
			{ID: "6", Name: "test-firewall-6"},
		}

		c := &testutils.MockCivo{
			FnDo: func(ctx context.Context, location, method string, output interface{}, params map[string]string) error {
				// ensure context is being passed
				testutils.AssertEqual(t, ctx, context.TODO())

				// ensure the appropriate endpoint is being called
				testutils.AssertEqual(t, location, Firewall{}.GetAPIEndpoint())

				// ensure the region is being passed
				testutils.AssertEqual(t, params["region"], "lon1")

				// ensure page is not being passed
				_, found := params["page"]
				if found {
					t.Fatalf("expected page not to be set, got %v", params["page"])
				}

				// Return the response
				*output.(*[]Firewall) = response

				return nil
			},

			// mock region for all requests
			FnGetRegion: func() string { return "lon1" },
		}

		firewalls, err := getAll[Firewall](ctx, c)
		testutils.AssertNoError(t, err)

		// Ensure all firewalls are returned
		testutils.AssertEqualf(t, len(firewalls), 6, "expected 6 firewalls, got %d", len(firewalls))

		// Ensure all firewalls are returned in the correct order
		for i, firewall := range firewalls {
			testutils.AssertEqualf(t, firewall, response[i], "expected firewall to be %v, got %v", response[i], firewall)
		}
	})

	t.Run("fetch single paged resource with no items", func(t *testing.T) {
		ctx := context.TODO()

		c := &testutils.MockCivo{
			FnDo: func(ctx context.Context, location, method string, output interface{}, params map[string]string) error {
				// ensure context is being passed
				testutils.AssertEqual(t, ctx, context.TODO())

				// ensure the appropriate endpoint is being called
				testutils.AssertEqual(t, location, Volume{}.GetAPIEndpoint())

				// ensure the region is being passed
				testutils.AssertEqual(t, params["region"], "lon1")

				// ensure page is not being passed
				_, found := params["page"]
				if found {
					t.Fatalf("expected page not to be set, got %v", params["page"])
				}

				// Return an empty response
				*output.(*[]Volume) = []Volume{}

				return nil
			},

			// mock region for all requests
			FnGetRegion: func() string { return "lon1" },
		}

		volumes, err := getAll[Volume](ctx, c)
		testutils.AssertNoError(t, err)

		// Ensure no volumes are returned
		testutils.AssertEqualf(t, len(volumes), 0, "expected 0 volumes, got %d", len(volumes))
	})

	t.Run("fetch single paged resource with failure", func(t *testing.T) {
		ctx := context.TODO()
		expectedError := errors.New("failure")

		c := &testutils.MockCivo{
			FnDo: func(ctx context.Context, location, method string, output interface{}, params map[string]string) error {
				// ensure context is being passed
				testutils.AssertEqual(t, ctx, context.TODO())

				// ensure the appropriate endpoint is being called
				testutils.AssertEqual(t, location, Volume{}.GetAPIEndpoint())

				// ensure the region is being passed
				testutils.AssertEqual(t, params["region"], "lon1")

				// ensure page is not being passed
				_, found := params["page"]
				if found {
					t.Fatalf("expected page not to be set, got %v", params["page"])
				}

				// Return the error
				return expectedError
			},

			// mock region for all requests
			FnGetRegion: func() string { return "lon1" },
		}

		_, err := getAll[Volume](ctx, c)
		testutils.AssertErrorEqual(t, err, expectedError)
	})
}
