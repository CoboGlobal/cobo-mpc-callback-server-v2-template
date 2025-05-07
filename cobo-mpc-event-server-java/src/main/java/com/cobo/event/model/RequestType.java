package com.cobo.callback.model;

import lombok.Getter;

@Getter
public enum RequestType {
    TYPE_PING(0),
    TYPE_KEY_GEN(1),
    TYPE_KEY_SIGN(2),
    TYPE_KEY_RESHARE(3);

    private final int value;

    RequestType(int value) {
        this.value = value;
    }

    public static RequestType fromValue(int value) {
        for (RequestType type : RequestType.values()) {
            if (type.value == value) {
                return type;
            }
        }
        throw new IllegalArgumentException("Unknown RequestType value: " + value);
    }
}
