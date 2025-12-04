package view

import (
	"fmt"

	"example.com/main/services/argocd"
	"example.com/main/services/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type AppView struct {
	App         *tview.Application
	Config      config.Config
	MainPage    *tview.Flex
	SideBar     *tview.Flex
	AppList     *tview.List
	MainContent *tview.TextView
	MainTable   *tview.Table
	StatusBox   *tview.Box
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
	mainContent := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true)

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

	mainContent.
		SetBorder(true).
		SetTitle(" Main Content ")

	sideBar := tview.NewFlex().SetDirection(tview.FlexRow)

	textList.
		SetBorder(true).
		SetTitle(" Applications ")

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

	mainTable := tview.NewTable()

	tableStyle := tcell.StyleDefault.
		Background(config.Background).
		Foreground(tcell.ColorDarkSlateGray).
		Bold(true)

	mainTable.SetSelectedStyle(tableStyle)

	mainTable.SetTitle(" Main Content ")
	mainTable.SetSelectable(true, false)

	sideBar.
		AddItem(textList, 0, 1, true).
		AddItem(bsBox, 0, 1, true)

	mainPage.
		AddItem(sideBar, 0, 1, true).
		AddItem(mainTable, 0, 3, false)

	return &AppView{
		App:         app,
		MainPage:    mainPage,
		SideBar:     sideBar,
		AppList:     textList,
		MainContent: mainContent,
		MainTable:   mainTable,
		StatusBox:   bsBox,
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
	row, _ := v.MainContent.GetScrollOffset()
	_, _, _, height := v.MainContent.Primitive.GetRect()
	offset := height / 2
	newScroll := row + offset*direction
	v.MainContent.ScrollTo(newScroll, 0)
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

	type TableCell struct {
		ColumnName  string
		ColumnValue string
	}

	columns := []string{"Name", "Kind", "Status"}

	for i, column := range columns {
		v.MainTable.SetCell(0, i,
			tview.NewTableCell(column).
				SetTextColor(tcell.ColorWhite).
				SetAlign(tview.AlignLeft))
	}

	for row, manifest := range resources {
		color := tcell.ColorLightGray

		if manifest.Health.Status == string(argocd.StatusDegraded) {
			color = tcell.ColorRed
		}

		for i, column := range columns {
			value := ""

			if column == "Name" {
				value = manifest.Name
			}
			if column == "Kind" {
				value = manifest.Kind
			}
			if column == "Status" {
				value = manifest.Health.Status
			}

			tableCell := tview.NewTableCell(value).
				SetTextColor(color).
				SetAlign(tview.AlignLeft)

			if manifest.Health.Status == string(argocd.StatusDegraded) {
				tableCell.
					SetSelectedStyle(
						tcell.StyleDefault.
							Background(tcell.ColorRed).
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
