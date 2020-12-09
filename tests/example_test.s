.advance 0xE000
            
            ; halt ( soft_reset : &fn() ) -> void
_halt:      ST      A, 1
            LH      X, #0xFF
            ADD     X, #0xFF
            LD      Y, #1
            ST      A, (X + 0)
            
            ; wait ( ) -> void
_wait:      LD      A, 16
            ST      A, 1
            ST      A, 2
            ST      A, 3
            ST      A, 4
            ST      A, 5
            ST      A, 6
            ST      A, 7
            LH      X, #0xFF
            ADD     X, #0xFF
            LD      Y, #1
            ST      Y, (X, 0)

_self_test: LD      X, #0
            SUB     X, #255
            LD      A, 0
            ST      A, (X, 4)
            ST      A, (X, 5)
            ST      A, (X, 6)
            ST      A, (X, 7)
            RET

_bcd_binary: .db 0, 0, 1, 0, 2, 0, 3, 0, 4, 0, 5, 0, 6, 0, 7, 0, 9, 0, 0, 1, 1, 1, 1, 2, 1, 3, 1, 4, 1, 5

            ; dec2bcd ( n : u16, ptr : &u16 ) -> void
_dec2bcd:   LD      Y, 4            ; i = 4; temp = n;
.L0:        ST      A, 18           ; while ( i != 0 ) {
            AND     A, #15          ;       nibble = temp & 15;
            ASL     A               ;       nibble <<= 2;
            ST      A, (X, 0)       ;       *ptr = nibble;
            ADD     X, #1           ;       ptr++;
            SUB     Y, #1           ;       i--;
            LD      A, 18           ;       temp >>= 4;
            LSR     A
            LSR     A
            LSR     A
            LSR     A
            ST      A, 18
            BNE     Y, .L0          ; }
            SUB     X, #4           ; ptr -= 4;
            ST      Y, 19           ; carry = 0;
            LD      Y, #4           ; i = 4;
            ST      Y, 18           ; 
            
            ; We have 4 nibbles multiplied by 2
                                    ; while ( i != 0 ) {
.L1:        LD      Y, (X, 0)       ;     conv_ptr = (uintptr*) *ptr;
            LD      A, (Y, 1)       ;     next_carry = conv_ptr[1];
            ST      A, 20           ;     //store for later
            LD      A, (Y, 0)       ;     digit = conv_ptr[0];
            ADD     A, 19           ;     digit += carry;
            LD      Y, 20           ;     //load next carry
            CMP     A, #10          ;     if ( digit > 10 ) {
            ICY     Y               ;         next_carry++;
            ST      Y, 19           ;         carry = next_carry;
            SUB     Y, 20           ;     
            BEQ     Y, .NO_CARRY
            SUB     A, 10           ;         digit -= 10;
.NO_CARRY:  LD      Y, 18           ;     }
            SUB     Y, #1           ;     i--;
            ST      A, (X, 0)       ;     *ptr = digit;
            ADD     X, #1           ;     ptr++;
            BNE     Y, .L1          ; }
            RET                     ; return;

            ; print7 ( a : &u16 )
_print7:    LD      X, #0
            SUB     X, #256
            LD      Y, (A, 0)
            ST      Y, (X, 4)
            LD      Y, (A, 1)
            ST      Y, (X, 5)
            LD      Y, (A, 2)
            ST      Y, (X, 6)
            LD      Y, (A, 3)
            ST      Y, (X, 7)
            

            ; memset ( z: u16, a : &u16, c : u16 ) -> void
_memset:    BEQ     Y, .L1
.L0:        ST      A, (X, 0)
            ADD     X, #1
            SUB     Y, #1
            BNE     Y, .L0
.L1:        RET

            ; memcpy ( a : &u16, b: &u16, c : u16 ) -> void
_memcpy:    BEQ     Y, .L1
            ST      L, 16
.L0:        LD      L, (A, 0)
            ST      L, (X, 0)
            ADD     A, #1
            ADD     X, #1
            SUB     Y, #1
            BNE     Y, .L0
            LD      L, 16
.L1:        RET

            ; lsl_u16( a: u16, f : u16 ) -> u16
_lsl_u16:   BEQ     X, .L1
.L0:        ASL     A
            SUB     X, #1
            BNE     X, .L0
.L1:        RET

            ; lsr_u16( a : u16, f : u16 ) -> u16
_lsr_u16:   BEQ     X, .L1
.L0:        LSR     A
            SUB     X, #1
            BNE     X, .L0
.L1:        RET

            ; mul_u16( a : u16, b : u16 ) -> u16
_mul_u16:   BNE     X, .L0
            AND     A, #0
            RET
.L0:        BNE     A, .setup
            RET
.setup:     ST      A, 17
            AND     A, #0
.L1:        AND     Y, #0
            LSR     X
            ICY     Y
            ASL     A
            BEQ     Y, .L2
            ST      A, 18
            ADD     A, 17
            ST      A, 17
            LD      A, 18
.L2:        BNE     X, .L1
            RET
            