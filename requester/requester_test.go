package requester

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRequester_GET(t *testing.T) {
	t.Run("sendingGetRequestWithoutQueryParams_returnsResources", func(t *testing.T) {
		// ARRANGE
		r := NewRequester(
			URL("https://catfact.ninja"),
		)

		// ACT
		type CatFact struct {
			Fact   string `json:"fact"`
			Length int    `json:"length"`
		}
		var cf CatFact
		resp, body, err := r.Send(&cf, Get("/fact"))

		// ASSERT
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotEmpty(t, body)
		require.NoError(t, resp.Body.Close())
		require.Equal(t, http.StatusOK, resp.StatusCode)

		require.NotEmpty(t, cf.Fact)
		require.NotEmpty(t, cf.Length)

		var jsonCf CatFact
		err = json.Unmarshal(body, &jsonCf)
		require.NoError(t, err)
		require.Equal(t, cf, jsonCf)
	})

	t.Run("sendingGetRequestUsingMultipleQueryParamOption_returnsResources", func(t *testing.T) {
		// ARRANGE
		r := NewRequester(
			URL("https://datausa.io"),
		)

		// ACT
		type Population struct {
			Nation string `json:"Nation"`
		}
		type Measure struct {
			Measures []string `json:"measures"`
		}
		type DataUSA struct {
			Data   []*Population `json:"data"`
			Source []*Measure    `json:"source"`
		}

		var dusa DataUSA
		resp, body, err := r.Send(&dusa,
			Get("/api"),
			AppendPath("/data"),
			QueryParam("drilldowns", "Nation"),
			QueryParam("measures", "Population"),
		)

		// ASSERT
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotEmpty(t, body)
		require.NoError(t, resp.Body.Close())
		require.Equal(t, http.StatusOK, resp.StatusCode)

		require.NotEmpty(t, dusa.Data)
		require.NotEmpty(t, dusa.Source)
		require.NotEmpty(t, dusa.Source[0].Measures)

		require.Equal(t, "United States", dusa.Data[0].Nation)
		require.Equal(t, "Population", dusa.Source[0].Measures[0])
	})

	t.Run("sendingGetRequestUsingASingleQueryParamsOption_returnsResources", func(t *testing.T) {
		// ARRANGE
		r := NewRequester(
			URL("https://datausa.io"),
		)

		// ACT
		type Population struct {
			Nation string `json:"Nation"`
		}
		type Measure struct {
			Measures []string `json:"measures"`
		}
		type DataUSA struct {
			Data   []*Population `json:"data"`
			Source []*Measure    `json:"source"`
		}

		v := url.Values{}
		v.Add("drilldowns", "Nation")
		v.Add("measures", "Population")
		var dusa DataUSA
		resp, body, err := r.Send(&dusa,
			Get("/api"),
			AppendPath("/data"),
			QueryParams(v),
		)

		// ASSERT
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotEmpty(t, body)
		require.NoError(t, resp.Body.Close())
		require.Equal(t, http.StatusOK, resp.StatusCode)

		require.NotEmpty(t, dusa.Data)
		require.NotEmpty(t, dusa.Source)
		require.NotEmpty(t, dusa.Source[0].Measures)

		require.Equal(t, "United States", dusa.Data[0].Nation)
		require.Equal(t, "Population", dusa.Source[0].Measures[0])
	})

	t.Run("sendingGetRequestWithInvalidPath_returnsError", func(t *testing.T) {
		// ARRANGE
		r := NewRequester(
			URL("https://not.existent.path"),
		)

		// ACT
		type Resource struct {
			Data string `json:"data"`
		}

		var resource Resource
		// nolint:bodyclose // resp is expected to be nil, so body cannot be closed
		resp, body, err := r.Send(&resource)

		// ASSERT
		require.Error(t, err)
		require.Nil(t, resp)
		require.Empty(t, body)

		require.Empty(t, resource)
	})
}

