package com.libreturtle.scholarbuddy.dto;

import com.libreturtle.scholarbuddy.model.ApprovalStatus;
import com.libreturtle.scholarbuddy.model.ListingType;

import java.time.Instant;
import java.util.Set;
import java.util.UUID;

public record ListingResponse(
        UUID id,
        String title,
        String description,
        Set<ListingType> types,
        String eligibility,
        String steps,
        String benefits,
        String link,
        String referral,
        UUID createdBy,
        ApprovalStatus approvalStatus,
        Long clickCount,
        Instant createdAt,
        Instant updatedAt
) {
}
