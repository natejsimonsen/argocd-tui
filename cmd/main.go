package main

import (
	"example.com/main/services/argocd"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

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
	app := tview.NewApplication()
	setTheme()

	mainPage := tview.NewFlex()
	mainContent := tview.NewBox().
		SetBorder(true).
		SetTitle(" Main Content ")

	result := argocd.ListApplications()
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
		textList.AddItem(item.Metadata["name"].(string), "", 0, nil)
	}

	sideBar := tview.NewFlex().SetDirection(tview.FlexRow)

	sideBar.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			return nil
		}

		if event.Key() == tcell.KeyBacktab {
			return nil
		}

		return event
	})

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
