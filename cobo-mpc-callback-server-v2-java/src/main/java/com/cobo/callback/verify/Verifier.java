package com.cobo.callback.verify;

import com.cobo.waas2.model.TSSCallbackRequest;

public interface Verifier {
    String verify(TSSCallbackRequest request);
}
