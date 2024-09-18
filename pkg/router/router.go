// Interface for a page
package router

import (
	"errors"

	"github.com/anti-raid/evil-befall/pkg/state"
	"github.com/rivo/tview"
)

var routes = []Route{}

func AddRoute(r Route) {
	routes = append(routes, r)
}

func GetRoute(id string) Route {
	for _, r := range routes {
		if r.ID() == id {
			return r
		}
	}

	return nil
}

// Get current route
func GetCurrentRoute(state *state.State) Route {
	for _, route := range routes {
		if state.CurrentLoc.ID == route.ID() {
			return route
		}
	}

	return nil
}

// Goto page based on the states current location
func GotoCurrent(state *state.State, app *tview.Application, pages *tview.Pages) (Route, error) {
	for _, route := range routes {
		if state.CurrentLoc.ID == route.ID() {
			return Goto(state, route.ID(), app, pages)
		}
	}

	return nil, errors.New("route not found")
}

// Goto page by ID
func Goto(state *state.State, id string, app *tview.Application, pages *tview.Pages) (Route, error) {
	// Get current route
	currentRoute := GetCurrentRoute(state)

	// If the current route is not nil, destroy it
	if currentRoute != nil {
		if err := currentRoute.Destroy(state); err != nil {
			return nil, err
		}
	}

	var route = GetRoute(id)

	if route == nil {
		return nil, errors.New("route not found")
	}

	if err := route.Setup(state); err != nil {
		return nil, err
	}

	page, err := route.Render(state, app, pages)

	if err != nil {
		return nil, err
	}

	pages.AddAndSwitchToPage(id, page, true)
	app.SetFocus(page)

	return route, nil
}

type Route interface {
	// The ID of the route
	ID() string

	// The title of the route
	Title() string

	// The description of the route
	Description() string

	// Given a current state, sets up all state for the route
	Setup(state *state.State) error

	// Called on destruction of the route
	Destroy(state *state.State) error

	// Renders the route returning a tview.Primitive
	Render(state *state.State, app *tview.Application, pages *tview.Pages) (tview.Primitive, error)
}
