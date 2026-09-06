package com.libreturtle.scholarbuddy.exception;

import lombok.AllArgsConstructor;
import lombok.Data;

@Data
@AllArgsConstructor
public class ErrorResponse {
    private Error error;

    @Data
    @AllArgsConstructor
    public static class Error {
        private String code;
        private String message;
        private String requestId;
    }
}
