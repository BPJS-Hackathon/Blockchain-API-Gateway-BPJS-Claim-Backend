package models

const (
	ClaimStatusSubmitted = "SUBMITTED" // Status awal, dilakukan bersamaan saat RS / faskes yang merawat mensubmit Rekam Medis
	ClaimStatusPending   = "PENDING"   // Claim masuk ke blockchain, smart contract memverifikasi apakah claim fiktif / ngga. Node BPJS trigger update db jika menemukan ada tx claim di block.
	ClaimStatusPaid      = "PAID"      // Admin BPJS update status claim setelah melakukan pembayaran
	ClaimStatusRejected  = "REJECTED"  // NOTE: BISA DILAKUKAN BLOCKCHAIN (Jika smart contract mark invalid) ATAU BPJS (verifikasi akhir)
	ClaimStatusFaked     = "FAKED"     // Jika smart contract menandai claim sebagai fiktif
)

type Claims struct {
	ClaimID      string `gorm:"primaryKey;type:uuid" json:"claim_id"`
	RekamMedisID string `gorm:"type:uuid;not null" json:"rekam_medis_id"`
	Amount       uint   `gorm:"not null" json:"amount"`
	Status       string `gorm:"type:varchar(20);not null" json:"status"`
}
