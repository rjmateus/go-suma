package download

import (
	"errors"
	"github.com/rjmateus/go-suma/models/rhnpackage"
	"gorm.io/gorm"
)

func GetDownloadPackage(db *gorm.DB, channel, pkgName, version, release, arch, checksum, epoch string) (rhnpackage.Package, error) {
	sql := `select p.id, p.path, pe.epoch as epoch
                from
                  rhnPackageArch pa,
                  rhnChannelPackage cp,
                  rhnPackage p,
                  rhnChecksum cs,
                  rhnPackageEVR pe,
                  rhnPackageName pn,
                  rhnChannel c
        where 1=1
            and c.label = ?
            and pn.name = ?
            and pe.version = ?
            and pe.release = ?
            and c.id = cp.channel_id
            and pa.label = ?
            and pn.id = p.name_id
            and p.id = cp.package_id
            and p.evr_id = pe.id
            and p.package_arch_id = pa.id
            and p.checksum_id = cs.id`

	parameter := []interface{}{channel, pkgName, version, release, arch}
	if len(checksum) > 0 {
		parameter = append(parameter, checksum)
		sql = sql + " and cs.checksum = ?"
	} else {
		sql = sql + " and cs.checksum is null"
	}

	var queryResult []rhnpackage.Package
	db.Raw(sql, parameter...).Scan(&queryResult)

	if len(queryResult) == 0 {
		return rhnpackage.Package{}, errors.New("no rhnpackage found")
	}

	result := make([]rhnpackage.Package, 0)

	for _, row := range queryResult {
		if len(epoch) > 0 {
			if epoch != row.Epoch {
				continue
			}
		}
		result = append(result, row)
	}
	if len(result) == 0 {
		return rhnpackage.Package{}, errors.New("no rhnpackage found")
	}
	return result[0], nil
}
