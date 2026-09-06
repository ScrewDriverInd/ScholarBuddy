package com.libreturtle.scholarbuddy.controller;

import com.libreturtle.scholarbuddy.dto.ApiResponseBody;
import com.libreturtle.scholarbuddy.dto.ListingRequest;
import com.libreturtle.scholarbuddy.dto.ListingResponse;
import com.libreturtle.scholarbuddy.dto.PageRequest;
import com.libreturtle.scholarbuddy.dto.PageResponse;
import com.libreturtle.scholarbuddy.model.ListingType;
import com.libreturtle.scholarbuddy.model.User;
import com.libreturtle.scholarbuddy.security.SecurityUtils;
import com.libreturtle.scholarbuddy.service.ListingService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.UUID;

@RestController
@RequestMapping("/api/v1/listings")
@RequiredArgsConstructor
public class ListingController {

    private final ListingService listingService;

    @GetMapping
    public ResponseEntity<ApiResponseBody<PageResponse<ListingResponse>>> list(
            @RequestParam(name = "type", required = false) ListingType type,
            @RequestParam(name = "page", defaultValue = "1") int page,
            @RequestParam(name = "per_page", defaultValue = "20") int perPage) {
        PageRequest pageRequest = new PageRequest(page, perPage);
        return ResponseEntity.ok(ApiResponseBody.of(listingService.listApproved(type, pageRequest.page(), pageRequest.perPage())));
    }

    @GetMapping("/{id}")
    public ResponseEntity<ApiResponseBody<ListingResponse>> get(@PathVariable UUID id) {
        return ResponseEntity.ok(ApiResponseBody.of(listingService.getAndRecordClick(id)));
    }

    @PostMapping
    public ResponseEntity<ApiResponseBody<ListingResponse>> create(@Valid @RequestBody ListingRequest request) {
        User user = SecurityUtils.getCurrentUser();
        return ResponseEntity.status(HttpStatus.CREATED)
                .body(ApiResponseBody.of(listingService.create(request, user)));
    }

    @PatchMapping("/{id}")
    public ResponseEntity<ApiResponseBody<ListingResponse>> update(
            @PathVariable UUID id,
            @Valid @RequestBody ListingRequest request) {
        User user = SecurityUtils.getCurrentUser();
        return ResponseEntity.ok(ApiResponseBody.of(listingService.update(id, request, user)));
    }
}
