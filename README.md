# Minimal C Compiler

A minimal, functional C compiler written in **Go**. Compiles a subset of C source code to machine code (via assembly), as an educational project and demonstration of compiler fundamentals.

---

## High-Level Architecture

```
  C source        Compiler (Go)                    Output
  (.c file)
       │
       ▼
   ┌────────┐
   │ Lexer  │
   └───┬────┘
       │ tokens
       ▼
   ┌────────┐
   │ Parser │
   └───┬────┘
       │ AST
       ▼
   ┌─────────────────┐
   │ Semantic        │
   │ Analysis        │
   └───┬─────────────┘
       │ checked AST
       ▼
   ┌─────────────────┐
   │ Code Generator  │
   └───┬─────────────┘
       │
       ▼
   Assembly (.s)  ──►  Object (.o)  ──►  Executable
   (as)                  (ld)
```

**Pipeline:** C source → **Lexer** (tokens) → **Parser** (AST) → **Semantic analysis** (checked AST + symbol table) → **Code generator** (assembly) → **Assembler** (object file) → **Linker** (executable).

---

## Development Strategy

Build and test **one phase at a time**; each phase consumes the previous phase’s output. Recommended order:

### Phase 1: Lexer (Tokenizer)

**Goal:** Turn raw source text into a stream of tokens.

- **Define tokens:** keywords (`int`, `return`, `if`, `else`, `while`, etc.), identifiers, literals (integers), operators (`+`, `-`, `*`, `/`, `==`, `!=`, `<`, `>`, `=`), punctuation (`;`, `,`, `(`, `)`, `{`, `}`).
- **Implement:** Single pass over the source; skip whitespace and comments; emit token type + value (e.g. line/column for errors).
- **Test:** Unit tests: given a string, assert the token slice (e.g. `"int x = 42;"` → `INT`, `IDENT("x")`, `ASSIGN`, `INT_LIT(42)`, `SEMICOLON`).
- **Deliverable:** A `Lexer` that returns `[]Token` (or a stream) and handles basic error reporting (e.g. invalid character).

---

### Phase 2: Parser

**Goal:** Turn token stream into an Abstract Syntax Tree (AST).

- **Define AST nodes:** Program, function declarations, statements (variable decl, assignment, `return`, `if`/`else`, `while`), expressions (binary ops, unary, identifiers, literals, function calls). Use Go structs and interfaces (e.g. `Expr`, `Stmt`).
- **Implement:** Recursive-descent parser driven by the grammar (e.g. expression precedence: additive → multiplicative → unary → primary). One function per non-terminal.
- **Test:** Unit tests: given tokens (or a small C snippet lexed first), assert the AST shape. Start with trivial programs (e.g. `int main() { return 0; }`).
- **Deliverable:** A `Parser` that produces an AST (e.g. `*ast.Program`) and reports syntax errors with location.

---

### Phase 3: Semantic Analysis

**Goal:** Validate program meaning and build symbol table.

- **Symbol table:** Map identifier names to types and scope (e.g. map + scope stack, or explicit scope structs). Support function and variable declarations.
- **Implement:** Single (or multiple) pass over AST: resolve names, check types (assignment compatibility, function call arity and types), ensure `return` type matches, no duplicate declarations in same scope.
- **Test:** Unit tests: valid programs pass; invalid ones (undefined variable, type mismatch, wrong return type) get clear errors.
- **Deliverable:** Checked AST (possibly annotated with types/symbols) and a symbol table; reject invalid programs with readable messages.

---

### Phase 4: Code Generator

**Goal:** Emit assembly for a target architecture (e.g. x86-64 or RISC-V).

- **Choose target:** e.g. x86-64 Linux (GNU assembler syntax) or RISC-V for simplicity. Document calling convention (e.g. SysV ABI).
- **Implement:** Walk the checked AST; for each construct emit the corresponding assembly (prologue/epilogue, load/store, arithmetic, branches, function calls). Use a simple register allocation strategy (e.g. fixed registers or stack-only for a minimal version).
- **Test:** Emit assembly for small programs; run `as` and `ld` (or `gcc`) and execute the binary; assert exit code or output.
- **Deliverable:** Code generator that produces a single `.s` file (or in-memory string) that assembles and links to a runnable executable.

---

### Phase 5: Driver & Toolchain Integration

**Goal:** End-to-end “compile this file → run executable.”

- **Implement:** Main entrypoint: read `.c` file → lex → parse → semantic check → codegen → write `.s` → invoke assembler → invoke linker (or `gcc -o out out.s`). Pass through exit codes and errors.
- **Test:** Integration tests: compile and run a few minimal C programs (e.g. `return 42`, simple arithmetic, one `if` and one `while`).
- **Deliverable:** CLI that compiles a C file to an executable (e.g. `go run ./cmd/compiler main.c -o main`).

---

## Suggested C Subset (MVP)

- Types: `int` only.
- Statements: variable declaration, assignment, `return`, `if`/`else`, `while`.
- Expressions: literals, identifiers, binary operators (`+`, `-`, `*`, `/`, `==`, `!=`, `<`, `>`), function calls (e.g. `main` only or one extra function).
- One function: `main()` returning `int`. Expand to multiple functions once the above works.

---

## Tech Stack

| Component        | Choice                          |
|------------------|---------------------------------|
| Language         | Go                              |
| Target code      | Assembly (e.g. x86-64 or RISC-V) |
| Assembler/Linker | System `as` / `ld` or `gcc`     |

---

## Repository Structure (Suggested)

```
.
├── cmd/
│   └── compiler/     # CLI entrypoint
├── internal/
│   ├── lexer/        # Tokenizer
│   ├── parser/       # Parser & AST
│   ├── semantic/     # Symbol table & type checking
│   └── codegen/      # AST → Assembly
├── go.mod
└── README.md
```

You can adjust package names and add `pkg/` or `ast/` as needed once you start coding.
