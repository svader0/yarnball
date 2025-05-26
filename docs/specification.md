# Yarnball Specification
## 1. Introduction

**Yarnball** is an esoteric, stack-based programming language whose every instruction reads like a crochet pattern.  You can literally hand this “pattern” to a crocheter and they’ll work through it—though what they produce won’t be much more than a random string of stitches.  Under the hood, each stitch manipulates a stack of integers.

---


## 2. Lexical & Structural Grammar


* **Tokens**
  * **Stitches**: `ch`, `sc`, `dc`, `hdc`, `tr`, `cl`, `inc`, `dec`, `add` (`bob`), `mul` (`dc`), `div` (`tr`), `sub` (`hdc`), `mod` (`cl`), `dup`, `swap`, `sl st`, `yo`, `pic`, `rep`, `FO`

  * **Literals**: non-negative integers (e.g. `0`, `42`)
  * **Grouping**: `(`, `)`, commas (`,`), optional labels (`Row N:`, `Round N:`)
  * **Comments**: start with `#`, go to end-of-line

* **Program** = sequence of lines; each line can begin with an **ignored** crochet label (`Row 1:`, etc.), then one or more instructions or literals separated by whitespace or commas.

* **`rep` syntax**

  ```
  rep [<count>] (<instr1>, <instr2>, …)
  ```

  * If `<count>` is present, it doesn't have to be a literal; it can be any expression that evaluates to a non-negative integer.
  * If omitted, the interpreter **pops** a count from the stack.

* **Termination**

  * The `FO` stitch immediately stops execution.
  * Falling off the end (no more instructions) also halts.
  * Only `yo` and `pic` produce visible output; `FO` does not implicitly print.

---

  

## 3. Instruction Set

### 3.1 Stack Manipulation

| Stitch         | Mnemonic | Effect                                                    |
| -------------- | -------- | --------------------------------------------------------- |
| Chain          | `ch n`   | Push literal `n` onto the stack.                          |
| Slip Stitch      | `sl st`    | Duplicate the top value (peek & push a copy).             |
| Swap           | `swap`   | Pop `a,b` then push in reverse order (`a` $\rightarrow$ top, `b` $\rightarrow$ next). |
| Single Crochet | `sc`     | Pop & discard the top value.                              |

### 3.2 Arithmetic Operations

All pop the required operands, compute, then push the integer result.

| Stitch          | Mnemonic | Operation                 | Example                       |
| --------------- | -------- | ------------------------- | ----------------------------- |
| Bobble          | `bob`    | Pop `a,b`; push `b + a`   | `ch 2 ch 3 bob` ⇒ stack `[5]` |
| Half-Double Cro | `hdc`    | Pop `a,b`; push `b − a`   | `ch 5 ch 2 hdc` ⇒ `[3]`       |
| Double Crochet  | `dc`    | Pop `a,b`; push `b × a`   | `ch 3 ch 4 dc` ⇒ `[12]`      |
| Treble Crochet  | `tr`    | Pop `a,b`; push `⌊b ÷ a⌋` | `ch 8 ch 3 tr` ⇒ `[2]`       |
| Cluster (mod)   | `cl`    | Pop `a,b`; push `b mod a` | `ch 5 ch 12 cl` ⇒ `[2]`      |
| Increase        | `inc`    | Pop `x`; push `x + 1`     | `ch 4 inc` ⇒ `[5]`            |
| Decrease        | `dec`    | Pop `x`; push `x − 1`     | `ch 4 dec` ⇒ `[3]`            |

### 3.3 I/O

| Stitch    | Mnemonic | Effect                                      |
| --------- | -------- | ------------------------------------------- |
| Yarn Over | `yo`     | Pop `n`; print decimal `n` + newline.       |
| Picot     | `pic`    | Pop `n`; print ASCII `chr(n)` (no newline). |

### 3.4 Control Flow: `rep`
```
rep [<count>] (<instr1>, <instr2>, …)
```

* **Literal count**: `rep 5 (…)` runs 5×.
* **Dynamic count**: `rep (…)` pops an integer from the stack and loops that many times.
* Non-positive counts ⇒ zero iterations.

---

## 4. (Optional) Function-Patterns

To avoid repetition, you can **define** and **invoke** named blocks:

```bnf

<pattern-def> ::= "pattern" <Name> [ "(" <param1>,… ")" ] "=" "(" <instr-list> ")"
<invoke>      ::= "use"     <Name> [ "(" <arg1>,… ")" ]

```

* **Definition**: parameters become local stack variables; arguments are pushed before execution.
* **Invocation**: executes the block in-place.
* Patterns can be called within other patterns, allowing for nested definitions.

**Example**:

```
pattern print_str(count) = (
  rep(count) ( pic )
)

# push 5 chars 'H','i','!','\n','\0'
ch 72 pic ch 105 pic ch 33 pic ch 10 pic ch 0 pic  
ch 5           # loop count
use print_str(5)
FO

```

---

## 5. Semantics & Error Handling

* **Stack underflow**: any pop on empty stack ⇒ runtime error, halt.
* **Divide/mod by zero** ⇒ runtime error.
* **Malformed `rep`** (missing parens or bad count) ⇒ parse error.
* **Unknown stitch** ⇒ parse error, listing the bad token + line.

---
## 6. Examples

### 6.1 “Hello, World!”
```
Row 1:  ch 72 pic    # 'H'
Row 2:  ch 101 pic   # 'e'
Row 3:  ch 108 pic   # 'l'
Row 4:  ch 108 pic   # 'l'
Row 5:  ch 111 pic   # 'o'
Row 6:  ch 44 pic    # ','
Row 7:  ch 32 pic    # ' '
Row 8:  ch 87 pic    # 'W'
Row 9:  ch 111 pic   # 'o'
Row 10: ch 114 pic   # 'r'
Row 11: ch 108 pic   # 'l'
Row 12: ch 100 pic   # 'd'
Row 13: ch 33 pic    # '!'
Row 14: yo          # print “!” + newline
Row 15: FO          # end

```

  

#### Condensed with loop

```
ch 72 ch 101 ch 108 ch 108 ch 111 ch 44 ch 32 ch 87 ch 111 ch 114 ch 108 ch 100 ch 33
ch 13       # push count
rep ( pic )
FO
```


### 6.2 Factorial

```

ch 5       # n
ch 1       # acc
rep 5 (
  dup      # keep n for next loop
  sc       # pop copy of n
  mul      # acc ← acc × n
  dec      # n ← n − 1
)
yo         # print 120
FO         # terminate
```
---
