

_ix0:                               .dw 0  ; index base 0 (caller saved)
_ix1:                               .dw 0  ; index base 1 (caller saved)
_ix2:                               .dw 0  ; index base 2 (caller saved)
_ix3:                               .dw 0  ; index base 3 (caller saved)

_arg0:                              .dw 0  ; call argument 0
_arg1:                              .dw 0  ; call argument 1
_arg2:                              .dw 0  ; call argument 2
_arg3:                              .dw 0  ; call argument 3

_g0:                                .dw 0  ; global location 0 (caller saved)
_g1:                                .dw 0  ; global location 1 (caller saved)
_g2:                                .dw 0  ; global location 2 (caller saved)
_g3:                                .dw 0  ; global location 3 (called saved)

_b0:                                .dw 0  ; temp 0 (callee saved)
_b1:                                .dw 0  ; temp 1 (callee saved)
_b2:                                .dw 0  ; temp 2 (callee saved)
_b3:                                .dw 0  ; temp 3 (callee saved)  

_ra:                                .dw 0  ; return address, used in pseudo instruction RET alias "JPI _ra"
_sp:                                .dw 0  ; stack pointer


_const_0001:                        .dw 1
_const_0002:                        .dw 2
_const_0004:                        .dw 4
_const_000B:                        .dw 11
_const_00FF:                        .dw 0xFF
_const_8000:                        .dw 0x8000

                
_ioptr:                             .dw 0xFFF0  ; ioptr base
_errno:                             .dw 0

_pow2_units:                        .dw 1, 2, 4, 8, 6, 2, 4, 8, 6, 2     ; bit to units
_pow2_tens:                         .dw 0, 0, 0, 0, 1, 3, 6, 2, 5, 1     ; bit to tens (0x0-0xA: 0, 0xB-0xF: nibble - 0xA)
_pow2_hundreds:                     .dw 0, 0, 0, 0, 0, 0, 0, 1, 2, 5     ; bit to hundreds     
    


__fatal_error:                      STA _errno
                                    JMP __fatal_error