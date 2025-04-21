package monitor_tview

import (
	"github.com/0xPolygon/polygon-cli/cmd/monitor_tview/router"
	"github.com/0xPolygon/polygon-cli/cmd/monitor_tview/ui/views"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type MonitorApp struct {
	app    *tview.Application
	router *router.Router
	rpcURL string
}

func NewMonitorApp(rpcURL string) *MonitorApp {
	app := tview.NewApplication()
	pages := tview.NewPages()
	router := router.NewRouter(app, pages)

	return &MonitorApp{
		app:    app,
		router: router,
		rpcURL: rpcURL,
	}
}

func (m *MonitorApp) Run() error {
	m.setGlobalKeybindings()
	m.RegisterModals()
	m.RegisterScreens()
	m.router.Navigate("network", m.rpcURL)
	return m.app.SetRoot(m.router.Pages(), true).EnableMouse(true).Run()
}

func (m *MonitorApp) setGlobalKeybindings() {
	m.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case '?':
				m.router.ShowModal("help")
				return nil
			case 'q':
				m.app.Stop()
				return nil
			}
		}
		return event
	})
}

func (m *MonitorApp) RegisterModals() {
	// Create modal once
	helpModal := views.NewHelpModal(m.app, func() {
		m.router.HideModal("help")
		m.app.SetFocus(m.router.LastFocused())
	})
	m.router.RegisterModal("help", helpModal)
}

func (m *MonitorApp) RegisterScreens() {
	m.router.RegisterScreen("network", func(id string) tview.Primitive {
		return views.NewNetworkView(m.app, m.router, id)
	})
	m.router.RegisterScreen("block", func(id string) tview.Primitive {
		return views.NewBlockView(m.app, m.router, id)
	})
	m.router.RegisterScreen("tx", func(id string) tview.Primitive {
		return views.NewTxView(m.app, m.router, id)
	})
}
