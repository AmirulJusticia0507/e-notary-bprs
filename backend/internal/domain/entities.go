package domain

import "time"

type User struct {
	ID                  int64      `db:"id" json:"id"`
	FullName            string     `db:"full_name" json:"full_name"`
	Email               string     `db:"email" json:"email"`
	PasswordHash        string     `db:"password_hash" json:"-"`
	Role                string     `db:"role" json:"role"`
	IsActive            bool       `db:"is_active" json:"is_active"`
	PhotoURL            string     `db:"photo_url" json:"photo_url,omitempty"`
	FailedLoginAttempts int        `db:"failed_login_attempts" json:"-"`
	LockedUntil         *time.Time `db:"locked_until" json:"-"`
	CreatedAt           time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt           *time.Time `db:"deleted_at" json:"-"`
}

type PasswordReset struct {
	ID        int64     `db:"id" json:"id"`
	UserID    int64     `db:"user_id" json:"user_id"`
	Email     string    `db:"email" json:"email,omitempty"`
	TokenHash string    `db:"token_hash" json:"-"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Notary struct {
	ID           int64      `db:"id" json:"id"`
	FullName     string     `db:"full_name" json:"full_name"`
	NotaryNumber string     `db:"notary_number" json:"notary_number"`
	WilayahKerja string     `db:"wilayah_kerja" json:"wilayah_kerja"`
	PhoneNumber  string     `db:"phone_number" json:"phone_number"`
	Email        string     `db:"email" json:"email"`
	IsAvailable  bool       `db:"is_available" json:"is_available"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at" json:"-"`
}

type FinancingApplication struct {
	ID                int64      `db:"id" json:"id"`
	UserID            *int64     `db:"user_id" json:"user_id,omitempty"`
	CustomerName      string     `db:"customer_name" json:"customer_name"`
	CustomerNIK       string     `db:"customer_nik" json:"customer_nik"`
	FinancingAmount   int64      `db:"financing_amount" json:"financing_amount"`
	CollateralType    string     `db:"collateral_type" json:"collateral_type"`
	CollateralDetails string     `db:"collateral_details" json:"collateral_details"`
	Status            string     `db:"status" json:"status"`
	SyncedAt          time.Time  `db:"synced_at" json:"synced_at"`
	CreatedAt         time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt         *time.Time `db:"deleted_at" json:"-"`
}

type Collateral struct {
	ID             int64      `db:"id" json:"id"`
	FinancingID    int64      `db:"financing_id" json:"financing_id"`
	Type           string     `db:"type" json:"type"`
	Details        string     `db:"details" json:"details"`
	EstimatedValue int64      `db:"estimated_value" json:"estimated_value"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt      *time.Time `db:"deleted_at" json:"-"`
}

type LegalOrder struct {
	ID          int64      `db:"id" json:"id"`
	OrderNumber string     `db:"order_number" json:"order_number"`
	FinancingID int64      `db:"financing_id" json:"financing_id"`
	NotaryID    int64      `db:"notary_id" json:"notary_id"`
	AssignedTo  int64      `db:"assigned_to" json:"assigned_to"`
	Status      string     `db:"status" json:"status"`
	SLADeadline time.Time  `db:"sla_deadline" json:"sla_deadline"`
	CompletedAt *time.Time `db:"completed_at" json:"completed_at"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at" json:"-"`
}

type LegalOrderLog struct {
	ID         int64     `db:"id" json:"id"`
	OrderID    int64     `db:"order_id" json:"order_id"`
	PrevStatus string    `db:"prev_status" json:"prev_status"`
	NewStatus  string    `db:"new_status" json:"new_status"`
	ChangedBy  int64     `db:"changed_by" json:"changed_by"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

type LegalDocument struct {
	ID                int64      `db:"id" json:"id"`
	OrderID           int64      `db:"order_id" json:"order_id"`
	FileURL           string     `db:"file_url" json:"file_url"`
	SHA256Hash        string     `db:"sha256_hash" json:"sha256_hash"`
	EMeteraiSN        string     `db:"e_meterai_sn" json:"e_meterai_sn"`
	ESignStatus       string     `db:"e_sign_status" json:"e_sign_status"`
	ActNumber         string     `db:"act_number" json:"act_number"`
	MinutesStatus     string     `db:"minutes_status" json:"minutes_status"`
	NotaryFee         int64      `db:"notary_fee" json:"notary_fee"`
	ProcessingStatus  string     `db:"processing_status" json:"processing_status"`
	NotaryProcessedAt *time.Time `db:"notary_processed_at" json:"notary_processed_at"`
	CreatedAt         time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt         *time.Time `db:"deleted_at" json:"-"`
}

type OrderDetail struct {
	ID           int64      `db:"id" json:"id"`
	OrderNumber  string     `db:"order_number" json:"order_number"`
	FinancingID  int64      `db:"financing_id" json:"financing_id"`
	NotaryID     int64      `db:"notary_id" json:"notary_id"`
	AssignedTo   int64      `db:"assigned_to" json:"assigned_to"`
	Status       string     `db:"status" json:"status"`
	SLADeadline  time.Time  `db:"sla_deadline" json:"sla_deadline"`
	CompletedAt  *time.Time `db:"completed_at" json:"completed_at"`
	CustomerName string     `db:"customer_name" json:"customer_name"`
	NotaryName   string     `db:"notary_name" json:"notary_name"`
	AssignedName string     `db:"assigned_name" json:"assigned_name"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
}
