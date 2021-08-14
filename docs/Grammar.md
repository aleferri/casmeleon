# Grammar with examples  

Assembly language specification for our code is written in the casm file, *casm* is a DSL that allow us to do pattern matching for the opcodes, mimicking a real assembler  

## Sets  
Sets are the base of the opcode format matching and set checking in the *casm* file  

    // Example:  
    .set Registers {  
        A;  
        B;  
    }  
    // this create a set named "Registers" that contains both a register called A (with implicit value of 0) and register called B (with implicit value of 1)  

## Numbers  
Assembly languages like their numbers, be it in the form with the final 'h' or the traditional form of 0x ; they use a lot of different format, especially for masks.  
In Casemeleon v1 there was a poor capability for handling numbers, the parser was mostly an hack. Now there is a easier version of number format definition
.hex, .bin, etc are no longer needed.  

    //Number format example ".num" keyword ; base in decimal (2, 3, 4, 5, etc) ; "QuotedPrefix" ; "QuotedSuffix" ; Example:  
    .num 16 "0x" "";  

    //What about my wonderful 6502 code that use the $ to identify hexadecimal?  
    .num 16 "$" "";  
    //See? Easy. Just don't put spaces inside the string or the his wrath will come upon you ("his" being the parser)  

## Common Code between opcodes  
In the beginning there was nothing. And the exact number of one person (that is: me) was forced by the programmer to rewrite the same code over and over again.  
And the user was mad at the programmer! But since the programmer was himself the only user, in truth the user was mad at himself!  
So it came Casmeleon v2 and the user avoided the development of a severe form of split personality.  
And so they were: Inlines. In others, less interesting languages, formulated by unoriginal authors, these would be called "private fuctions" (see how the name immediately destroy the aura of misteriousness around the concept).  
Let's stop digressing further and let's dive inside the wonderful worlds of inlines!

    // Note: first of all you say you want to declare an "inline", and you do it with a keyword: ".inline",  
    // i hope no one was surprised by that  
    .inline HELLO_WORLD                 // Ok, tecnically speaking this is not wat the .inline was meant to be, 
                                        // but i couldn't resist on the temptation of an "hello, world!" tutorial  
    .with ( dummy : Ints )  ->  {       // Parameter list, leave empty if you don't have parameters, note "->" separator, 
                                        // it seems to be the latest trend, note that the dummy parameter is useless    
        .error dummy, "Hello, World!" ; // And here we will give a big welcome to the World! Since this is an error we  
                                        // need to send the wrong parameter so that the user can  see what is wrong with the program  
                                        // Of course, if an error is issued, the compilation stop and no further processing is done.  
                                        // Since the .error is unavoidable if we call the inline HELLO_WORLD from one of the opcodes,  
                                        // the very next second after our "hello" is printed, we will crush the binary    
                                        // aspirations of our poor user. So this is essentially a big f**k you to the assembly 
                                        // programmer, as such i don't recommend writing code like this in a casm file  
        .return 0 ;                     // This return will never be executed, of course  
    }  

    // Ok, now that our holy duty of an "Hello, World" is done, we can skip to useful user interaction and  
    // doing actual processing with the given parameters. Next example please!  
    .inline DO_DIV  
    .with ( q : Ints, d : Ints ) -> {  
        .if d == 0 {                                // Useless parens are useless. The ".if" can be followed by an ".else"  
            .error d, "Divisor cannot be zero";     // Finally a interesting user interaction, we give the user an error  
                                                    // because one of opcodes called this inline with invalid parameters  
        }  
        .return q / d;                              // A return that will be actually executed  
    }  

    // Just to be sure, let's make another example, a dive in the realm of instructions encoding   
    // We may have designed a backward instruction set full of prefixes and suffixes bytes, so we don't want to recode  
    // the logic for the suffix in every opcode. The best solution would be dropping the instruction set and burn it with fire.
    // But for some obscure reason we really want to encode it: so we let the opcode capture the pattern, but leave the processing to  
    // the specified inline. It is purely academical of course: no one would dare to invent something so bizantine    
    .inline SIB_BYTE
    .with ( base : Register, index : Register, scale : Ints ) -> {  
        .return (scale << 6) + ( index << 3 ) + base;     // scale[7:6], index[5:3], base[2:0]
    }  


## Opcodes

Finally we are here, the last element of our beautiful color palette. The Opcode. An hacked version of perl like regex, made language construct.  
Between "{{" and "}}", called "double curly braces for pattern matching", abbreviated: "ugly delimiters" for friends and foe alike  

    // The opcode format is:  
    .opcode move {{ dest, mod ptr [ segm:base + index * scaled ] }}  
    .with ( dest : Register, mod : x86Modifiers, ptr : PtrKeyword, segm : Segments, base : Register, index : Register, scaled : Ints ) -> {  
        .if scaled != 1 && scaled != 2 && scaled != 4 && scaled != 8 {  
            .error scaled, "Only 1 | 2 | 4 | 8 allowed as scale";   
        }  
        .out [ .expr SEGMENT_PREFIX(segm), 0x8B, .expr MAKE_RM( dest, + 0b100 ), .expr MAKE_SIB(base, index, scaled) ];  
    }   

Inlines are called from opcodes using the ".expr" syntax. I did this because i was lazy, i didn't want to check for "open round parens" before jumping
in the "inline call" parser branch,

## Conclusions  

Congratulations, now you know the grammar, you are free to paint casmeleon with every possible color mimicking every assembler out of here   
