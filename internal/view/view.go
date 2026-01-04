package view

import (
	"fmt"
	"strings"

	"example.com/main/internal/model"
	"example.com/main/services/argocd"
	"example.com/main/services/config"
	"example.com/main/services/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
)

type AppView struct {
	App                  *tview.Application
	Pages                *tview.Pages
	Config               *config.Config
	MainPage             *tview.Flex
	SideBar              *tview.Flex
	AppTable             *tview.Table
	HelpModal            tview.Primitive
	HelpPage             *tview.List
	CommandBar           *tview.Flex
	SearchInput          *tview.InputField
	MainContentContainer *tview.Flex
	MainPageContainer    *tview.Flex
	MainTable            *tview.Table
	StatusBox            *tview.Box
	Logger               *logrus.Logger
}

func NewAppView(app *tview.Application, config *config.Config, logger *logrus.Logger) *AppView {
	theme := tview.Theme{
		PrimitiveBackgroundColor: config.Background,
		BorderColor:              config.Border,
		TitleColor:               config.Text,
		PrimaryTextColor:         config.Text,
		SecondaryTextColor:       config.Progressing,
		InverseTextColor:         config.Background,
	}

	tview.Styles = theme

	tableStyle := tcell.StyleDefault.
		Bold(true)

	mainPage := tview.NewFlex()
	mainContentContainer := tview.NewFlex()
	sideBar := tview.NewFlex().
		SetDirection(tview.FlexRow)
	commandBar := tview.NewFlex()
	mainTable := tview.NewTable()
	appTable := tview.NewTable()
	mainPageContainer := tview.NewFlex().
		SetDirection(tview.FlexRow)

	mainContentContainer.
		SetBorder(true).
		SetTitle(" Main Content ")

	bsBox := tview.NewBox().
		SetBorder(true).
		SetTitle(" FooBarBaz ")

	bsBox.
		SetFocusFunc(func() {
			bsBox.SetBorderColor(config.Selected)
		}).
		SetBlurFunc(func() {
			bsBox.SetBorderColor(config.Border)
		})

	appTable.
		SetTitle(" Applications ").
		SetBorder(true)

	appTable.
		SetBlurFunc(func() {
			appTable.SetBorderColor(config.Border)
		}).
		SetFocusFunc(func() {
			appTable.SetBorderColor(config.Selected)
		})

	appTable.
		SetSelectedStyle(tableStyle).
		SetSelectable(true, false)

	mainTable.SetSelectedStyle(tableStyle)
	mainTable.SetSelectable(true, false)

	sideBar.
		AddItem(appTable, 0, 1, true).
		AddItem(bsBox, 0, 1, true)

	mainContentContainer.
		AddItem(mainTable, 0, 1, true).
		SetBorder(true)

	mainContentContainer.
		SetBlurFunc(func() {
			mainContentContainer.SetBorderColor(config.Border)
		}).
		SetFocusFunc(func() {
			mainContentContainer.SetBorderColor(config.Selected)
		})

	mainPage.
		AddItem(sideBar, 0, 1, true).
		AddItem(mainContentContainer, 0, 3, false)

	commandBar.
		SetBorder(true)

	mainPageContainer.
		AddItem(mainPage, 0, 1, true)

	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(p, height, 1, true).
				AddItem(nil, 0, 1, false), width, 1, true).
			AddItem(nil, 0, 1, false)
	}

	helpPage := tview.NewList().
		SetHighlightFullLine(true).
		ShowSecondaryText(false).
		SetSelectedBackgroundColor(config.Selected).
		SetSelectedTextColor(config.Background)

	helpPage.
		SetBorder(true).
		SetTitle("Help Page")

	helpPage.
		SetBlurFunc(func() {
			helpPage.SetBorderColor(config.Border)
		}).
		SetFocusFunc(func() {
			helpPage.SetBorderColor(config.Selected)
		})

	helpModal := modal(helpPage, 80, 40)

	pages := tview.NewPages().
		AddPage("main page", mainPageContainer, true, true).
		AddPage("help page", helpModal, true, false)

	appView := &AppView{
		App:                  app,
		Pages:                pages,
		MainPage:             mainPage,
		SideBar:              sideBar,
		HelpPage:             helpPage,
		HelpModal:            helpModal,
		AppTable:             appTable,
		MainContentContainer: mainContentContainer,
		MainPageContainer:    mainPageContainer,
		CommandBar:           commandBar,
		MainTable:            mainTable,
		StatusBox:            bsBox,
		Config:               config,
		Logger:               logger,
	}

	// appView.AddSearchInput()

	return appView
}

