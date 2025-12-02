package view

import (
	"fmt"

	"example.com/main/services/argocd"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type AppView struct {
	App         *tview.Application
	MainPage    *tview.Flex
	SideBar     *tview.Flex
	AppList     *tview.List
	MainContent *tview.TextView
	MainTable   *tview.Table
	StatusBox   *tview.Box
}

func NewAppView(app *tview.Application) *AppView {
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

	mainTable := tview.NewTable().SetBorders(true)

	mainTable.SetTitle(" Main Content ")
	mainTable.SetSelectable(true, true)

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
	v.MainTable.Select(row+offset*direction, 0)
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
				SetAlign(tview.AlignCenter))
		return
	}

	v.MainTable.SetCell(0, 0,
		tview.NewTableCell("Name").
			SetTextColor(tcell.ColorWhite).
			SetAlign(tview.AlignCenter))
	v.MainTable.SetCell(0, 1,
		tview.NewTableCell("Kind").
			SetTextColor(tcell.ColorWhite).
			SetAlign(tview.AlignCenter))

	for row, manifest := range resources {
		v.MainTable.SetCell(row+1, 0,
			tview.NewTableCell(manifest.Name).
				SetTextColor(tcell.ColorWhite).
				SetAlign(tview.AlignCenter))
		v.MainTable.SetCell(row+1, 1,
			tview.NewTableCell(manifest.Kind).
				SetTextColor(tcell.ColorWhite).
				SetAlign(tview.AlignCenter))
	}
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
