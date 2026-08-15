package services

import (
	"errors"
	"money-tracker-api/internal/models"

	"gorm.io/gorm"
)

// CreateTransaction menangani pembuatan transaksi dan memperbarui saldo wallet.
func CreateTransaction(tx *gorm.DB, transaction *models.Transaction) error {
	// 1. Cari wallet yang dikaitkan dan pastikan milik user terkait
	var wallet models.Wallet
	if err := tx.First(&wallet, "id = ? AND user_id = ?", transaction.WalletID, transaction.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("wallet tidak ditemukan")
		}
		return err
	}

	// 2. Sesuaikan saldo wallet
	if transaction.Type == "income" {
		wallet.Balance += transaction.Amount
	} else if transaction.Type == "expense" {
		wallet.Balance -= transaction.Amount
	} else {
		return errors.New("tipe transaksi tidak valid")
	}

	// 3. Simpan perubahan saldo wallet
	if err := tx.Save(&wallet).Error; err != nil {
		return err
	}

	// 4. Buat data transaksi
	if err := tx.Create(transaction).Error; err != nil {
		return err
	}

	return nil
}

// UpdateTransaction menangani perubahan transaksi dan menyesuaikan saldo wallet.
func UpdateTransaction(tx *gorm.DB, oldTx *models.Transaction, newTxInput *models.Transaction) error {
	// 1. Kembalikan saldo wallet lama (revert efek transaksi lama)
	var oldWallet models.Wallet
	if err := tx.First(&oldWallet, "id = ? AND user_id = ?", oldTx.WalletID, oldTx.UserID).Error; err != nil {
		return err
	}

	if oldTx.Type == "income" {
		oldWallet.Balance -= oldTx.Amount
	} else if oldTx.Type == "expense" {
		oldWallet.Balance += oldTx.Amount
	}

	if err := tx.Save(&oldWallet).Error; err != nil {
		return err
	}

	// 2. Ambil wallet baru (bisa sama atau berbeda dengan wallet lama)
	var newWallet models.Wallet
	if err := tx.First(&newWallet, "id = ? AND user_id = ?", newTxInput.WalletID, oldTx.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("wallet baru tidak ditemukan")
		}
		return err
	}

	// 3. Terapkan saldo transaksi baru ke wallet baru
	if newTxInput.Type == "income" {
		newWallet.Balance += newTxInput.Amount
	} else if newTxInput.Type == "expense" {
		newWallet.Balance -= newTxInput.Amount
	} else {
		return errors.New("tipe transaksi baru tidak valid")
	}

	if err := tx.Save(&newWallet).Error; err != nil {
		return err
	}

	// 4. Perbarui data transaksi
	oldTx.WalletID = newTxInput.WalletID
	oldTx.CategoryID = newTxInput.CategoryID
	oldTx.Amount = newTxInput.Amount
	oldTx.Type = newTxInput.Type
	oldTx.Description = newTxInput.Description
	oldTx.Date = newTxInput.Date

	if err := tx.Save(oldTx).Error; err != nil {
		return err
	}

	return nil
}

// DeleteTransaction menangani penghapusan transaksi dan mengembalikan saldo wallet.
func DeleteTransaction(tx *gorm.DB, transaction *models.Transaction) error {
	// 1. Cari wallet terkait
	var wallet models.Wallet
	if err := tx.First(&wallet, "id = ? AND user_id = ?", transaction.WalletID, transaction.UserID).Error; err != nil {
		return err
	}

	// 2. Kembalikan efek transaksi dari saldo wallet
	if transaction.Type == "income" {
		wallet.Balance -= transaction.Amount
	} else if transaction.Type == "expense" {
		wallet.Balance += transaction.Amount
	}

	if err := tx.Save(&wallet).Error; err != nil {
		return err
	}

	// 3. Hapus data transaksi
	if err := tx.Delete(transaction).Error; err != nil {
		return err
	}

	return nil
}
