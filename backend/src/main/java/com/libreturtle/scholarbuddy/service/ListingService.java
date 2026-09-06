package com.libreturtle.scholarbuddy.service;

import com.libreturtle.scholarbuddy.dto.ListingRequest;
import com.libreturtle.scholarbuddy.dto.ListingResponse;
import com.libreturtle.scholarbuddy.dto.PageResponse;
import com.libreturtle.scholarbuddy.exception.ApiException;
import com.libreturtle.scholarbuddy.model.ApprovalStatus;
import com.libreturtle.scholarbuddy.model.Listing;
import com.libreturtle.scholarbuddy.model.ListingType;
import com.libreturtle.scholarbuddy.model.User;
import com.libreturtle.scholarbuddy.repository.ListingRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

@Service
@RequiredArgsConstructor
public class ListingService {

    private final ListingRepository listingRepository;

    @Transactional(readOnly = true)
    public PageResponse<ListingResponse> listApproved(ListingType type, int page, int perPage) {
        Pageable pageable = PageRequest.of(page - 1, perPage);
        Page<Listing> listings;

        if (type == null) {
            listings = listingRepository.findByApprovalStatus(ApprovalStatus.APPROVED, pageable);
        } else {
            listings = listingRepository.findByApprovalStatusAndType(ApprovalStatus.APPROVED, type, pageable);
        }

        return new PageResponse<>(
            listings.map(this::mapToResponse).getContent(),
            page,
            perPage,
            listings.getTotalElements()
        );
    }

    @Transactional(readOnly = true)
    public PageResponse<ListingResponse> listPending(int page, int perPage) {
        Pageable pageable = PageRequest.of(page - 1, perPage);
        Page<Listing> listings = listingRepository.findByApprovalStatus(ApprovalStatus.PENDING, pageable);

        return new PageResponse<>(
            listings.map(this::mapToResponse).getContent(),
            page,
            perPage,
            listings.getTotalElements()
        );
    }

    @Transactional
    public ListingResponse getAndRecordClick(UUID id) {
        int updated = listingRepository.incrementClickCount(id);
        if (updated == 0) {
            throw ApiException.notFound("listing was not found");
        }
        Listing listing = listingRepository.findApprovedById(id)
                .orElseThrow(() -> ApiException.notFound("listing was not found"));
        return mapToResponse(listing);
    }

    @Transactional
    public ListingResponse create(ListingRequest request, User user) {
        Listing listing = new Listing();
        listing.setTitle(request.title());
        listing.setDescription(request.description());
        listing.setTypes(request.types());
        listing.setEligibility(request.eligibility());
        listing.setSteps(request.steps());
        listing.setBenefits(request.benefits());
        listing.setLink(request.link());
        listing.setReferral(request.referral());
        listing.setCreatedBy(user);
        listing.setApprovalStatus(ApprovalStatus.PENDING);
        listing.setClickCount(0L);

        listing = listingRepository.save(listing);
        return mapToResponse(listing);
    }

    @Transactional
    public ListingResponse update(UUID id, ListingRequest request, User user) {
        Listing listing = listingRepository.findByIdAndCreatedById(id, user.getId())
                .orElseThrow(() -> ApiException.forbidden("you can only update listings you created"));

        listing.setTitle(request.title());
        listing.setDescription(request.description());
        listing.setTypes(request.types());
        listing.setEligibility(request.eligibility());
        listing.setSteps(request.steps());
        listing.setBenefits(request.benefits());
        listing.setLink(request.link());
        listing.setReferral(request.referral());
        listing.setApprovalStatus(ApprovalStatus.PENDING);

        listing = listingRepository.save(listing);
        return mapToResponse(listing);
    }

    @Transactional
    public ListingResponse approve(UUID id) {
        Listing listing = listingRepository.findById(id)
                .filter(o -> o.getApprovalStatus() == ApprovalStatus.PENDING)
                .orElseThrow(() -> ApiException.notFound("pending listing was not found"));

        listing.setApprovalStatus(ApprovalStatus.APPROVED);
        listing = listingRepository.save(listing);
        return mapToResponse(listing);
    }

    @Transactional
    public void delete(UUID id) {
        if (!listingRepository.existsById(id)) {
            throw ApiException.notFound("listing was not found");
        }
        listingRepository.deleteById(id);
    }

    private ListingResponse mapToResponse(Listing listing) {
        return new ListingResponse(
                listing.getId(),
                listing.getTitle(),
                listing.getDescription(),
                listing.getTypes(),
                listing.getEligibility(),
                listing.getSteps(),
                listing.getBenefits(),
                listing.getLink(),
                listing.getReferral(),
                listing.getCreatedBy().getId(),
                listing.getApprovalStatus(),
                listing.getClickCount(),
                listing.getCreatedAt(),
                listing.getUpdatedAt()
        );
    }
}
