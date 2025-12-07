package view

import (
	"fmt"
	"strings"

	"example.com/main/services/argocd"
	"example.com/main/services/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type AppView struct {
	App                  *tview.Application
	Config               config.Config
	MainPage             *tview.Flex
	SideBar              *tview.Flex
	AppList              *tview.List
	CommandBar           *tview.Flex
	SearchInput          *tview.InputField
	MainContentContainer *tview.Flex
	MainPageContainer    *tview.Flex
	MainTable            *tview.Table
	StatusBox            *tview.Box
}

func NewAppView(app *tview.Application, config *config.Config) *AppView {
	theme := tview.Theme{
		PrimitiveBackgroundColor:    tcell.ColorDefault,
		ContrastBackgroundColor:     tcell.ColorBlack,
		MoreContrastBackgroundColor: tcell.ColorBlack,
		BorderColor:                 tcell.ColorGray,
		TitleColor:                  tcell.ColorGreen,
		GraphicsColor:               tcell.ColorYellow,
		PrimaryTextColor:            tcell.ColorWhite,
		SecondaryTextColor:          tcell.ColorAqua,
		TertiaryTextColor:           tcell.ColorFuchsia,
		InverseTextColor:            tcell.ColorBlack,
		ContrastSecondaryTextColor:  tcell.ColorBlack,
	}

	tview.Styles = theme

	mainPage := tview.NewFlex()
	mainContentContainer := tview.NewFlex()
	sideBar := tview.NewFlex().
		SetDirection(tview.FlexRow)
	commandBar := tview.NewFlex()
	mainTable := tview.NewTable()
	searchInput := tview.NewInputField()
	mainPageContainer := tview.NewFlex().
		SetDirection(tview.FlexRow)

	textList := tview.NewList().
		SetHighlightFullLine(true).
		ShowSecondaryText(false).
		SetSelectedBackgroundColor(tcell.ColorAqua).
		SetSelectedTextColor(tcell.ColorBlack)

	textList.
		SetBlurFunc(func() {
			textList.SetBorderColor(tview.Styles.BorderColor)
		}).
		SetFocusFunc(func() {
			textList.SetBorderColor(tcell.ColorAquaMarine)
		})

	textList.
		SetBorder(true).
		SetTitle(" Applications ")

	mainContentContainer.
		SetBorder(true).
		SetTitle(" Main Content ")

	bsBox := tview.NewBox().
		SetBorder(true).
		SetTitle(" FooBarBaz ")

	bsBox.
		SetFocusFunc(func() {
			bsBox.SetBorderColor(tcell.ColorAquaMarine)
		}).
		SetBlurFunc(func() {
			bsBox.SetBorderColor(tview.Styles.BorderColor)
		})

	tableStyle := tcell.StyleDefault.
		Background(config.Background).
		Foreground(tcell.ColorDarkSlateGray).
		Bold(true)

	mainTable.SetSelectedStyle(tableStyle)
	mainTable.SetSelectable(true, false)

	sideBar.
		AddItem(textList, 0, 1, true).
		AddItem(bsBox, 0, 1, true)

	commandBar.
		AddItem(searchInput, 0, 1, true)

	mainContentContainer.
		AddItem(mainTable, 0, 1, true).
		SetBorder(true)

	mainPage.
		AddItem(sideBar, 0, 1, true).
		AddItem(mainContentContainer, 0, 3, false)

	// command bar stuff
	commandBar.
		SetBorder(true)

	mainPageContainer.
		AddItem(mainPage, 0, 1, true)

	return &AppView{
		App:                  app,
		MainPage:             mainPage,
		SideBar:              sideBar,
		AppList:              textList,
		MainContentContainer: mainContentContainer,
		MainPageContainer:    mainPageContainer,
		SearchInput:          searchInput,
		CommandBar:           commandBar,
		MainTable:            mainTable,
		StatusBox:            bsBox,
	}
}

func (v *AppView) UpdateTitles(index, prevIndex int, text, prevText string) {
	v.AppList.SetItemText(prevIndex, prevText, "")
	v.AppList.SetItemText(index, "[::b]"+text, "")
}

func (v *AppView) UpdateAppList(apps []argocd.ApplicationItem) {
	for _, app := range apps {
		colorTag := GetColorTag(app.Status.Health.Status)

		name := fmt.Sprintf("%s%s", colorTag, app.Metadata.Name)
		v.AppList.AddItem(name, "", 0, nil)
	}
}

