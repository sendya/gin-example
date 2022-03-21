package controller

import "go.uber.org/fx"

var Modules = fx.Invoke(
	// demo controller
	NewHelloController,

	NewUserController,
)
