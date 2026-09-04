package com.libreturtle.scholarbuddy.model;

import jakarta.persistence.*;
import lombok.Getter;
import lombok.Setter;
import org.hibernate.annotations.CreationTimestamp;
import org.hibernate.annotations.UpdateTimestamp;

import java.time.Instant;
import java.util.HashSet;
import java.util.Set;
import java.util.UUID;

@Entity
@Table(name = "opportunities")
@Getter
@Setter
public class Opportunity {

    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private UUID id;

    @Column(nullable = false, length = 250)
    private String title;

    @Column(nullable = false, length = 10000)
    private String description;

    @ElementCollection(fetch = FetchType.EAGER)
    @CollectionTable(name = "opportunity_types", joinColumns = @JoinColumn(name = "opportunity_id"))
    @Column(name = "type")
    @Enumerated(EnumType.STRING)
    private Set<OpportunityType> types = new HashSet<>();

    @Column(nullable = false, columnDefinition = "TEXT")
    private String eligibility = "";

    @Column(nullable = false, columnDefinition = "TEXT")
    private String steps = "";

    @Column(nullable = false, columnDefinition = "TEXT")
    private String benefits = "";

    @Column(nullable = false, columnDefinition = "TEXT")
    private String link = "";

    @Column(nullable = false, columnDefinition = "TEXT")
    private String referral = "";

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "created_by", nullable = false)
    private User createdBy;

    @Enumerated(EnumType.STRING)
    @Column(name = "approval_status", nullable = false)
    private ApprovalStatus approvalStatus = ApprovalStatus.PENDING;

    @Column(name = "click_count", nullable = false)
    private Long clickCount = 0L;

    @CreationTimestamp
    @Column(name = "created_at", nullable = false, updatable = false)
    private Instant createdAt;

    @UpdateTimestamp
    @Column(name = "updated_at", nullable = false)
    private Instant updatedAt;
}