func (v *AppView) ToggleCommandBar() {
	if v.MainPageContainer.GetItemCount() > 1 {
		v.MainPageContainer.Clear()
		v.MainPageContainer.AddItem(v.MainPage, 0, 1, false)
		return
	}

	v.MainPageContainer.Clear()
	v.MainPageContainer.AddItem(v.CommandBar, 3, 0, true)
	v.MainPageContainer.AddItem(v.MainPage, 0, 1, false)
}

func (v *AppView) HorizontallyScrollMainTable(direction int) {
	row, col := v.MainTable.GetSelection()
	offset := 1
	newCol := col + offset*direction

	if newCol == 0 {
		newCol++
	}

	if newCol == v.MainTable.GetRowCount() {
		return
	}

	v.MainTable.Select(row, newCol)
}

func (v *AppView) ScrollMainContent(direction int) {
	row, _ := v.MainTable.GetSelection()
	offset := 1
	newRow := row + offset*direction

	if newRow == 0 {
		newRow++
	}

	if newRow == v.MainTable.GetRowCount() {
		return
	}

	v.MainTable.Select(newRow, 0)
}

func (v *AppView) PageMainContent(direction int) {
	row, _ := v.MainTable.GetSelection()
	offset := 1
	newRow := row + offset*direction

	if newRow == 0 {
		newRow++
	}

	if newRow == v.MainTable.GetRowCount() {
		return
	}

	v.MainTable.Select(newRow, 0)
}

func (v *AppView) UpdateMainContent(resources []argocd.ApplicationNode) {
	v.MainTable.Clear()

	if len(resources) == 0 {
		v.MainTable.SetCell(0, 0,
			tview.NewTableCell("No data").
				SetTextColor(tcell.ColorWhite).
				SetAlign(tview.AlignLeft))
		return
	}

	columns := []string{"Name", "Kind", "Health", "Namespace", "Version", "Resource Version", "Images"}

	for i, column := range columns {
		v.MainTable.SetCell(0, i,
			tview.NewTableCell(column).
				SetTextColor(tcell.ColorWhite).
				SetAlign(tview.AlignLeft))
	}

	for row, manifest := range resources {
		color := tcell.ColorWhiteSmoke

		if manifest.Health.Status == string(argocd.StatusDegraded) {
			color = tcell.ColorRed
		}

		if manifest.Health.Status == string(argocd.StatusHealthy) {
			color = tcell.ColorLightGreen
		}

		if manifest.Health.Status == string(argocd.StatusProgressing) {
			color = tcell.ColorLightBlue
		}

		if manifest.Health.Status == string(argocd.StatusMissing) {
			color = tcell.ColorLightYellow
		}

		for i, column := range columns {
			value := ""

			if column == "Name" {
				value = manifest.Name
			}
			if column == "Kind" {
				value = manifest.Kind
			}
			if column == "Health" {
				value = manifest.Health.Status
			}
			if column == "Namespace" {
				value = manifest.Namespace
			}
			if column == "Version" {
				value = manifest.Version
			}
			if column == "Resource Version" {
				value = manifest.ResourceVersion
			}
			if column == "Images" {
				value = strings.Join(manifest.Images, ", ")
			}

			tableCell := tview.NewTableCell(value).
				SetTextColor(color).
				SetAlign(tview.AlignLeft)

			if manifest.Health.Status != "" {
				tableCell.
					SetSelectedStyle(
						tcell.StyleDefault.
							Background(color).
							Foreground(tcell.ColorDarkSlateGray).
							Bold(true),
					)
			}

			v.MainTable.SetCell(row+1, i, tableCell)
		}
	}

	v.MainTable.Select(1, 0)
}

func GetColorTag(status argocd.ApplicationHealthStatus) string {
	colorTag := ""
	if string(status) == string(argocd.StatusHealthy) {
		colorTag = "[green]"
	}

	if string(status) == string(argocd.StatusDegraded) {
		colorTag = "[red]"
	}

	if string(status) == string(argocd.StatusMissing) {
		colorTag = "[yellow]"
	}

	if string(status) == string(argocd.StatusUnknown) {
		colorTag = "[grey]"
	}

	if string(status) == string(argocd.StatusProgressing) {
		colorTag = "[blue]"
	}

	return colorTag
}
