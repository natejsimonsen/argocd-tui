package view

import (
	"example.com/main/services/argocd"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type AppView struct {
	App         *tview.Application
	MainPage    *tview.Flex
	SideBar     *tview.Flex
	AppList     *tview.List
	MainContent *tview.TextView
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
		SetDynamicColors(true)

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

	sideBar.
		AddItem(textList, 0, 1, true).
		AddItem(bsBox, 0, 1, true)

	mainPage.
		AddItem(sideBar, 0, 1, true).
		AddItem(mainContent, 0, 3, false)

	return &AppView{
		App:         app,
		MainPage:    mainPage,
		SideBar:     sideBar,
		AppList:     textList,
		MainContent: mainContent,
		StatusBox:   bsBox,
	}
}

func (v *AppView) UpdateAppList(apps []argocd.ApplicationItem) {
	for _, app := range apps {
		colorTag := GetColorTag(app.Status.Health.Status)

		name := fmt.Sprintf("%s%s", colorTag, app.Metadata.Name)
		v.AppList.AddItem(name, "", 0, nil)
	}
}

func (v *AppView) UpdateMainContent(resources []argocd.ApplicationNode) {
	text := ""

	for _, manifest := range resources {
		text += fmt.Sprintf("Kind: [aqua]%s[end]\nName: [green]%s[end]\n\n", manifest.Kind, manifest.Name)
	}

	v.MainContent.SetText(text)
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
