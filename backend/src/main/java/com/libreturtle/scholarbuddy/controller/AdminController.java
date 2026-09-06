package com.libreturtle.scholarbuddy.controller;

import com.libreturtle.scholarbuddy.dto.ApiResponseBody;
import com.libreturtle.scholarbuddy.dto.ListingResponse;
import com.libreturtle.scholarbuddy.dto.PageRequest;
import com.libreturtle.scholarbuddy.dto.PageResponse;
import com.libreturtle.scholarbuddy.dto.UserResponse;
import com.libreturtle.scholarbuddy.service.ListingService;
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

    private final ListingService listingService;
    private final UserService userService;

    @GetMapping("/listings")
    public ResponseEntity<ApiResponseBody<PageResponse<ListingResponse>>> listPending(
            @RequestParam(name = "page", defaultValue = "1") int page,
            @RequestParam(name = "per_page", defaultValue = "20") int perPage) {
        PageRequest pageRequest = new PageRequest(page, perPage);
        return ResponseEntity.ok(ApiResponseBody.of(listingService.listPending(pageRequest.page(), pageRequest.perPage())));
    }

    @PatchMapping("/listings/{id}/approve")
    public ResponseEntity<ApiResponseBody<ListingResponse>> approve(@PathVariable UUID id) {
        return ResponseEntity.ok(ApiResponseBody.of(listingService.approve(id)));
    }

    @DeleteMapping("/listings/{id}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    public void delete(@PathVariable UUID id) {
        listingService.delete(id);
    }

    @PostMapping("/users/{id}/admin")
    public ResponseEntity<ApiResponseBody<UserResponse>> grantAdmin(@PathVariable UUID id) {
        return ResponseEntity.ok(ApiResponseBody.of(userService.grantAdminRole(id)));
    }
}
