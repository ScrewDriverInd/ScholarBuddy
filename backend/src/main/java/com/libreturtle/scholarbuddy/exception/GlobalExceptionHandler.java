package com.libreturtle.scholarbuddy.exception;

import com.libreturtle.scholarbuddy.dto.ErrorResponse;
import jakarta.servlet.http.HttpServletRequest;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.http.converter.HttpMessageNotReadableException;
import org.springframework.validation.FieldError;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;
import org.springframework.web.method.annotation.MethodArgumentTypeMismatchException;
import org.springframework.web.HttpRequestMethodNotSupportedException;
import org.springframework.web.servlet.resource.NoResourceFoundException;

import java.util.Map;
import java.util.stream.Collectors;

@RestControllerAdvice
@Slf4j
public class GlobalExceptionHandler {

    @ExceptionHandler(ApiException.class)
    public ResponseEntity<Map<String, ErrorResponse>> handleApiException(ApiException ex, HttpServletRequest request) {
        String requestId = request.getHeader("X-Request-ID");
        ErrorResponse error = new ErrorResponse(ex.getCode(), ex.getMessage(), requestId);
        return ResponseEntity.status(ex.getStatus()).body(Map.of("error", error));
    }

    @ExceptionHandler(MethodArgumentNotValidException.class)
    public ResponseEntity<Map<String, ErrorResponse>> handleValidationException(
            MethodArgumentNotValidException ex, HttpServletRequest request) {
        String requestId = request.getHeader("X-Request-ID");
        String message = ex.getBindingResult().getFieldErrors().stream()
                .map(FieldError::getDefaultMessage)
                .collect(Collectors.joining(", "));
        ErrorResponse error = new ErrorResponse("validation_error", message, requestId);
        return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(Map.of("error", error));
    }

    @ExceptionHandler(MethodArgumentTypeMismatchException.class)
    public ResponseEntity<Map<String, ErrorResponse>> handleTypeMismatch(
            MethodArgumentTypeMismatchException ex, HttpServletRequest request) {
        String requestId = request.getHeader("X-Request-ID");
        String message = "invalid parameter";
        if (ex.getName().equals("id")) {
            message = "ID must be a valid UUID";
        }
        ErrorResponse error = new ErrorResponse("invalid_parameter", message, requestId);
        return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(Map.of("error", error));
    }

    @ExceptionHandler(HttpRequestMethodNotSupportedException.class)
    public ResponseEntity<Map<String, ErrorResponse>> handleMethodNotSupported(
            HttpRequestMethodNotSupportedException ex, HttpServletRequest request) {
        String requestId = request.getHeader("X-Request-ID");
        ErrorResponse error = new ErrorResponse(
                "method_not_allowed", "use PATCH to approve a listing", requestId);
        return ResponseEntity.status(HttpStatus.METHOD_NOT_ALLOWED).body(Map.of("error", error));
    }

    @ExceptionHandler(HttpMessageNotReadableException.class)
    public ResponseEntity<Map<String, ErrorResponse>> handleUnreadableMessage(
            HttpMessageNotReadableException ex, HttpServletRequest request) {
        String requestId = request.getHeader("X-Request-ID");
        ErrorResponse error = new ErrorResponse(
                "invalid_request", "request body must be valid JSON", requestId);
        return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(Map.of("error", error));
    }

    @ExceptionHandler(NoResourceFoundException.class)
    public ResponseEntity<Map<String, ErrorResponse>> handleNoResourceFound(
            NoResourceFoundException ex, HttpServletRequest request) {
        String requestId = request.getHeader("X-Request-ID");
        ErrorResponse error = new ErrorResponse("not_found", "Resource not found", requestId);
        return ResponseEntity.status(HttpStatus.NOT_FOUND).body(Map.of("error", error));
    }

    @ExceptionHandler(Exception.class)
    public ResponseEntity<Map<String, ErrorResponse>> handleGenericException(Exception ex, HttpServletRequest request) {
        String requestId = request.getHeader("X-Request-ID");
        log.error("Unexpected error occurred", ex);
        ErrorResponse error = new ErrorResponse("internal_error", "An unexpected error occurred", requestId);
        return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(Map.of("error", error));
    }
}
