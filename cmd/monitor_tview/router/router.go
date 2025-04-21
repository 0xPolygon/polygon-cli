package router

import (
	"github.com/rivo/tview"
)

type Router struct {
	app         *tview.Application
	pages       *tview.Pages
	screens     map[string]func(id string) tview.Primitive
	current     string
	lastFocused tview.Primitive
}

func NewRouter(app *tview.Application, pages *tview.Pages) *Router {
	return &Router{
		app:     app,
		pages:   pages,
		screens: make(map[string]func(string) tview.Primitive),
	}
}

func (r *Router) Current() string {
	return r.current
}

func (r *Router) Pages() *tview.Pages {
	return r.pages
}

func (r *Router) RegisterScreen(screen string, factory func(id string) tview.Primitive) {
	r.screens[screen] = factory
}

func (r *Router) RegisterModal(screen string, modal tview.Primitive) {
	r.pages.AddPage(screen, modal, true, false)
}

func (r *Router) Navigate(screen string, id string) {
	factory, ok := r.screens[screen]
	if !ok {
		return
	}
	view := factory(id)

	if r.pages.HasPage(screen) {
		r.pages.RemovePage(screen)
	}
	r.pages.AddPage(screen, view, true, true)
	r.app.SetFocus(view)
	r.lastFocused = view
	r.current = screen
}

func (r *Router) LastFocused() tview.Primitive {
	return r.lastFocused
}

func (r *Router) ShowModal(screen string) {
	if r.pages.HasPage(screen) {
		r.pages.ShowPage(screen)
	}
}

func (r *Router) HideModal(screen string) {
	if r.pages.HasPage(screen) {
		r.pages.HidePage(screen)
	}
}
