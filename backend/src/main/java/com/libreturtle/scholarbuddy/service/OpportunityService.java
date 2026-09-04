package com.libreturtle.scholarbuddy.service;

import com.libreturtle.scholarbuddy.dto.OpportunityRequest;
import com.libreturtle.scholarbuddy.dto.OpportunityResponse;
import com.libreturtle.scholarbuddy.dto.PageResponse;
import com.libreturtle.scholarbuddy.exception.ApiException;
import com.libreturtle.scholarbuddy.model.ApprovalStatus;
import com.libreturtle.scholarbuddy.model.Opportunity;
import com.libreturtle.scholarbuddy.model.OpportunityType;
import com.libreturtle.scholarbuddy.model.User;
import com.libreturtle.scholarbuddy.repository.OpportunityRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

@Service
@RequiredArgsConstructor
public class OpportunityService {

    private final OpportunityRepository opportunityRepository;

    @Transactional(readOnly = true)
    public PageResponse<OpportunityResponse> listApproved(OpportunityType type, int page, int perPage) {
        Pageable pageable = PageRequest.of(page - 1, perPage);
        Page<Opportunity> opportunities;

        if (type == null) {
            opportunities = opportunityRepository.findByApprovalStatus(ApprovalStatus.APPROVED, pageable);
        } else {
            opportunities = opportunityRepository.findByApprovalStatusAndType(ApprovalStatus.APPROVED, type, pageable);
        }

        return new PageResponse<>(
            opportunities.map(this::mapToResponse).getContent(),
            page,
            perPage,
            opportunities.getTotalElements()
        );
    }

    @Transactional(readOnly = true)
    public PageResponse<OpportunityResponse> listPending(int page, int perPage) {
        Pageable pageable = PageRequest.of(page - 1, perPage);
        Page<Opportunity> opportunities = opportunityRepository.findByApprovalStatus(ApprovalStatus.PENDING, pageable);

        return new PageResponse<>(
            opportunities.map(this::mapToResponse).getContent(),
            page,
            perPage,
            opportunities.getTotalElements()
        );
    }

    @Transactional
    public OpportunityResponse getAndRecordClick(UUID id) {
        int updated = opportunityRepository.incrementClickCount(id);
        if (updated == 0) {
            throw ApiException.notFound("opportunity was not found");
        }
        Opportunity opportunity = opportunityRepository.findApprovedById(id)
                .orElseThrow(() -> ApiException.notFound("opportunity was not found"));
        return mapToResponse(opportunity);
    }

    @Transactional
    public OpportunityResponse create(OpportunityRequest request, User user) {
        Opportunity opportunity = new Opportunity();
        opportunity.setTitle(request.title());
        opportunity.setDescription(request.description());
        opportunity.setTypes(request.types());
        opportunity.setEligibility(request.eligibility());
        opportunity.setSteps(request.steps());
        opportunity.setBenefits(request.benefits());
        opportunity.setLink(request.link());
        opportunity.setReferral(request.referral());
        opportunity.setCreatedBy(user);
        opportunity.setApprovalStatus(ApprovalStatus.PENDING);
        opportunity.setClickCount(0L);

        opportunity = opportunityRepository.save(opportunity);
        return mapToResponse(opportunity);
    }

    @Transactional
    public OpportunityResponse update(UUID id, OpportunityRequest request, User user) {
        Opportunity opportunity = opportunityRepository.findByIdAndCreatedById(id, user.getId())
                .orElseThrow(() -> ApiException.forbidden("you can only update opportunities you created"));

        opportunity.setTitle(request.title());
        opportunity.setDescription(request.description());
        opportunity.setTypes(request.types());
        opportunity.setEligibility(request.eligibility());
        opportunity.setSteps(request.steps());
        opportunity.setBenefits(request.benefits());
        opportunity.setLink(request.link());
        opportunity.setReferral(request.referral());
        opportunity.setApprovalStatus(ApprovalStatus.PENDING);

        opportunity = opportunityRepository.save(opportunity);
        return mapToResponse(opportunity);
    }

    @Transactional
    public OpportunityResponse approve(UUID id) {
        Opportunity opportunity = opportunityRepository.findById(id)
                .filter(o -> o.getApprovalStatus() == ApprovalStatus.PENDING)
                .orElseThrow(() -> ApiException.notFound("pending opportunity was not found"));

        opportunity.setApprovalStatus(ApprovalStatus.APPROVED);
        opportunity = opportunityRepository.save(opportunity);
        return mapToResponse(opportunity);
    }

    @Transactional
    public void delete(UUID id) {
        if (!opportunityRepository.existsById(id)) {
            throw ApiException.notFound("opportunity was not found");
        }
        opportunityRepository.deleteById(id);
    }

    private OpportunityResponse mapToResponse(Opportunity opportunity) {
        return new OpportunityResponse(
                opportunity.getId(),
                opportunity.getTitle(),
                opportunity.getDescription(),
                opportunity.getTypes(),
                opportunity.getEligibility(),
                opportunity.getSteps(),
                opportunity.getBenefits(),
                opportunity.getLink(),
                opportunity.getReferral(),
                opportunity.getCreatedBy().getId(),
                opportunity.getApprovalStatus(),
                opportunity.getClickCount(),
                opportunity.getCreatedAt(),
                opportunity.getUpdatedAt()
        );
    }
}
