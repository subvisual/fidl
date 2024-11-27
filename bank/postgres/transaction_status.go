package postgres

type TransactionStatus int8

const (
	TransactionPending TransactionStatus = iota + 1
	TransactionCompleted
)

func (a TransactionStatus) String() string {
	switch a {
	case TransactionPending:
		return "Pending"
	case TransactionCompleted:
		return "Completed"
	default:
		return "Unknown" // nolint:goconst
	}
}
