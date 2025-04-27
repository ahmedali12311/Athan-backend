package wallet_controller

import (
	"app/controller"
	"app/models/wallet"
)

func (m *Controllers) SetBasicRoutes(
	d *controller.RouterDependencies,
) {
	// TLync auto-confirm route triggered by tlync backend-url callback
	d.E.POST(
		"/wallet-transactions/auto-confirm",
		m.Basic.TylncAutoConfirm,
	).Name = "wallet-transactions:auto-confirm:public"

	f := d.E.Group("/wallets")
	r := d.Requires(
		wallet.ScopeCustomer,
		wallet.ScopeAdmin,
	)

	f.GET("/me", m.Basic.Wallet, r).Name = "wallets:index:admin,customer"

	f.GET("", m.Basic.Index, r).Name = "wallets:index:customer,admin"
	f.GET("/:id", m.Basic.Show, r).Name = "wallets:show:customer,admin"

	r1 := d.Requires(wallet.ScopeAdmin)

	f.DELETE("/:id", m.Basic.Destroy, r1).Name = "wallets:destroy:admin"
	f.POST("", m.Basic.Store, r1).Name = "wallets:store:customer,admin"

	tylnc := d.E.Group("/tlync")

	tylnc.POST("/initiate", m.Basic.TylncInitiate, r).Name = "wallets:tylnc-initiate:customer,admin"
	tylnc.POST("/:id/confirm", m.Basic.TylncConfirm, r).Name = "wallets:tylnc-confirm:customer,admin"

	masarat := d.E.Group("/masarat")

	masarat.POST("/initiate", m.Basic.MasaratInitiate, r).Name = "wallets:masarat-initiate:admin,customer"
	masarat.POST("/confirm", m.Basic.MasaratConfirm, r).Name = "wallets:masarat-confirm:admin,customer"

	edfali := d.E.Group("/edfali")
	edfali.POST("/initiate", m.Basic.EdfaliInitiate, r).Name = "wallets:edfali-initiate:admin,customer"
	edfali.POST("/confirm", m.Basic.EdfaliConfirm, r).Name = "wallets:edfali-confirm:admin,customer"

	sadad := f.Group("/sadad")

	sadad.POST("/initiate", m.Basic.SadadInitiate, r).Name = "wallets:sadad-initiate:customer,admin"
	sadad.POST("/confirm", m.Basic.SadadConfirm, r).Name = "wallets:sadad-confirm:customer,admin"
	sadad.POST("/resend", m.Basic.SadadResendOTP, r).Name = "wallets:sadad-resend:customer,admin"
}
