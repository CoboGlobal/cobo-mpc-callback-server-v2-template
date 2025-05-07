package com.cobo.callback.model;

import lombok.Getter;

@Getter
public enum TssProtocol {
    UNKNOWN_PROTOCOL(0),
    GG18(1),
    LINDELL(2),
    EDDSA_TSS(3);

    private final int value;

    TssProtocol(int value) {
        this.value = value;
    }

}
