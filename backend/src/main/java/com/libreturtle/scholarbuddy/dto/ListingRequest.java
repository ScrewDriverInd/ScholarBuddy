package com.libreturtle.scholarbuddy.dto;

import com.libreturtle.scholarbuddy.model.ListingType;
import com.libreturtle.scholarbuddy.validation.ValidUrl;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotEmpty;
import jakarta.validation.constraints.Size;

import java.util.Set;

public record ListingRequest(
        @NotBlank(message = "title is required")
        @Size(min = 1, max = 250, message = "title must be between 1 and 250 characters")
        String title,

        @NotBlank(message = "description is required")
        @Size(min = 1, max = 10000, message = "description must be between 1 and 10000 characters")
        String description,

        @NotEmpty(message = "types must contain at least one listing type")
        @Size(max = 5, message = "types cannot contain more than 5 listing types")
        Set<ListingType> types,

        String eligibility,

        String steps,

        String benefits,

        @ValidUrl
        String link,

        String referral
) {
    public ListingRequest {
        if (eligibility == null) eligibility = "";
        if (steps == null) steps = "";
        if (benefits == null) benefits = "";
        if (link == null) link = "";
        if (referral == null) referral = "";
    }
}
