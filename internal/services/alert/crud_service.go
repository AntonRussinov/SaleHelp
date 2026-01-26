package alert

import (
	"context"
	"errors"
	"log/slog"
	"salehelp/domain"

	"gorm.io/gorm"
)

// AlertCRUDRepo интерфейс для CRUD операций с алертами
type AlertCRUDRepo interface {
	Create(ctx context.Context, a *domain.Alert) error
	Update(ctx context.Context, a *domain.Alert) error
	Delete(ctx context.Context, id int, shopID int) error
	GetByIDAndShop(ctx context.Context, id int, shopID int) (*domain.Alert, error)
	ListWithFilters(ctx context.Context, shopID int, productID *int, alertType *string, status *string, limit, offset int) ([]*domain.Alert, int, error)
}

// AlertService сервис для CRUD операций с алертами
type AlertService struct {
	repo AlertCRUDRepo
	log  *slog.Logger
}

// NewAlertService создает новый AlertService
func NewAlertService(repo AlertCRUDRepo, log *slog.Logger) *AlertService {
	return &AlertService{
		repo: repo,
		log:  log.With("service", "alert_crud"),
	}
}

// CreateAlert создает новый алерт
func (s *AlertService) CreateAlert(ctx context.Context, shopID int, alert *domain.Alert) (*domain.Alert, error) {
	alert.ShopID = shopID
	alert.Status = "active" // По умолчанию активный

	if err := s.repo.Create(ctx, alert); err != nil {
		s.log.Error("failed to create alert", "error", err, "shopID", shopID)
		return nil, err
	}

	s.log.Info("alert created", "alertID", alert.ID, "shopID", shopID, "type", alert.Type)
	return alert, nil
}

// UpdateAlert обновляет существующий алерт
func (s *AlertService) UpdateAlert(ctx context.Context, alertID int, shopID int, alert *domain.Alert) (*domain.Alert, error) {
	// Проверяем, что алерт существует и принадлежит магазину
	existing, err := s.repo.GetByIDAndShop(ctx, alertID, shopID)
	if err != nil {
		s.log.Error("failed to get alert", "error", err, "alertID", alertID, "shopID", shopID)
		return nil, err
	}

	if existing == nil {
		s.log.Warn("alert not found", "alertID", alertID, "shopID", shopID)
		return nil, gorm.ErrRecordNotFound
	}

	// Обновляем поля
	alert.ID = alertID
	alert.ShopID = shopID

	if err := s.repo.Update(ctx, alert); err != nil {
		s.log.Error("failed to update alert", "error", err, "alertID", alertID, "shopID", shopID)
		return nil, err
	}

	s.log.Info("alert updated", "alertID", alertID, "shopID", shopID)
	return alert, nil
}

// DeleteAlert удаляет алерт
func (s *AlertService) DeleteAlert(ctx context.Context, alertID int, shopID int) error {
	err := s.repo.Delete(ctx, alertID, shopID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.log.Warn("alert not found for deletion", "alertID", alertID, "shopID", shopID)
			return err
		}
		s.log.Error("failed to delete alert", "error", err, "alertID", alertID, "shopID", shopID)
		return err
	}

	s.log.Info("alert deleted", "alertID", alertID, "shopID", shopID)
	return nil
}

// GetAlert получает алерт по ID
func (s *AlertService) GetAlert(ctx context.Context, alertID int, shopID int) (*domain.Alert, error) {
	alert, err := s.repo.GetByIDAndShop(ctx, alertID, shopID)
	if err != nil {
		s.log.Error("failed to get alert", "error", err, "alertID", alertID, "shopID", shopID)
		return nil, err
	}

	if alert == nil {
		s.log.Warn("alert not found", "alertID", alertID, "shopID", shopID)
		return nil, gorm.ErrRecordNotFound
	}

	return alert, nil
}

// ListAlerts получает список алертов с фильтрацией
func (s *AlertService) ListAlerts(ctx context.Context, shopID int, productID *int, alertType *string, status *string, limit, offset int) ([]*domain.Alert, int, error) {
	alerts, total, err := s.repo.ListWithFilters(ctx, shopID, productID, alertType, status, limit, offset)
	if err != nil {
		s.log.Error("failed to list alerts", 
			"error", err, 
			"shopID", shopID,
			"productID", productID,
			"alertType", alertType, 
			"status", status,
			"limit", limit,
			"offset", offset,
		)
		return nil, 0, err
	}

	s.log.Debug("alerts listed", 
		"shopID", shopID,
		"count", len(alerts),
		"total", total,
		"limit", limit,
		"offset", offset,
	)

	return alerts, total, nil
}