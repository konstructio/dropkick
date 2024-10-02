package sdk

import (
	"context"
	"errors"
	"strconv"
	"testing"
)

func Test_GetAll(t *testing.T) {
	t.Run("fetch paginated resource", func(t *testing.T) {
		ctx := context.TODO()
		perPage := 2

		// Matrix of paginated responses
		responsePages := [][]Instance{
			{
				{ID: "1", Name: "test-instance-1"},
				{ID: "2", Name: "test-instance-2"},
			},
			{
				{ID: "3", Name: "test-instance-3"},
				{ID: "4", Name: "test-instance-4"},
			},
			{
				{ID: "5", Name: "test-instance-5"},
			},
		}

		c := &MockCivo{
			FnDo: func(ctx context.Context, location, method string, output interface{}, params map[string]string) error {
				// ensure context is being passed
				assertEqual(t, ctx, context.TODO())

				// ensure the appropriate endpoint is being called
				assertEqual(t, location, Instance{}.GetAPIEndpoint())

				// ensure the region is being passed
				assertEqual(t, params["region"], "lon1")

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
				switch pageInt {
				case 1, 2, 3:
					// Response for the expected pages
					resp := paginatedResponse[Instance]{
						Page:    pageInt,
						PerPage: perPage,
						Pages:   3,
						Items:   responsePages[pageInt-1],
					}
					*output.(*paginatedResponse[Instance]) = resp
				default:
					// Civo responds with the first page if you get
					// an overflown page number
					resp := paginatedResponse[Instance]{
						Page:    1,
						PerPage: perPage,
						Pages:   3,
						Items:   responsePages[0],
					}
					*output.(*paginatedResponse[Instance]) = resp
				}

				return nil
			},

			// mock region for all requests
			FnGetRegion: func() string { return "lon1" },
		}

		instances, err := GetAll[Instance](ctx, c)
		assertNoError(t, err)

		// Ensure all instances are returned
		assertEqualf(t, len(instances), 5, "expected 5 instances, got %d", len(instances))

		// Ensure all instances are returned in the correct order
		for i, instance := range instances {
			// Since results are paginated, we need to calculate the index
			// of the instance in the responsePages matrix, we do this by
			// dividing the index by the number of items per page
			// and using the remainder to get the index of the item in the
			// current page.
			currInstance := responsePages[i/perPage][i%perPage]
			assertEqualf(t, instance, currInstance, "expected instance to be %v, got %v", currInstance, instance)
		}
	})

	t.Run("fetch paginated resource with no items", func(t *testing.T) {
		ctx := context.TODO()

		c := &MockCivo{
			FnDo: func(ctx context.Context, location, method string, output interface{}, params map[string]string) error {
				// ensure context is being passed
				assertEqual(t, ctx, context.TODO())

				// ensure the appropriate endpoint is being called
				assertEqual(t, location, Instance{}.GetAPIEndpoint())

				// ensure the region is being passed
				assertEqual(t, params["region"], "lon1")

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
				assertEqualf(t, pageInt, 1, "expected page to be 1, got %v", pageInt)

				// Civo responds with the first page if you get
				// an overflown page number
				resp := paginatedResponse[Instance]{
					Page:    1,
					PerPage: 100,
					Pages:   1,
					Items:   []Instance{},
				}
				*output.(*paginatedResponse[Instance]) = resp

				return nil
			},

			// mock region for all requests
			FnGetRegion: func() string { return "lon1" },
		}

		instances, err := GetAll[Instance](ctx, c)
		assertNoError(t, err)

		// Ensure no instances are returned
		assertEqualf(t, len(instances), 0, "expected 0 instances, got %d", len(instances))
	})

	t.Run("fetch paginated resource with failure on first page", func(t *testing.T) {
		ctx := context.TODO()
		expectedError := errors.New("failure")

		c := &MockCivo{
			FnDo: func(ctx context.Context, location, method string, output interface{}, params map[string]string) error {
				// ensure context is being passed
				assertEqual(t, ctx, context.TODO())

				// ensure the appropriate endpoint is being called
				assertEqual(t, location, Instance{}.GetAPIEndpoint())

				// ensure the region is being passed
				assertEqual(t, params["region"], "lon1")

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
				assertEqualf(t, pageInt, 1, "expected page to be 1, got %v", pageInt)

				// Return an error on the first page
				return expectedError
			},

			// mock region for all requests
			FnGetRegion: func() string { return "lon1" },
		}

		// Ensure the error is the expected error
		_, err := GetAll[Instance](ctx, c)
		assertErrorEqual(t, err, expectedError)
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

		c := &MockCivo{
			FnDo: func(ctx context.Context, location, method string, output interface{}, params map[string]string) error {
				// ensure context is being passed
				assertEqual(t, ctx, context.TODO())

				// ensure the appropriate endpoint is being called
				assertEqual(t, location, Firewall{}.GetAPIEndpoint())

				// ensure the region is being passed
				assertEqual(t, params["region"], "lon1")

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

		firewalls, err := GetAll[Firewall](ctx, c)
		assertNoError(t, err)

		// Ensure all firewalls are returned
		assertEqualf(t, len(firewalls), 6, "expected 6 firewalls, got %d", len(firewalls))

		// Ensure all firewalls are returned in the correct order
		for i, firewall := range firewalls {
			assertEqualf(t, firewall, response[i], "expected firewall to be %v, got %v", response[i], firewall)
		}
	})

	t.Run("fetch single paged resource with no items", func(t *testing.T) {
		ctx := context.TODO()

		c := &MockCivo{
			FnDo: func(ctx context.Context, location, method string, output interface{}, params map[string]string) error {
				// ensure context is being passed
				assertEqual(t, ctx, context.TODO())

				// ensure the appropriate endpoint is being called
				assertEqual(t, location, Volume{}.GetAPIEndpoint())

				// ensure the region is being passed
				assertEqual(t, params["region"], "lon1")

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

		volumes, err := GetAll[Volume](ctx, c)
		assertNoError(t, err)

		// Ensure no volumes are returned
		assertEqualf(t, len(volumes), 0, "expected 0 volumes, got %d", len(volumes))
	})

	t.Run("fetch single paged resource with failure", func(t *testing.T) {
		ctx := context.TODO()
		expectedError := errors.New("failure")

		c := &MockCivo{
			FnDo: func(ctx context.Context, location, method string, output interface{}, params map[string]string) error {
				// ensure context is being passed
				assertEqual(t, ctx, context.TODO())

				// ensure the appropriate endpoint is being called
				assertEqual(t, location, Volume{}.GetAPIEndpoint())

				// ensure the region is being passed
				assertEqual(t, params["region"], "lon1")

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

		_, err := GetAll[Volume](ctx, c)
		assertErrorEqual(t, err, expectedError)
	})
}
