package com.cobo.callback.model;

import lombok.Getter;

@Getter
public enum CurveType {
    SECP256K1(0),
    ED25519(2);

    private final int value;

    CurveType(int value) {
        this.value = value;
    }
}
