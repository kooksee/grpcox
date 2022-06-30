package pf

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// PasswordModal is a modal which allows to user to enter a password
type PasswordModal struct {
	app.Compo

	Title         string // Title of the modal
	WrongPassword bool   // Whether the previously entered password was wrong

	ClearWrongPassword func()                // Handler to call when clearing the password
	OnSubmit           func(password string) // Handler to call to submit the password
	OnCancel           func()                // Handler to call when closing/cancelling the modal

	password string
}

func (c *PasswordModal) Render() app.UI {
	return &Modal{
		ID:           "password-modal",
		Title:        c.Title,
		DisableFocus: true,
		Body: []app.UI{
			app.Form().
				Class("pf-c-form").
				ID("password-modal-form").
				OnSubmit(func(ctx app.Context, e app.Event) {
					e.PreventDefault()

					// Submit the form
					c.OnSubmit(
						c.password,
					)

					c.clear()
				}).
				Body(
					app.Div().
						Class("pf-c-form__group").
						Body(
							app.Div().
								Class("pf-c-form__group-control").
								Body(
									&Autofocused{
										Component: app.Input().
											Class("pf-c-form-control").
											Required(true).
											Type("password").
											Placeholder("Password").
											Aria("label", "Password").
											Aria("invalid", c.WrongPassword).
											Aria("describedby", func() string {
												if c.WrongPassword {
													return "password-helper"
												}

												return ""
											}).
											OnInput(func(ctx app.Context, e app.Event) {
												c.password = ctx.JSSrc().Get("value").String()

												if c.ClearWrongPassword != nil {
													c.ClearWrongPassword()
												}
											}).
											Value(c.password),
									},
									app.If(
										c.WrongPassword,
										app.P().
											Class("pf-c-form__helper-text pf-m-error").
											ID("password-helper").
											Aria("live", "polite").
											Body(
												app.Text("The password is incorrect. Please try again."),
											),
									),
								),
						),
				),
		},
		Footer: []app.UI{
			app.Button().
				Class("pf-c-button pf-m-primary").
				Type("submit").
				Form("password-modal-form").
				Text("Continue"),
			app.Button().
				Class("pf-c-button pf-m-link").
				Type("button").
				Text("Cancel").
				OnClick(func(ctx app.Context, e app.Event) {
					c.clear()
					c.OnCancel()
				}),
		},
		OnClose: func() {
			c.clear()
			c.OnCancel()
		},
	}
}

func (c *PasswordModal) clear() {
	c.password = ""
}
