# Yarnball Language Specification (v2)
## 1. Introduction

**Yarnball** is an esoteric, stack-based programming language whose instructions read like a crochet pattern. This version prioritizes a simple, minimal core while keeping the surface syntax crochet-like.

**Breaking change:** This spec replaces the earlier `subpattern`/`rep` syntax. See examples in `/examples`.

---

## 2. Program styling
- Comments use `#` and can appear anywhere.
- Whitespace and commas are ignored.
- Optional headers are allowed; parsing starts after `STITCH GUIDE:` or `INSTRUCTIONS:` if present.
- `Row N:` and `Round N:` prefixes are ignored.
- Common filler words are ignored (e.g., `in`, `next`, `st`, `to`, `from`, `and`, `then`, `around`, `times`).

Example:
```yarnball
COZY HAT PATTERN
Author: You

STITCH GUIDE:
stitch brim = (
  sc 3
)

INSTRUCTIONS:
Row 1: ch 5
Row 2: brim
```

---

## 3. Stitch definitions

### Definition
```
stitch <name> = (
  ...instructions...
)
```

The stitch name may only be alphabetic characters.

### Call
Write the stitch name directly, or use the optional `use <name>` form.

---

## 4. Counts and repeats

### Repeating a single stitch
- Prefix: `3 sc`
- Postfix: `sc 3`

These expand to three consecutive `sc` instructions.

### Repeating a block
Use crochet-style blocks with `* ... *` (or `[ ... ]`) and `repeat`:
```
* sc inc * repeat 3
* sc inc * repeat while
* sc inc * repeat until
```

`repeat while` runs the block while the top of the stack is non-zero.  
`repeat until` runs the block while the top of the stack is zero.  
The condition is **peeked**, not popped.

---

## 5. Instruction set

### Stack + arithmetic
- **ch `<n>`**: push number
- **sc**: pop (discard)
- **sl st**: duplicate top
- **swap**: swap top two
- **over**: copy second item to top
- **pick `<n>`**: copy item at depth `<n>` to top (0 = top)
- **roll `<n>`**: rotate item at depth `<n>` to top
- **inc / dec**: increment / decrement top
- **bob**: add top two
- **hdc**: subtract top from second
- **dc**: multiply top two
- **tr**: divide second by top
- **cl**: modulo (second % top)
- **turn**: rotate top three (n1 n2 n3 -> n2 n3 n1)

### Comparisons
- **>**, **<**, **eq**, **neq**: compare top two and push 1 (true) or 0 (false)

### Output
- **pic**: pop and print ASCII character
- **yo**: pop and print number
- **fo**: halt immediately

---

## 6. Control flow

### if / else / end
```
if
  ...then...
else
  ...else...
end
```
Condition is popped; non-zero is true.

### repeat blocks
See Section 4. Use `repeat 3`, `repeat while`, or `repeat until`.

---

## 7. Step limit
The evaluator enforces a default step limit of 1,000,000 to prevent runaway programs.  
Configure it with the `YARNBALL_STEP_LIMIT` environment variable.
