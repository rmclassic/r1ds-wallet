package mysql

import (
	"errors"
	"wallet/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (c *MysqlDatabase) UseDiscount(walletId int, discountId int) error {
	err := c.db.Transaction(func(tx *gorm.DB) error {
		wallet := models.Wallet{ID: walletId}
		if err := tx.Clauses(clause.Locking{
			Strength: "UPDATE",
			Options:  "NOWAIT",
		}).Find(&wallet, &wallet); err != nil {
			return err.Error
		}

		discount := models.Discount{ID: discountId}
		if err := tx.Clauses(clause.Locking{
			Strength: "UPDATE",
			Options:  "NOWAIT",
		}).Find(&discount, &discount); err != nil {
			return err.Error
		}

		if discount.QuantityLeft <= 0 {
			return errors.New("no more discounts left")
		}

		discount.QuantityLeft--
		if err := tx.Save(discount); err != nil {
			return err.Error
		}

		// TODO: Add discount usage row
		wallet.Balance += discount.Amount
		if err := tx.Save(wallet); err != nil {
			return err.Error
		}
		return nil
	})

	return err
}

func (db *MysqlDatabase) WalletAdd(wallet *models.Wallet) error {
	return db.db.Create(wallet).Error
}

func (db *MysqlDatabase) WalletUpdate(wallet *models.Wallet) error {
	return db.db.Save(wallet).Error
}

func (db *MysqlDatabase) WalletGet(wallet *models.Wallet) error {
	return db.db.Preload("User").Find(wallet, wallet).Error
}
