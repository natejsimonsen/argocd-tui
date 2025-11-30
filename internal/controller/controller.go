package controller

import (
	"example.com/main/internal/model"
	"example.com/main/internal/view"
	"example.com/main/services/utils"
	"github.com/gdamore/tcell/v2"
)

type AppController struct {
	Model *model.AppModel
	View  *view.AppView
}

func NewAppController(m *model.AppModel, v *view.AppView) *AppController {
	return &AppController{
		Model: m,
		View:  v,
	}
}

func (c *AppController) SetupEventHandlers() {
	// application on change
	c.View.AppList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		appName := utils.StripTags(mainText)

		c.Model.LoadResources(appName)

		c.View.UpdateMainContent(c.Model.SelectedAppResources)
		c.View.UpdateTitles(index, c.Model.PrevIndex, appName, c.Model.PrevText)

		c.Model.PrevText = mainText
		c.Model.PrevIndex = index
	})

	// application vim-like navigation
	c.View.AppList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			if event.Rune() == 'g' {
				c.View.AppList.SetCurrentItem(0)
			}

			if event.Rune() == 'G' {
				c.View.AppList.SetCurrentItem(c.View.AppList.GetItemCount() - 1)
			}

			if event.Rune() == 'k' {
				if c.View.AppList.GetCurrentItem() == 0 {
					c.View.AppList.SetCurrentItem(-1)
					return event
				}

				c.View.AppList.SetCurrentItem(c.View.AppList.GetCurrentItem() - 1)
			}

			if event.Rune() == 'j' {
				if c.View.AppList.GetCurrentItem()+1 == c.View.AppList.GetItemCount() {
					c.View.AppList.SetCurrentItem(0)
					return event
				}

				c.View.AppList.SetCurrentItem(c.View.AppList.GetCurrentItem() + 1)
			}
		}

		return event
	})

	// TODO: improve tab / shift tab key logic
	c.View.SideBar.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			c.View.App.SetFocus(c.View.StatusBox)
			return nil
		}

		if event.Key() == tcell.KeyBacktab {
			c.View.App.SetFocus(c.View.AppList)
			return nil
		}

		return event
	})

	c.View.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			if event.Rune() == 'q' {
				c.View.App.Stop()
				return nil
			}
		}

		return event
	})
}

func (c *AppController) Start() error {
	c.SetupEventHandlers()
	c.Model.LoadApplications()
	c.View.UpdateAppList(c.Model.Applications)
	c.View.App.SetRoot(c.View.MainPage, true)
	return c.View.App.Run()
}