func (v *AppView) AddSearchInput() {
	searchInput := tview.NewInputField()
	searchInput.SetFieldBackgroundColor(v.Config.Background).
		SetBorder(false)

	v.SearchInput = searchInput

	v.CommandBar.
		AddItem(searchInput, 0, 1, true)
}

func (v *AppView) ToggleHelp(commands map[model.Context]map[model.KeyStroke]*model.Command) {
	if page, _ := v.Pages.GetFrontPage(); page == "help page" {
		v.RemoveHelp()
		return
	}

	v.Pages.ShowPage("help page")

	for ctx, cmdMap := range commands {
		for trigger, cmd := range cmdMap {
			v.HelpPage.AddItem(fmt.Sprintf("%c - %-10s - %-10s", trigger, ctx, cmd), "", 0, nil)
		}
	}
}

func (v *AppView) RemoveHelp() {
	v.Pages.HidePage("help page")
	v.HelpPage.Clear()
}

func (v *AppView) UpdateAppTable(apps []argocd.ApplicationItem) {
	v.AppTable.Clear()

	if len(apps) == 0 {
		v.AppTable.SetCell(0, 0,
			tview.NewTableCell("No data").
				SetTextColor(v.Config.Text).
				SetAlign(tview.AlignLeft))
		return
	}

	for i, app := range apps {
		color := v.Config.Progressing

		switch app.Status.Health.Status {
		case argocd.StatusDegraded:
			color = v.Config.Degraded
		case argocd.StatusHealthy:
			color = v.Config.Healthy
		case argocd.StatusProgressing:
			color = v.Config.Progressing
		case argocd.StatusMissing:
			color = v.Config.Missing
		}

		tableCell := tview.NewTableCell(app.Metadata.Name).
			SetTextColor(color).
			SetAlign(tview.AlignLeft).
			SetExpansion(1)

		tableCell.
			SetSelectedStyle(
				tcell.StyleDefault.
					Background(color).
					Foreground(utils.GetContrastColor(color)).
					Bold(true),
			)

		v.AppTable.SetCell(i, 0, tableCell)
	}

	v.AppTable.Select(0, 0)
}

func (v *AppView) RemoveSearchBar() {
	v.MainPageContainer.Clear()
	v.MainPageContainer.AddItem(v.MainPage, 0, 1, false)
	if v.SearchInput != nil {
		v.SearchInput.SetText("")
		v.CommandBar.RemoveItem(v.SearchInput)
	}
	v.MainContentContainer.SetTitle(" Main Content ")
	v.App.SetFocus(v.AppTable)
}

func (v *AppView) AddSearchBar() {
	v.MainPageContainer.Clear()
	v.AddSearchInput()
	v.MainPageContainer.AddItem(v.CommandBar, 3, 0, true)
	v.MainPageContainer.AddItem(v.MainPage, 0, 1, false)
	v.App.SetFocus(v.CommandBar)
}

func (v *AppView) Scroll(dir int) {
	prim := v.App.GetFocus()

	switch t := prim.(type) {
	case *tview.List:
		if dir == 1 {
			if t.GetCurrentItem() == 0 {
				t.SetCurrentItem(-1)
				return
			}

			t.SetCurrentItem(t.GetCurrentItem() - 1)
		}
		if dir == -1 {
			if t.GetCurrentItem()+1 == t.GetItemCount() {
				t.SetCurrentItem(0)
				return
			}

			t.SetCurrentItem(t.GetCurrentItem() + 1)
		}
	case *tview.Table:
		row, _ := t.GetSelection()
		offset := 1
		newRow := row + offset*-1*dir

		if newRow <= 0 && t == v.MainTable {
			newRow = 1
		}

		if newRow >= t.GetRowCount() {
			return
		}

		t.Select(max(newRow, 0), 0)
	}
}

