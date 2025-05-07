package com.cobo.callback.model;

import lombok.Getter;

@Getter
public enum SignatureType {
    UNKNOWN_TYPE(0),
    ECDSA(1),
    EDDSA(2),
    SCHNORR(3);

    private final int value;

    SignatureType(int value) {
        this.value = value;
    }

}
