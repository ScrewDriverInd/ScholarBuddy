package com.libreturtle.scholarbuddy.dto;

public record ApiResponseBody<T>(T data) {
    public static <T> ApiResponseBody<T> of(T data) {
        return new ApiResponseBody<>(data);
    }
}