func (v *AppView) ScrollTo(row int) {
	prim := v.App.GetFocus()

	switch t := prim.(type) {
	case *tview.List:
		t.SetCurrentItem(row)
	case *tview.Table:
		if row == 0 && t == v.MainTable {
			t.Select(1, 0)
			return
		}
		if row < 0 {
			rows := t.GetRowCount()
			t.Select(rows-1, 0)
			return
		}

		t.Select(row, 0)
	}
}

func (v *AppView) ToggleCommandBar() {
	if v.MainPageContainer.GetItemCount() > 1 {
		v.RemoveSearchBar()
		return
	}
	v.AddSearchBar()
}

func (v *AppView) ClearSearch() {
	v.RemoveSearchBar()
}

func (v *AppView) SetSearchTitle(search string) {
	title := strings.Split(v.MainContentContainer.GetTitle(), "/")[0]
	v.MainContentContainer.SetTitle(fmt.Sprintf("%s / %s ", title, search))
}

// TODO: refactor to global func
// func (v *AppView) PageMainContent(direction int) {
// 	rowNums := v.MainTable.GetRowCount()
// 	row, _ := v.MainTable.GetSelection()
// 	offset := rowNums / 2
// 	newRow := row + offset*direction
//
// 	if newRow <= 0 {
// 		newRow = 1
// 	}
//
// 	if newRow >= rowNums {
// 		v.MainTable.Select(rowNums-1, 0)
// 		return
// 	}
//
// 	v.MainTable.Select(newRow, 0)
// }

func (v *AppView) UpdateMainContent(resources []argocd.ApplicationNode) {
	v.MainTable.Clear()

	if len(resources) == 0 {
		v.MainTable.SetCell(0, 0,
			tview.NewTableCell("No data").
				SetTextColor(v.Config.Text).
				SetAlign(tview.AlignLeft))
		return
	}

	columns := []string{
		"Name",
		"Kind",
		"Health",
		"Namespace",
		"Version",
		"Resource Version",
		"Images",
	}

	for i, column := range columns {
		v.MainTable.SetCell(
			0,
			i,
			tview.NewTableCell(column).
				SetTextColor(v.Config.Header).
				SetAlign(tview.AlignLeft),
		).
			SetFixed(1, i)
	}

	for row, manifest := range resources {
		color := v.Config.Progressing

		switch manifest.Health.Status {
		case string(argocd.StatusDegraded):
			color = v.Config.Degraded
		case string(argocd.StatusHealthy):
			color = v.Config.Healthy
		case string(argocd.StatusProgressing):
			color = v.Config.Progressing
		case string(argocd.StatusMissing):
			color = v.Config.Missing
		}

		for i, column := range columns {
			value := ""

			switch column {
			case "Name":
				value = manifest.Name
			case "Kind":
				value = manifest.Kind
			case "Health":
				value = manifest.Health.Status
			case "Namespace":
				value = manifest.Namespace
			case "Version":
				value = manifest.Version
			case "Resource Version":
				value = manifest.ResourceVersion
			case "Images":
				value = strings.Join(manifest.Images, ", ")
			}

			tableCell := tview.NewTableCell(value).
				SetTextColor(color).
				SetAlign(tview.AlignLeft)

			tableCell.
				SetSelectedStyle(
					tcell.StyleDefault.
						Background(color).
						Foreground(utils.GetContrastColor(color)).
						Bold(true),
				)

			if i == len(columns)-1 {
				tableCell.SetExpansion(1)
			}

			v.MainTable.SetCell(row+1, i, tableCell)
		}
	}

	v.MainTable.Select(1, 0).ScrollToBeginning()
}

func (v *AppView) GetColorTag(status argocd.ApplicationHealthStatus) (string, tcell.Color) {
	colorTag := ""
	var color tcell.Color

	switch status {
	case argocd.StatusHealthy:
		color = v.Config.Healthy
	case argocd.StatusDegraded:
		color = v.Config.Degraded
	case argocd.StatusMissing:
		color = v.Config.Missing
	case argocd.StatusUnknown:
		color = v.Config.Header
	case argocd.StatusProgressing:
		color = v.Config.Progressing
	}

	colorTag = fmt.Sprintf("[%s]", color.CSS())

	return colorTag, color
}
