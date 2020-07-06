Casmeleon name is not a spelling error, it is a pun.

It is the combination of cameleon and asm intended as assembly source code file.

So it is an assembler that hide itself in the internet without being noticed by anyone... of course not, it is an assembler
that mimic the syntax of other assembler for different cpu architecture, a sort of universal assembler

What it has:

  * flexible opcode syntax
  * variable length opcode binary output
  * labels (global and locals)
  * enums to define registers and similar
  * flexible number format
  * include other file in assembly part
  * variable list arguments (as byte list) in opcodes

What it lacks (lot of things):

  * macro in user program
  * variable list arguments for opcodes
  * expression in program
  * few other things

The assembler require a least 2 files: a definition of the language in .casm file and a source file in any extension as long as it is text

casm definition language bnf

     <definition> ::= { <numeric format> | <opcode definition> | <enum definition> }

     <numeric format> ::= .number <number base> <position> <single quoted string> '\n'

     <number base> ::= .hex | .bin | .oct | .dec

     <position> ::= prefix | suffix

     <enum definition> ::= .enum <identifier> '{' <identifier list> '}';

     <identifier list> ::= <identifier> [, <identifier list>]

     <opcode definition> ::= .opcode <identifier> <syntax definition> -> <block>

     <syntax definition> ::= ε | <arg format>

     <arg format> ::= { <symbol> | <arg> | <number> }

     <arg> ::= <identifier>

     <block> ::= '{' <statement list> '}'

     <statement list> ::= <statement>; [<statement list>]

     <statement> ::= <deposit> <expression> | <if statement> | <error statement> | <loop statement>  

     <deposit> ::= .db | .dw | .dd

     <expression> ::= <operand> | <operator> <expression> | <expression> <operator> <expression> | '(' <expression> ')'

     <operand> ::= <number> | <identifier>

     <operator> ::= '+' | '-' | '/' | '%' | '*' | '>>' | '<<' | '&' | '|' | '^' | '~' | '&&' | '||' | '!=' | '!' | '==' | '>' | '<' | '<=' | '>=' | .in

     <if statement> ::= if <expression> <block> [ else <block> ]

     <error statement> ::= .error <source> <double quoted string>

     <loop statement> ::= for <identifier> until <expression> <deposit> <expression>

     <source> ::= <identifier>

     <single quoted string> ::= "'" <string> "'"

     <double quoted string> ::= "\\"" <string> "\\""

     <string> ::= { . }

     <symbol> ::= <separator> | <operator>

     <separator> ::= '(' | ')' | '[' | ']' | '{' | '}' | ',' | '@' | '$' | ';' | '#'

     <number> ::= <binary number> | <decimal number> | <octal number> | <hexadecimal number>

     <binary number> ::= 0b { <binary digit> }

     <binary digit> ::= 0 | 1

     <decimal number> ::= { <decimal digit> }

     <decimal digit> ::= 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9

     <octal number> ::= 0o { <octal digit> }

     <octal digit> ::= 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7

     <hexadecimal number> ::= 0x { <hexadecimal digit> }

     <hexadecimal digit> ::= 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | A | B | C | D | E | F | a | b | c | d | e | f

     <identifier> ::= <char> { <char> | <decimal digit> }

char is anything that is not a separator an operator or a digit

Comment marker is '//'. Comment are removed before parsing

Operators precedences are not included in the grammar to keep it short  
This is their precedence:  

  * 0: "||", "&&"
  * 1: "!=", "==", ">=", "<=", "<", ">", ".in"
  * 2: "<<", ">>", "+", "-", "|"
  * 3: "*", "/", "%", "^", "&"
  * 4: "!", "~"
 
Some example:  

    //org implementation example  
    .opcode org addr -> {  
        for i until addr db 0; //i start from this_address and loop until the desired address padding with 0  
    }  

    //add immediate  
    .opcode add #imm8 -> {  
        if imm8 \> 255 {  
            .error imm8 "immediate must be less than 256";  
        }  
        .db 0x10;  
        .db imm8;  
    }  
    //space after symbols are stripped, so "add # 26" is the same as "add #26"  

    //this syntax can also be used for meta-opcode, example:  
    .opcode db imm8 -> {  
        .db imm8;  
    } //db now deposit values in the program  

    //number format  
    .number .hex suffix 'h' //in the user program all number ending with h are hexadecimal  

Program file is a list of labels and opcodes, labels are local or global  
Local labels are normal labels with '.' prefix so .done is a local label  
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

Comment marker for user program is semicolon ';'

Flags and usage

casmeleon.exe -lang=lang-name -debug=true/false file  
debug flag is optional  
output file name is file - extension + .bin  
debug cause the individual opcode to be print to video in the form address: opcode-name expanded args -> list of bytes
