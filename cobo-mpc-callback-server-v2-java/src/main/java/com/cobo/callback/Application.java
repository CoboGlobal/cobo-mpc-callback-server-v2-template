package com.cobo.callback;

import java.util.HashMap;
import java.util.Map;

import com.cobo.waas2.model.*;
import com.cobo.callback.config.AppConfig;
import com.cobo.callback.service.JwtService;
import com.cobo.callback.verify.TssVerifier;
import com.cobo.callback.verify.Verifier;
import com.hellokaton.blade.Blade;
import com.hellokaton.blade.annotation.*;
import com.hellokaton.blade.annotation.request.Form;
import com.hellokaton.blade.annotation.route.GET;
import com.hellokaton.blade.annotation.route.POST;
import com.hellokaton.blade.mvc.ui.ResponseType;

import io.jsonwebtoken.JwtException;
import lombok.extern.slf4j.Slf4j;

@Slf4j
@Path
public class Application {
    public static final int STATUS_OK = 0;
    public static final int STATUS_INVALID_REQUEST = 10;
    public static final int STATUS_INVALID_TOKEN = 20;
    public static final int STATUS_INTERNAL_ERROR = 30;

    private static JwtService jwtService;
    private static AppConfig appConfig;
    private static final Verifier verifier = TssVerifier.create();

    @GET(value = "/ping", responseType = ResponseType.JSON)
    public Map<String, String> ping() {
        Map<String, String> response = new HashMap<>();
        response.put("server", appConfig.getServiceName());
        response.put("timestamp", String.valueOf(System.currentTimeMillis()));
        return response;
    }

    @POST(value = "/v2/check", responseType = ResponseType.TEXT)
    public String riskControl(@Form String TSS_JWT_MSG) {
        try {
            if (TSS_JWT_MSG.isEmpty()) {
                throw new IllegalArgumentException("Missing TSS_JWT_MSG parameter");
            }

            String requestData = jwtService.verifyToken(TSS_JWT_MSG);
            TSSCallbackRequest request = TSSCallbackRequest.fromJson(requestData);

            TSSCallbackResponse response = processRequest(request);
            String responseJson = response.toJson();
            return jwtService.createToken(responseJson);

        } catch (JwtException e) {
            return handleError(STATUS_INVALID_TOKEN, e.getMessage());
        } catch (Exception e) {
            log.error("Failed to process request", e);
            return handleError(STATUS_INTERNAL_ERROR, e.getMessage());
        }
    }

    private String handleError(int status, String message) {
        try {
            TSSCallbackResponse response = new TSSCallbackResponse();
            response.setStatus(status);
            response.setError(message);
            return jwtService.createToken(response.toJson());
        } catch (Exception ex) {
            log.error("Error handling error response", ex);
            throw new RuntimeException(ex);
        }
    }

    private static TSSCallbackResponse processRequest(TSSCallbackRequest request) {
        String error = verifier.verify(request);
        if (error != null) {
            TSSCallbackResponse response = new TSSCallbackResponse();
            response.setStatus(STATUS_INVALID_REQUEST);
            response.setRequestId(request.getRequestId());
            response.setAction(TSSCallbackActionType.REJECT);
            response.setError(error);
            return response;
        }

        TSSCallbackResponse response = new TSSCallbackResponse();
        response.setStatus(STATUS_OK);
        response.setRequestId(request.getRequestId());
        response.setAction(TSSCallbackActionType.APPROVE);
        return response;
    }

    public static void main(String[] args) {
        appConfig = AppConfig.loadConfig(args);
        log.info("Loaded configuration: {}", appConfig);

        jwtService = new JwtService(appConfig);

        String[] hostAndPort = appConfig.getEndpoint().split(":");
        int port = Integer.parseInt(hostAndPort[1]);
        log.info("Starting server on port {}", port);

        Blade.create().listen(port).start(Application.class, args);
    }
}
