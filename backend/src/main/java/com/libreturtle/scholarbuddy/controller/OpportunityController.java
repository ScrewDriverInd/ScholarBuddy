package com.libreturtle.scholarbuddy.controller;

import com.libreturtle.scholarbuddy.dto.ApiResponseBody;
import com.libreturtle.scholarbuddy.dto.OpportunityRequest;
import com.libreturtle.scholarbuddy.dto.OpportunityResponse;
import com.libreturtle.scholarbuddy.dto.PageRequest;
import com.libreturtle.scholarbuddy.dto.PageResponse;
import com.libreturtle.scholarbuddy.model.OpportunityType;
import com.libreturtle.scholarbuddy.model.User;
import com.libreturtle.scholarbuddy.security.SecurityUtils;
import com.libreturtle.scholarbuddy.service.OpportunityService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.UUID;

@RestController
@RequestMapping("/api/v1/opportunities")
@RequiredArgsConstructor
public class OpportunityController {

    private final OpportunityService opportunityService;

    @GetMapping
    public ResponseEntity<ApiResponseBody<PageResponse<OpportunityResponse>>> list(
            @RequestParam(name = "type", required = false) OpportunityType type,
            @RequestParam(name = "page", defaultValue = "1") int page,
            @RequestParam(name = "per_page", defaultValue = "20") int perPage) {
        PageRequest pageRequest = new PageRequest(page, perPage);
        return ResponseEntity.ok(ApiResponseBody.of(opportunityService.listApproved(type, pageRequest.page(), pageRequest.perPage())));
    }

    @GetMapping("/{id}")
    public ResponseEntity<ApiResponseBody<OpportunityResponse>> get(@PathVariable UUID id) {
        return ResponseEntity.ok(ApiResponseBody.of(opportunityService.getAndRecordClick(id)));
    }

    @PostMapping
    public ResponseEntity<ApiResponseBody<OpportunityResponse>> create(@Valid @RequestBody OpportunityRequest request) {
        User user = SecurityUtils.getCurrentUser();
        return ResponseEntity.status(HttpStatus.CREATED)
                .body(ApiResponseBody.of(opportunityService.create(request, user)));
    }

    @PatchMapping("/{id}")
    public ResponseEntity<ApiResponseBody<OpportunityResponse>> update(
            @PathVariable UUID id,
            @Valid @RequestBody OpportunityRequest request) {
        User user = SecurityUtils.getCurrentUser();
        return ResponseEntity.ok(ApiResponseBody.of(opportunityService.update(id, request, user)));
    }
}
