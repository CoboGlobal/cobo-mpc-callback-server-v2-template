package com.cobo.callback.verify;

import com.cobo.callback.model.Request;

public interface Verifier {
    String verify(Request request);
}
