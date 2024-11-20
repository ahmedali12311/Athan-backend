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
		m.Basic.AutoConfirm,
	).Name = "wallet-transactions:auto-confirm:public"

	f := d.E.Group("/wallets")
	r := d.Requires(
		wallet.ScopeCustomer,
		wallet.ScopeAdmin,
	)

	f.GET("/me", m.Basic.Wallet, r).Name = "wallets:index:admin,customer"

	f.GET("", m.Basic.Index, r).Name = "wallets:index:customer,admin"
	f.GET("/:id", m.Basic.Show, r).Name = "wallets:show:customer,admin"
	f.POST("/initiate", m.Basic.Initiate).Name = "wallets:initiate:customer,admin"
	f.POST("/:id/confirm", m.Basic.Confirm).Name = "wallets:confirm:customer,admin"

	r1 := d.Requires(wallet.ScopeAdmin)

	f.DELETE("/:id", m.Basic.Destroy, r1).Name = "wallets:destroy:admin"
	f.POST("", m.Basic.Store, r1).Name = "wallets:store:customer,admin"
}
