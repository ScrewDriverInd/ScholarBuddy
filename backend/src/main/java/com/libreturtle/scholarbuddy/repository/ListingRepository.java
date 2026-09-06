package com.libreturtle.scholarbuddy.repository;

import com.libreturtle.scholarbuddy.model.ApprovalStatus;
import com.libreturtle.scholarbuddy.model.Listing;
import com.libreturtle.scholarbuddy.model.ListingType;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.util.Optional;
import java.util.UUID;

@Repository
public interface ListingRepository extends JpaRepository<Listing, UUID> {

    @Query("SELECT o FROM Listing o WHERE o.approvalStatus = :status ORDER BY o.clickCount DESC, o.createdAt DESC")
    Page<Listing> findByApprovalStatus(@Param("status") ApprovalStatus status, Pageable pageable);

    @Query("SELECT o FROM Listing o JOIN o.types t WHERE o.approvalStatus = :status AND t = :type ORDER BY o.clickCount DESC, o.createdAt DESC")
    Page<Listing> findByApprovalStatusAndType(@Param("status") ApprovalStatus status, @Param("type") ListingType type, Pageable pageable);

    @Modifying(clearAutomatically = true)
    @Query("UPDATE Listing o SET o.clickCount = o.clickCount + 1 WHERE o.id = :id AND o.approvalStatus = 'APPROVED'")
    int incrementClickCount(@Param("id") UUID id);

    @Query("SELECT o FROM Listing o WHERE o.id = :id AND o.approvalStatus = 'APPROVED'")
    Optional<Listing> findApprovedById(@Param("id") UUID id);

    Optional<Listing> findByIdAndCreatedById(UUID id, UUID userId);
}
