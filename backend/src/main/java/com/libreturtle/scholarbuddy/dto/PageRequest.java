package com.libreturtle.scholarbuddy.dto;

import com.libreturtle.scholarbuddy.exception.ApiException;

public record PageRequest(int page, int perPage) {

    public PageRequest {
        if (page < 1) {
            throw ApiException.badRequest("invalid_page", "page must be greater than 0");
        }
        if (perPage < 1 || perPage > 100) {
            throw ApiException.badRequest("invalid_per_page", "per_page must be between 1 and 100");
        }
    }
}
