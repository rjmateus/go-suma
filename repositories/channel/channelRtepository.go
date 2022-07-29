package channel

import (
	"database/sql"
	"gorm.io/gorm"
)

func IsAccessibleBy(db *gorm.DB, channel string, org int) bool {
	sqlQuery := `SELECT COUNT(*) AS count
          FROM rhnChannel c
          JOIN rhnAvailableChannels ac ON c.id = ac.channel_id
          WHERE c.label = @channel
            AND ac.org_id = @org`

	var queryResult int64
	db.Raw(sqlQuery,
		sql.Named("channel", channel),
		sql.Named("org", org)).Scan(&queryResult)

	return queryResult > 0
}
