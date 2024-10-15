package postgres

type BankService struct {
	db *DB
}

func NewBankService(db *DB) *BankService {
	return &BankService{
		db: db,
	}
}
