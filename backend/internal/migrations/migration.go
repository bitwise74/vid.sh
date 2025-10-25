package migrations

import "gorm.io/gorm"

func DelCol(dst interface{}, name string, d *gorm.DB) error {
	return d.Migrator().DropColumn(dst, name)
}

func DelTab(dst interface{}, d *gorm.DB) error {
	return d.Migrator().DropTable(dst)
}
