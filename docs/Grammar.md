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
Assembly languages like their numbers, be it in the form with the final 'h' or the traditional form of 0x they use a lot of different format, especially for masks  
In Casemeleon v1 there was a poor capability for handling numbers, the parser was mostly an hack. Now there is a easier version of number format definition  
.hex, .bin, etc are no longer needed.  

    //Number format example ".num" keyword ; base in decimal (2, 3, 4, 5, etc) ; "QuotedPrefix" ; "QuotedSuffix" ; Example:  
    .num 16 "0x" "";  

    //What about my wonderful 6502 code that use the $ to identify decimal?  
    .num 16 "$" "";  
    //See? Easy. Just don't put spaces inside the string or the his wrath will come upon you ("his" being the parser)  

## Common Code between opcodes  
In the beginning there was nothing. And the exact number of one person (that is: me) was forced by the programmer to rewrite the same code over and over again.  
And the user was mad at the programmer! But since the programmer was himself the only user, in truth the user was mad at himself!  
So it came Casmeleon v2 and the user avoided the development of a severe form of split personality.  
And so they were: Inlines. In others, less interesting languages, formulated by unoriginal authors, these would be called "private fuctions" (see how the name immediately destroy the aura of misteriousness around the concept).  
Let's stop digressing further and let's dive inside the wonderful worlds of inlines!

    // Note: first of all you say you want to declare an "inline", and you do it with a keyword: ".inline", i hope no one was surprised by that  
    .inline HELLO_WORLD                 // Ok, tecnically speaking this is not wat the .inline was meant to be, 
                                        // but i couldn't resist on the temptation of an "hello, world!" tutorial  
    .with ( dummy : Ints )  ->  {       // Parameter list, leave empty if you don't have parameters, rust inspired "->" separator, it seems to be the latest trend  
                                        // The dummy parameter here is not really used, but i need it a few line later  
        .error dummy, "Hello, World!" ; // And here we will give a big welcome to the World! Since this is an error we need to send the wrong parameter so that the user can  
                                        // see what is wrong with the program  
                                        // Of course, if an error is issued, the compilation stop and no further processing is done. Since the .error is unavoidable if we   
                                        // call the inline HELLO_WORLD from one of the opcodes, the very next second after our "hello" is printed, we will crush the binary    
                                        // aspirations of our poor user. So this is essentially a big f**k you to the assembly programmer, as such i don't recommend   
                                        // writing code like this in a casm file that is meant to be shared  
        .return 0 ;                     // This return will never be executed, of course  
    }  

    // Ok, now that our holy duty of an "Hello, World" is done, we can skip to useful user interaction and doing actual processing with the given parameters  
    // Next example inline  
    .inline DO_DIV  
    .with ( q : Ints, d : Ints ) -> {  
        .if d == 0 {                                // Useless parens are useless. The ".if" can be followed by an ".else", totally unexpected right?  
            .error d, "Divisor cannot be zero";     // Finally a interesting user interaction, we give the user an error because one of opcodes called this inline with   
                                                    // invalid parameters  
        }  
        .return q / d;                              // A return that will be actually executed  
    }  

    // Should be clear by now the usefulness but just to be sure, let's make another example, a dive in the realm of instruction encoding   
    // We may have designed an ADD that can be encoded in 10 different modes, but maybe the field that say "do add", the opcode syntax will capture the pattern,  
    // but in this case is better for them to pass the actual encoding to the inlines, unless the field are shuffled around like astronauts during training  
    .inline ENCODE_ADD
    .with ( base : Ints, dest : Register ) -> {          // Invented  
        return base + ( dest << 4 ) + ( 0x10 << 2 );     // base opcode family + destination register + add operation  
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
        .out [ SEGMENT_PREFIX(segm), 0x8B, MAKE_RM( dest, + 0b100 ), .expr MAKE_SIB(base, index, scaled) ];  
    }  
    // Here we again can see how inlines are useful, unless you want to rewrite the logic to encode the SIB and R/M byte every instruction of course  
    // Good luck when a bug arises  

## Conclusions  

Congratulations, now you know the grammar, you are free to paint casmeleon with every possible color mimicking every assembler out of here   