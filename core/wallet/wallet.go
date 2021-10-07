package main

import "encoding/binary"

func CreditOrDebit(found bool, hash Hash, b *Bucket, item int64, param int64) OperationResult {
	if found {
		acc := b.ReadItem(item)
		balance := int64(binary.LittleEndian.Uint64(acc[size:]))
		newbalance := balance + param
		if newbalance > 0 {
			return OperationResult{
				result: QueryResult{ok: true, data: acc},
			}
		} else if newbalance == 0 {
			return OperationResult{
				deleted: &Item{bucket: b, item: item},
				result:  QueryResult{ok: true, data: acc},
			}
		} else {
			return OperationResult{
				result: QueryResult{ok: false},
			}

		}
	} else {
		if param >= 0 {
			acc := make([]byte, size+8)
			binary.LittleEndian.PutUint64(acc[size:], uint64(param))
			copy(acc[0:size], hash[:])
			b.WriteItem(item, acc)
			return OperationResult{
				added:  &Item{bucket: b, item: item},
				result: QueryResult{ok: false, data: acc},
			}
		} else {
			return OperationResult{
				result: QueryResult{
					ok: false,
				},
			}
		}
	}
}

func NewMemoryWalletStore()
