package com.libreturtle.scholarbuddy.dto;

public record ErrorResponse(String code, String message, String requestId) {
}
