# ormc

Go code generator for `tinywasm/model` definitions: reads hand-written
`model.Definition` literals via AST and emits the concrete struct, schema,
codec, validation, and typed query helpers the runtime consumes.

> Quick-start and CLI usage will land with the phase-B implementation.

## Documentation

- [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) — what ormc is post-split,
  input contract (`Type:` constructor expressions), AST-only mechanism,
  storage resolution order, generated-code contract.
- [docs/DESIGN.md](docs/DESIGN.md) — decision record: why storage resolves
  via the cached dependency probe (directive and marker alternatives
  rejected).
- [docs/SYNC_DESIGN.md](docs/SYNC_DESIGN.md) — rationale and trade-offs of
  tool-driven dev schema sync.
- [docs/WHY_GENERATED_CODE_IS_FREE.md](docs/WHY_GENERATED_CODE_IS_FREE.md) —
  why the generated volume costs zero bytes in WASM binaries.
- [docs/WHY_PACKAGE_LEVEL_SCHEMA.md](docs/WHY_PACKAGE_LEVEL_SCHEMA.md) —
  why generated `Schema()` returns a package-level variable (zero alloc).
- [docs/diagrams/ORMC_FLOW.md](docs/diagrams/ORMC_FLOW.md) — Definition →
  generated code flow, error table, runtime validation.
- [docs/diagrams/DB_SYNC.md](docs/diagrams/DB_SYNC.md) — dev database
  schema-sync flow (tool-driven).
