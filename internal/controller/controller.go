package controller

import (
	"strings"

	"example.com/main/internal/model"
	"example.com/main/internal/view"
	"example.com/main/services/argocd"
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

	c.View.SideBar.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		isShiftPressed := event.Modifiers()&tcell.ModShift != 0

		if event.Key() == tcell.KeyRune {
			if event.Rune() == 'J' {
				c.View.ScrollMainContent(1)
				return nil
			}

			if event.Rune() == 'l' {
				c.View.HorizontallyScrollMainTable(1)
				return nil
			}

			if event.Rune() == 'h' {
				c.View.HorizontallyScrollMainTable(-1)
				return nil
			}

			if event.Rune() == 'K' {
				c.View.ScrollMainContent(-1)
				return nil
			}

			if event.Rune() == 'D' {
				c.View.PageMainContent(1)
				return nil
			}

			if event.Rune() == 'U' {
				c.View.PageMainContent(-1)
				return nil
			}

			if event.Rune() == '/' {
				c.View.ToggleCommandBar()
				return nil
			}
		}

		if event.Key() == tcell.KeyPgDn {
			c.View.PageMainContent(1)
			return nil
		}

		if event.Key() == tcell.KeyPgUp {
			c.View.PageMainContent(-1)
			return nil
		}

		if event.Key() == tcell.KeyDown && isShiftPressed {
			c.View.ScrollMainContent(1)
			return nil
		}

		if event.Key() == tcell.KeyUp && isShiftPressed {
			c.View.ScrollMainContent(-1)
			return nil
		}

		// TODO: improve tab / shift tab key logic
		if event.Key() == tcell.KeyTab {
			c.View.App.SetFocus(c.View.StatusBox)
			return nil
		}

		if event.Key() == tcell.KeyBacktab {
			c.View.App.SetFocus(c.View.AppList)
			return nil
		}

		if event.Key() == tcell.KeyEsc {
			if c.Model.SearchString != "" {
				c.Model.LoadResources(c.Model.SelectedAppName)
				c.View.UpdateMainContent(c.Model.SelectedAppResources)
			}

			c.Model.SearchString = ""
			c.View.ClearSearch()
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

	c.View.CommandBar.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			if event.Rune() == '/' {
				c.View.ToggleCommandBar()
				return nil
			}
		}

		if event.Key() == tcell.KeyEsc {
			c.View.ToggleCommandBar()
			return nil
		}

		if event.Key() == tcell.KeyEnter {
			c.Model.SearchString = c.View.SearchInput.GetText()
			c.View.ToggleCommandBar()
			c.Model.SelectedAppResources = c.FilterContent(c.Model.SearchString)
			c.View.UpdateMainContent(c.Model.SelectedAppResources)
			return nil
		}

		return event
	})
}

func (c *AppController) FilterContent(search string) []argocd.ApplicationNode {
	var filteredResources []argocd.ApplicationNode

	for _, app := range c.Model.SelectedAppResources {
		if strings.Contains(app.Name, c.Model.SearchString) {
			filteredResources = append(filteredResources, app)
		}
	}

	return filteredResources
}

func (c *AppController) Start() error {
	c.SetupEventHandlers()
	c.Model.LoadApplications()
	c.View.UpdateAppList(c.Model.Applications)
	c.View.App.SetRoot(c.View.MainPageContainer, true)
	return c.View.App.Run()
}
