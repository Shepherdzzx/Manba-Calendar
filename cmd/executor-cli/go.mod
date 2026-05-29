module github.com/Shepherdzzx/manba-alert/cmd/executor-cli

go 1.22

require (
	github.com/Shepherdzzx/manba-alert/executor v0.0.0
	github.com/Shepherdzzx/manba-alert/parser v0.0.0
)

replace github.com/Shepherdzzx/manba-alert/executor => ../../executor
replace github.com/Shepherdzzx/manba-alert/parser => ../../parser
