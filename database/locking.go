package database

func ObtainLock(dbox DbBox) (bool, error) {
	if dbox.IsPostgres {
		return postgresObtainLock(dbox)
	} else if dbox.IsMysql {
		return mysqlObtainLock(dbox)
	} else if dbox.IsSqlite {
		return sqliteObtainLock(dbox)
	} else {
		panic("unknown database: " + dbox.Type)
	}
}

func postgresObtainLock(dbox DbBox) (bool, error) {
	result, err := dbox.Db.Exec(`UPDATE migrations_lock SET is_locked = 1 WHERE index = 1 AND is_locked = 0;`)
	if err != nil {
		return false, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return affected == 1, nil
}

func mysqlObtainLock(dbox DbBox) (bool, error) {
	result, err := dbox.Db.Exec(`UPDATE migrations_lock SET is_locked = 1 WHERE ` + "`index`" + ` = 1 AND is_locked = 0;`)
	if err != nil {
		return false, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return affected == 1, nil
}

func sqliteObtainLock(dbox DbBox) (bool, error) {
	result, err := dbox.Db.Exec(`UPDATE migrations_lock SET is_locked = 1 WHERE "index" = 1 AND is_locked = 0;`)
	if err != nil {
		return false, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return affected == 1, nil
}

func ReleaseLock(dbox DbBox) (bool, error) {
	if dbox.IsPostgres {
		return postgresReleaseLock(dbox)
	} else if dbox.IsMysql {
		return mysqlReleaseLock(dbox)
	} else if dbox.IsSqlite {
		return sqliteReleaseLock(dbox)
	} else {
		panic("unknown database: " + dbox.Type)
	}
}

func postgresReleaseLock(dbox DbBox) (bool, error) {
	result, err := dbox.Db.Exec(`UPDATE migrations_lock SET is_locked = 0 WHERE index = 1 AND is_locked = 1;`)
	if err != nil {
		return false, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return affected == 1, nil
}

func mysqlReleaseLock(dbox DbBox) (bool, error) {
	result, err := dbox.Db.Exec(`UPDATE migrations_lock SET is_locked = 0 WHERE ` + "`index`" + ` = 1 AND is_locked = 1;`)
	if err != nil {
		return false, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return affected == 1, nil
}

func sqliteReleaseLock(dbox DbBox) (bool, error) {
	result, err := dbox.Db.Exec(`UPDATE migrations_lock SET is_locked = 0 WHERE "index" = 1 AND is_locked = 1;`)
	if err != nil {
		return false, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return affected == 1, nil
}
