# Yarnball Language Specification
## 1. Introduction

**Yarnball** is an esoteric, stack-based programming language whose every instruction reads like a crochet pattern.  You can literally hand this “pattern” to a crochet artist and they’ll work through it. Though, what they produce won’t be much more than a random string of stitches.  

Under the hood, each instruction manipulates a single stack of integers. This is to be able to perform complex operations by combining a minimal set of very simple instructions.

For program examples, see [here](https://github.com/svader0/yarnball/tree/master/examples).

---

## Instruction Set

### Simple Instructions

- **ch (chain)**  
  *Usage:* `ch <number>`  
  *Behavior:* Converts the provided argument to an integer and pushes it onto the stack.

- **pic (picot stitch)**  
  *Usage:* `pic`  
  *Behavior:* Pops the top of the stack and prints the corresponding character (based on its ASCII value).  
  *Error:* Fails if the stack is empty.

- **yo (yarn over)**  
  *Usage:* `yo`  
  *Behavior:* Pops the top value from the stack and prints it as a number.  
  *Error:* Fails if the stack is empty.

- **fo (finish off)**  
  *Usage:* `fo`  
  *Behavior:* Halts the program immediately.

- **sc (slip stitch)**  
  *Usage:* `sc`  
  *Behavior:* Pops (discards) the top value from the stack.  
  *Error:* Fails if the stack is empty.

- **dc (double crochet)**  
  *Usage:* `dc`  
  *Behavior:* Pops the top two values, multiplies them, and pushes the product onto the stack.  
  *Error:* Fails if the stack has fewer than two values or if the second value is zero (to avoid issues noted in the implementation).

- **bob (bobble stitch)**  
  *Usage:* `bob`  
  *Behavior:* Pops the top two values, adds them, and pushes the sum onto the stack.  (n1, n2) -> (n1 + n2)
  *Error:* Fails if the stack has fewer than two values.

- **hdc (half double crochet)**  
  *Usage:* `hdc`  
  *Behavior:* Pops the top two values, subtracts the first popped value from the second, and pushes the result.  
  *Error:* Fails if the stack has fewer than two values.

- **tr (treble crochet)**  
  *Usage:* `tr`  
  *Behavior:* Pops the top two values, divides the second by the first, and pushes the quotient.  
  *Error:* Fails if the stack has fewer than two values or if division by zero occurs.

- **cl (cluster stitch)**  
  *Usage:* `cl`  
  *Behavior:* Pops the top two values, computes the modulo (second % top), and pushes the result.  
  *Error:* Fails if the stack has fewer than two values or if modulo by zero occurs.

- **sl st (slip stitch)**  
  *Usage:* `sl st`  
  *Behavior:* Duplicates the top value of the stack.  
  *Error:* Fails if the stack is empty.

- **swap**  
  *Usage:* `swap`  
  *Behavior:* Pops the top two elements and pushes them back in reverse order, effectively swapping them.  
  *Error:* Fails if the stack has fewer than two values.

- **inc (increase)**   
  *Usage:* `inc`  
  *Behavior:* Pops the top value, increments it by one, and pushes the result.  
  *Error:* Fails if the stack is empty.

- **dec (decrease)**  
  *Usage:* `dec`  
  *Behavior:* Pops the top value, decrements it by one, and pushes the result.  
  *Error:* Fails if the stack is empty.

---

### Comparison Instructions

These instructions assume that two values can be popped from the stack to perform a comparison.

- **>**  
  *Usage:* `>`  
  *Behavior:* Pops two values; if the second is greater than the first, pushes `1` (true) otherwise `0` (false).  
  *Error:* Fails if the stack has fewer than two values.

- **<**  
  *Usage:* `<`  
  *Behavior:* Pops two values; if the second is less than the first, pushes `1` (true) otherwise `0` (false).  
  *Error:* Fails if the stack has fewer than two values.

- **eq**  
  *Usage:* `eq`  
  *Behavior:* Pops two values; if they are equal, pushes `1` (true) otherwise `0` (false).  
  *Error:* Fails if the stack has fewer than two values.

- **neq**  
  *Usage:* `neq`  
  *Behavior:* Pops two values; if they are not equal, pushes `1` (true) otherwise `0` (false).  
  *Error:* Fails if the stack has fewer than two values.

---

### Stack Manipulation

- **turn**  
  *Usage:* `turn`  
  *Behavior:* Rotates the top three elements of the stack.  
  *Details:*  
    1. Pops the top three values (let's call them _top_, _second_, and _third_ in the order they are popped).  
    2. Pushes them back in the order: _second_, then _top_, then _third_.  
    This effectively rotates the top three items of the stack, similar to FORTH'S `rot` instruction. 
  *Error:* Fails if the stack contains fewer than three values.

---

### Block and Control Flow Instructions

- **rep**  
  *Usage:*  
    - With a count expression: `*[ ... block instructions ... ]; rep from * [count] times`  
    - Without a count expression: `*[ ... block instructions ... ]; rep from *` (in which case the count is popped from the stack)  
  *Behavior:*  
    Executes the provided block of instructions a specified number of times. If a count expression is provided, it is parsed into an integer; otherwise, the evaluator pops the count from the stack.  
  *Error:* Fails if the provided count is invalid or if the stack underflows when a count is needed.

- **if**  
  *Usage:* `if [ ... if-body ... ] else [ ... else-body ... ] end`  
  *Behavior:*  
    Pops the top value from the stack, which should be either `0` (false) or `1` (true). If the value is `1`, the evaluator executes the if-body; if `0`, the else-body is executed.  
  *Error:* Fails if the condition is not `0` or `1` or if the stack is empty.

  Example:  
  ```yarnball
  if 
    ch 65 pic 
  else 
    ch 66 pic 
  end
  ```
  
### Stitches (Functions)

- **use**
  *Usage:* `use <subpattern name>`  
  *Behavior:*  
    Invokes a predefined subpattern. The evaluator retrieves the subpattern by name and executes its sequence of instructions.  
  *Error:* Fails if the subpattern is not defined.

- **Subpattern Definitions**  
  *Usage:* `subpattern <subpattern name> = ([ ... instructions ... ]) end`  
  *Behavior:*  
    Defines a reusable subpattern that can be invoked later using the `use` instruction. By convention, subpattern names should be defined under the `STITCH GUIDE:` header. Names should be unique and contain no special characters. Names are also case-insensitive.
  *Error:* Fails if the subpattern name is invalid or if the definition is malformed.
    - Example:  
  ```yarnball
    subpattern printHello = (
      ch 72 pic  # H
      ch 101 pic # e
      ch 108 pic # l
      ch 108 pic # l
      ch 111 pic # o
    )
  ```

---