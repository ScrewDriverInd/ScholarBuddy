package com.libreturtle.scholarbuddy.controller;

import com.libreturtle.scholarbuddy.dto.ApiResponseBody;
import com.libreturtle.scholarbuddy.dto.OpportunityResponse;
import com.libreturtle.scholarbuddy.dto.PageRequest;
import com.libreturtle.scholarbuddy.dto.PageResponse;
import com.libreturtle.scholarbuddy.dto.UserResponse;
import com.libreturtle.scholarbuddy.service.OpportunityService;
import com.libreturtle.scholarbuddy.service.UserService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.UUID;

@RestController
@RequestMapping("/api/v1/abbujaan")
@RequiredArgsConstructor
public class AdminController {

    private final OpportunityService opportunityService;
    private final UserService userService;

    @GetMapping("/opportunities")
    public ResponseEntity<ApiResponseBody<PageResponse<OpportunityResponse>>> listPending(
            @RequestParam(name = "page", defaultValue = "1") int page,
            @RequestParam(name = "per_page", defaultValue = "20") int perPage) {
        PageRequest pageRequest = new PageRequest(page, perPage);
        return ResponseEntity.ok(ApiResponseBody.of(opportunityService.listPending(pageRequest.page(), pageRequest.perPage())));
    }

    @PatchMapping("/opportunities/{id}/approve")
    public ResponseEntity<ApiResponseBody<OpportunityResponse>> approve(@PathVariable UUID id) {
        return ResponseEntity.ok(ApiResponseBody.of(opportunityService.approve(id)));
    }

    @DeleteMapping("/opportunities/{id}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    public void delete(@PathVariable UUID id) {
        opportunityService.delete(id);
    }

    @PostMapping("/users/{id}/admin")
    public ResponseEntity<ApiResponseBody<UserResponse>> grantAdmin(@PathVariable UUID id) {
        return ResponseEntity.ok(ApiResponseBody.of(userService.grantAdminRole(id)));
    }
}
