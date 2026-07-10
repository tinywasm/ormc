// Separate test module: isolates codegen fixture deps (orm runtime for
// generated models) so the root github.com/tinywasm/ormc module stays
// fmt + model + modfind + ddlc only.
module github.com/tinywasm/ormc/tests

go 1.25.2

require (
	github.com/tinywasm/model v0.0.8
	github.com/tinywasm/orm v0.9.27
	github.com/tinywasm/ormc v0.0.1
)

require (
	github.com/tinywasm/ddlc v0.0.4 // indirect
	github.com/tinywasm/fmt v0.25.1 // indirect
	github.com/tinywasm/modfind v0.0.4 // indirect
)

replace github.com/tinywasm/ormc => ..
