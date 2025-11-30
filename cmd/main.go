package main

import (
	"fmt"

	"example.com/main/services/argocd"
	"example.com/main/services/logger"
	"example.com/main/services/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

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

func setTheme() {
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
}

func main() {
	PrevText := ""
	PrevIndex := -1

	app := tview.NewApplication()
	setTheme()

	logger := logger.SetupLogger()
	argocdSvc := argocd.NewService(logger)

	mainPage := tview.NewFlex()
	mainContent := tview.NewTextView().
		SetDynamicColors(true)

	result := argocdSvc.ListApplications()
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

	textList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		newText := utils.StripTags(mainText)
		textList.SetItemText(index, newText, "")

		if PrevText != "" {
			textList.SetItemText(PrevIndex, PrevText, "")
		}

		PrevText = mainText
		PrevIndex = index

		manifests := argocdSvc.GetResourceTree(utils.StripTags(mainText))
		text := ""

		for _, manifest := range manifests {
			text += fmt.Sprintf("Kind: [aqua]%s[end]\nName: [green]%s[end]\n\n", manifest.Kind, manifest.Name)
		}

		mainContent.SetText(text)
	})

	mainContent.
		SetBorder(true).
		SetTitle(" Main Content ")

	textList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			if event.Rune() == 'g' {
				textList.SetCurrentItem(0)
			}

			if event.Rune() == 'G' {
				textList.SetCurrentItem(textList.GetItemCount() - 1)
			}

			if event.Rune() == 'k' {
				if textList.GetCurrentItem() == 0 {
					textList.SetCurrentItem(-1)
					return event
				}

				textList.SetCurrentItem(textList.GetCurrentItem() - 1)
			}

			if event.Rune() == 'j' {
				if textList.GetCurrentItem()+1 == textList.GetItemCount() {
					textList.SetCurrentItem(0)
					return event
				}

				textList.SetCurrentItem(textList.GetCurrentItem() + 1)
			}
		}

		return event
	})

	for _, item := range result.Items {
		colorTag := GetColorTag(item.Status.Health.Status)

		name := fmt.Sprintf("%s%s", colorTag, item.Metadata.Name)
		textList.AddItem(name, "", 0, nil)
	}

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

	sideBar.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			app.SetFocus(bsBox)
			return nil
		}

		if event.Key() == tcell.KeyBacktab {
			app.SetFocus(textList)
			return nil
		}

		return event
	})

	sideBar.
		AddItem(textList, 0, 1, true).
		AddItem(bsBox, 0, 1, true)

	mainPage.
		AddItem(sideBar, 0, 1, true).
		AddItem(mainContent, 0, 3, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			if event.Rune() == 'q' {
				app.Stop()
				return nil
			}
		}

		return event
	})

	if err := app.SetRoot(mainPage, true).Run(); err != nil {
		panic(err)
	}
}
