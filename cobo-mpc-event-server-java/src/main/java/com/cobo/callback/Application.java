package com.cobo.callback;

import java.util.HashMap;
import java.util.Map;

import com.cobo.callback.config.AppConfig;
import com.cobo.callback.model.Request;
import com.cobo.callback.model.Response;
import com.cobo.callback.service.JwtService;
import com.cobo.callback.verify.TssVerifier;
import com.cobo.callback.verify.Verifier;
import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.PropertyNamingStrategies;
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
    private static final ObjectMapper MAPPER = new ObjectMapper()
            .setPropertyNamingStrategy(PropertyNamingStrategies.SNAKE_CASE)
            .configure(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES, false)
            .configure(DeserializationFeature.FAIL_ON_NULL_FOR_PRIMITIVES, false)
            .configure(DeserializationFeature.READ_UNKNOWN_ENUM_VALUES_AS_NULL, true)
            .configure(DeserializationFeature.ACCEPT_EMPTY_STRING_AS_NULL_OBJECT, true);

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
            Request request = MAPPER.readValue(requestData, Request.class);

            Response response = processRequest(request);
            String responseJson = MAPPER.writeValueAsString(response);
            return jwtService.createToken(responseJson);

        } catch (JwtException e) {
            return handleError(Response.STATUS_INVALID_TOKEN, e.getMessage());
        } catch (Exception e) {
            log.error("Failed to process request", e);
            return handleError(Response.STATUS_INTERNAL_ERROR, e.getMessage());
        }
    }

    private String handleError(int status, String message) {
        try {
            Response response = Response.builder()
                    .status(status)
                    .errStr(message)
                    .build();
            return jwtService.createToken(MAPPER.writeValueAsString(response));
        } catch (Exception ex) {
            log.error("Error handling error response", ex);
            throw new RuntimeException(ex);
        }
    }

    private static Response processRequest(Request request) {
        String error = verifier.verify(request);
        if (error != null) {
            return Response.builder()
                    .status(Response.STATUS_INVALID_REQUEST)
                    .requestId(request.getRequestId())
                    .action(Response.ACTION_REJECT)
                    .errStr(error)
                    .build();
        }

        return Response.builder()
                .status(Response.STATUS_OK)
                .requestId(request.getRequestId())
                .action(Response.ACTION_APPROVE)
                .build();
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