func TestRequester_POST(t *testing.T) {
	t.Run("sendingPostRequestWithBody_returnsCreatedResource", func(t *testing.T) {
		// ARRANGE
		r := NewRequester(
			URL("https://reqres.in"),
			ContentType("application/json"),
		)

		// ACT
		type User struct {
			ID        string    `json:"id"`
			Name      string    `json:"name"`
			Job       string    `json:"job"`
			CreatedAt time.Time `json:"createdAt"`
		}
		var user User
		resp, body, err := r.Send(&user,
			Post("/api/users"),
			Body(&struct {
				Name string `json:"name"`
				Job  string `json:"job"`
			}{
				Name: "utils",
				Job:  "developer",
			}),
		)

		// ASSERT
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotEmpty(t, body)
		require.NoError(t, resp.Body.Close())
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		require.NotEmpty(t, user.ID, "ID")
		require.NotEmpty(t, user.Name, "Name")
		require.NotEmpty(t, user.Job, "Job")
		require.NotEmpty(t, user.CreatedAt, "CreatedAt")

		var jsonUser User
		err = json.Unmarshal(body, &jsonUser)
		require.NoError(t, err)
		require.Equal(t, user, jsonUser)
	})

	t.Run("sendingPostRequestWithBodyMissingParameters_returnsStatusBadRequest", func(t *testing.T) {
		// ARRANGE
		r := NewRequester(
			URL("https://reqres.in"),
			ContentType("application/json"),
		)

		// ACT
		type Register struct {
			ID    int    `json:"id"`
			Token string `json:"token"`
		}
		var register Register
		resp, body, err := r.Send(&register,
			Post("/api/register"),
			Body(&struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}{
				Email: "eve.holt@reqres.in",
				// Missing password
			}),
		)

		// ASSERT
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotEmpty(t, body)
		require.NoError(t, resp.Body.Close())
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		require.Empty(t, register.ID, "ID")
		require.Empty(t, register.Token, "Token")
	})
}

func TestRequester_PUT(t *testing.T) {
	t.Run("sendingPutRequestWithBody_returnsCreatedResource", func(t *testing.T) {
		// ARRANGE
		r := NewRequester(
			URL("https://reqres.in"),
			ContentType("application/json"),
		)

		// ACT
		type User struct {
			Name      string    `json:"name"`
			Job       string    `json:"job"`
			UpdatedAt time.Time `json:"updatedAt"`
		}
		var user User
		resp, body, err := r.Send(&user,
			Put("/api/users/2"),
			Body(&struct {
				Name string `json:"name"`
				Job  string `json:"job"`
			}{
				Name: "morpheus",
				Job:  "zion resident",
			}),
		)

		// ASSERT
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotEmpty(t, body)
		require.NoError(t, resp.Body.Close())
		require.Equal(t, http.StatusOK, resp.StatusCode)

		require.NotEmpty(t, user.Name, "Name")
		require.NotEmpty(t, user.Job, "Job")
		require.NotEmpty(t, user.UpdatedAt, "UpdatedAt")

		var jsonUser User
		err = json.Unmarshal(body, &jsonUser)
		require.NoError(t, err)
		require.Equal(t, user, jsonUser)
	})
}

func TestRequester_PATCH(t *testing.T) {
	t.Run("sendingPatchRequestWithBody_returnsCreatedResource", func(t *testing.T) {
		// ARRANGE
		r := NewRequester(
			URL("https://reqres.in"),
			ContentType("application/json"),
		)

		// ACT
		type User struct {
			Name      string    `json:"name"`
			Job       string    `json:"job"`
			UpdatedAt time.Time `json:"updatedAt"`
		}
		var user User
		resp, body, err := r.Send(&user,
			Patch("/api/users/2"),
			Body(&struct {
				Name string `json:"name"`
				Job  string `json:"job"`
			}{
				Name: "morpheus",
				Job:  "zion resident",
			}),
		)

		// ASSERT
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotEmpty(t, body)
		require.NoError(t, resp.Body.Close())
		require.Equal(t, http.StatusOK, resp.StatusCode)

		require.NotEmpty(t, user.Name, "Name")
		require.NotEmpty(t, user.Job, "Job")
		require.NotEmpty(t, user.UpdatedAt, "UpdatedAt")

		var jsonUser User
		err = json.Unmarshal(body, &jsonUser)
		require.NoError(t, err)
		require.Equal(t, user, jsonUser)
	})
}
