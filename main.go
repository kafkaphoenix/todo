package main

import (
	"net/http"
	"os"

	"github.com/calvinmclean/babyapi"
	"github.com/calvinmclean/babyapi/extensions"
	"github.com/calvinmclean/babyapi/html"

	"github.com/go-chi/render"
)

type TODO struct {
	babyapi.DefaultResource

	Title       string
	Description string
	Completed   *bool
}

func (t *TODO) HTML(r *http.Request) string {
	return todoRow.Render(r, t)
}

type AllTODOs struct {
	babyapi.ResourceList[*TODO]
}

func (at AllTODOs) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (at AllTODOs) HTML(r *http.Request) string {
	return allTODOs.Render(r, at.Items)
}

func createAPI() *babyapi.API[*TODO] {
	api := babyapi.NewAPI("TODOs", "/todos", func() *TODO { return &TODO{} })

	api.AddCustomRootRoute(http.MethodGet, "/", http.RedirectHandler("/todos", http.StatusFound))

	// Use AllTODOs in the GetAll response since it implements HTMLer
	api.SetGetAllResponseWrapper(func(todos []*TODO) render.Renderer {
		return AllTODOs{ResourceList: babyapi.ResourceList[*TODO]{todos}}
	})

	api.ApplyExtension(extensions.HTMX[*TODO]{})

	// Add SSE handler endpoint which will receive events on the returned channel and write them to the front-end
	todoChan := api.AddServerSentEventHandler("/listen")

	// Push events onto the SSE channel when new TODOs are created
	api.SetOnCreateOrUpdate(func(_ http.ResponseWriter, r *http.Request, t *TODO) *babyapi.ErrResponse {
		if r.Method != http.MethodPost {
			return nil
		}

		select {
		case todoChan <- &babyapi.ServerSentEvent{Event: "newTODO", Data: t.HTML(r)}:
		default:
			logger := babyapi.GetLoggerFromContext(r.Context())
			logger.Info("no listeners for server-sent event")
		}
		return nil
	})

	// Optionally setup redis storage if environment variables are defined
	api.ApplyExtension(extensions.KeyValueStorage[*TODO]{
		KVConnectionConfig: extensions.KVConnectionConfig{
			RedisHost:     os.Getenv("REDIS_HOST"),
			RedisPassword: os.Getenv("REDIS_PASS"),
			Filename:      os.Getenv("STORAGE_FILE"),
			Optional:      true,
		},
	})

	html.SetMap(map[string]string{
		string(allTODOs): allTODOsTemplate,
		string(todoRow):  todoRowTemplate,
	})

	return api
}

func main() {
	api := createAPI()
	api.RunCLI()
}
