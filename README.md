Casmeleon name is not a spelling error, it is a pun.

It is the combination of cameleon and asm intended as assembly source code file.

So it is an assembler that hide itself in the internet without being noticed by anyone... of course not, it is an assembler
that mimic the syntax of other assembler for different cpu architecture, a sort of universal assembler

What it has:

  * Flexible opcode syntax with checked arguments and pattern matching
  * Variable length opcode binary output, with a special instruction to output bytes in reverse order
  * Labels (global and locals)
  * Sets to define registers and similar
  * Straightforward number format
  * Inlude directive to inject other file in a specified position assembly
  * DB directive with support of quoted strings and variable length lists
  * DW directive with support of quoted strings (UTF-16 16 bit sized only) and variable lenght lists
  * Common functions (called inlines) that allow one to factor the instruction decoding logic in a small number of places
  * Advance directive to pad the generated file with minimal effort
  * Org directive to change the address without generating padding bytes

What it lacks:

  * macro in user program

The assembler require a least 2 files: a definition of the language in .casm file and a source file in any extension as long as it is text

casm definition language bnf

     <definition> ::= { <numeric format definition> | <opcode definition> | <inline definition> | <set definition> }

     <numeric format> ::= '.num' <number> <quoted stirng> <quoted string> ';'

     <set definition> ::= '.set' <identifier> '{' <identifier list> '}';

     <identifier list> ::= <identifier> [';' <identifier list>]

     <opcode definition> ::= '.opcode' <identifier> '{{' <syntax definition> '}}' '->' <block>

     <syntax definition> ::= ε | <arg format>

     <arg format> ::= { <symbol> | <arg> | <number> }

     <arg> ::= <identifier>

     <block> ::= '{' <statement list> '}'

     <statement list> ::= <statement>; [<statement list>]

     <statement> ::= <return expression> | <if statement> | <error statement> | <warning statement> | <out statement>  

     <expression> ::= <operand> | <operator> <expression> | <expression> <operator> <expression>

     <operand> ::= <number> | <identifier> | '.expr' <identifier> '(' <expression> ')' | '(' <expression> ')'

     <operator> ::= '+' | '-' | '/' | '%' | '*' | '>>' | '<<' | '&' | '|' | '^' | '~' | '&&' | '||' | '!=' | '!' | '==' | '>' | '<' | '<=' | '>='

     <if statement> ::= 'if' <expression> <block> [ 'else' <block> ]

     <error statement> ::= '.error' <source> ',' <double quoted string>
     
     <warning statement> ::= '.warning' <source> ',' <double quoted string>

     <source> ::= <identifier>

     <single quoted string> ::= "'" <string> "'"

     <double quoted string> ::= "\\"" <string> "\\""

     <string> ::= { . }

     <symbol> ::= <separator> | <operator>

     <separator> ::= '(' | ')' | '[' | ']' | '{' | '}' | ',' | '@' | ';' | '#' | ':' 

     <number> ::= <binary number> | <decimal number> | <octal number> | <hexadecimal number>

     <binary number> ::= '0b' { <binary digit> }

     <binary digit> ::= '0' | '1'

     <decimal number> ::= { <decimal digit> }

     <decimal digit> ::= '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9'

     <hexadecimal number> ::= '0x' { <hexadecimal digit> }

     <hexadecimal digit> ::= <decimal digit> | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f'

     <identifier> ::= <char> { <char> | <decimal digit> }

char is anything that is not a separator an operator or a digit

Comment marker is '//' or /* */. Comment are removed before parsing

Operators precedences are not included in the grammar to keep it short  
This is their precedence:  

  * 0: "||", "&&"
  * 1: "!=", "==", ">=", "<=", "<", ">"
  * 2: "<<", ">>", "+", "-", "|"
  * 3: "*", "/", "%", "^", "&"
  * 4: "!", "~"
 
For examples watch docs/grammar.md

Program file is a list of labels, opcodes and directives, labels are local or global  
Local labels are normal labels with '.' prefix. For example '.done' is a local label  
Local labels can be called outside their scope, example:  

     _f1:  
        ...  
    .loop:  
        ...  
    .ret:  
        ret  
  
    _f2:  
        jmp _f1.loop  

Include file example:

    .include "fileName.s"  

Advance to address example:

    .advance 5000 ; advance to the address 5000

Store bytes or words:

    .db "My list for the supermarket: even emojii are supported", 1, 20, 0x0A, 0x0D
    .dw "String of UTF-16 runes cut to 16 bit", 10, 20, 5000, 24678

Comment marker for user program is semicolon ';'

Flags and usage

casmeleon.exe -lang=lang-name file  
debug flag was provided in the v1, but was temporanely removed in v2, pending further reorganization of debug experience  
output file name is file - extension + .bin  

"Program oscillation" message mean that there was some symbol that wasn't known when first referenced (e.g. future labels) or that the subsequent reassemble list caused some of the symbol to change their address. In comparison of the last version there are internally guards that trigger a partial re-evaluation of the input after a change of address for a referenced symbol. Performance are strictly better, because the precedent version iterated the whole source multiple time until the outut was stable. In fixed encoding instruction set it is guaranteed to complete in 2 passes (1° pass whole source, 2° pass triggered revaluations), more complex instructions set encodings can require a few more passes. 
