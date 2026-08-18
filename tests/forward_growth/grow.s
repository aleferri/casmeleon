        .advance 0x0100
start:  LD      #big         ; forma corta alla prima passata, lunga alla seconda
        JSR     target       ; deve puntare al primo byte di target
target: RET
        .advance 0x0200
big:    RET
        .dw start, target, big
        .db 1, 2, start
