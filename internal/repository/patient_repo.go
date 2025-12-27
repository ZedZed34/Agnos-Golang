package repository

import (
	"agnos-gin/internal/domain"
	"gorm.io/gorm"
)

type patientRepo struct {
	db *gorm.DB
}

func NewPatientRepository(db *gorm.DB) domain.PatientRepository {
	return &patientRepo{db: db}
}

func (r *patientRepo) Create(p *domain.Patient) error {
	return r.db.Create(p).Error
}

func (r *patientRepo) Search(filters map[string]interface{}, staffHospital string) ([]domain.Patient, error) {
	var patients []domain.Patient
	query := r.db.Model(&domain.Patient{})

	// Requirement #3: Enforce hospital isolation
	query = query.Where("hospital_name = ?", staffHospital)

	// Requirement #4: Dynamic Filters
	if val, ok := filters["national_id"]; ok && val != "" {
		query = query.Where("national_id = ?", val)
	}
	if val, ok := filters["passport_id"]; ok && val != "" {
		query = query.Where("passport_id = ?", val)
	}
	if val, ok := filters["first_name"]; ok && val != "" {
		query = query.Where("first_name_en ILIKE ? OR first_name_th ILIKE ?", "%"+val.(string)+"%", "%"+val.(string)+"%")
	}

	err := query.Find(&patients).Error
	return patients, err
}