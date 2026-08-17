 

; memcopy: src, dest, n
__memcopy:      LDA     _arg2
                JZE     .end
.loop:          STA     _arg2
                LDI     _arg0
                STI     _arg1
                INC     _arg0
                INC     _arg1
                LDA     _arg2
                SUB     _const_0001
                JNZ     .loop
.end:           JPI     _ra
                
                
; memset: c: int, dest: ptr, n: int
__memset:       LDA     _arg2
                JZE     .end
.loop:          STA     _arg2
                LDA     _arg0
                STI     _arg1
                INC     _arg1
                LDA     _arg2
                SUB     _const_0001
                JNZ     .loop
.end:           JPI     _ra
                

; split binary number in tens and units
; bin2str: bin: char, dest: ptr
__split_bin:    LIT     0x7F
                AND     _arg0
                STA     _arg0
                LIT     _pow2_units         ; prepare units table
                STA     _ix0
                LIT     _pow2_tens          ; prepare tens table
                STA     _ix1
                LIT     0
                STA     _g0                 ; units
                STA     _g1                 ; tens
.loop:          LDA     _arg0
                AND     _const_0001         ; check if n & 1 is true, 
                JZE     .skip               ; otherwise skip to the next bit
                LDI     _ix0                ; read units table entry at [ix0]
                ADD     _g0                 ; add to the units
                STA     _g0
                LDI     _ix1                ; read tens table entry at [ix1]
                ADD     _g1                 ; add to the tens
                STA     _g1
.skip:          INC    _ix0                 ; advance units table pointer
                INC    _ix1                 ; advance tens table pointer
                SHR    _arg0                ; shift right n
                LDA    _arg0
                JNZ    .loop                ; test not zero
.adjust_test:   LDA    _g0
                SUB    _const_000B          ; subtract 11 to units
                AND    _const_8000          ; test sign bit
                JNZ    .write               ; no overflow, skip adjustement
.adjust:        LDA    _g0
                SUB    _g3
                STA    _g0
                INC    _g1
                JMP    .adjust_test
.write:         LDA    _g1                  ; write tens
                STI    _arg1
                INC    _arg1                ; increment pointer
                LDA    _g0                  ; write units
                STI    _arg1
                INC    _arg1
                JPI    _ra                  ; return
                