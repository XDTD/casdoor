// Copyright 2021 The Casdoor Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package object

import (
	"fmt"

	"github.com/casdoor/casdoor/util"
	"github.com/xorm-io/core"
)

type Cert struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`

	DisplayName     string `xorm:"varchar(100)" json:"displayName"`
	Scope           string `xorm:"varchar(100)" json:"scope"`
	Type            string `xorm:"varchar(100)" json:"type"`
	CryptoAlgorithm string `xorm:"varchar(100)" json:"cryptoAlgorithm"`
	BitSize         int    `json:"bitSize"`
	ExpireInYears   int    `json:"expireInYears"`

	Certificate            string `xorm:"mediumtext" json:"certificate"`
	PrivateKey             string `xorm:"mediumtext" json:"privateKey"`
	AuthorityPublicKey     string `xorm:"mediumtext" json:"authorityPublicKey"`
	AuthorityRootPublicKey string `xorm:"mediumtext" json:"authorityRootPublicKey"`
}

func GetMaskedCert(cert *Cert) *Cert {
	if cert == nil {
		return nil
	}

	return cert
}

func GetMaskedCerts(certs []*Cert) []*Cert {
	for _, cert := range certs {
		cert = GetMaskedCert(cert)
	}
	return certs
}

func GetCertCount(owner, field, value string) int {
	session := GetSession("", -1, -1, field, value, "", "")
	count, err := session.Where("owner = ? or owner = ? ", "admin", owner).Count(&Cert{})
	if err != nil {
		panic(err)
	}

	return int(count)
}

func GetCertCountByOwners(owners []string, field, value string) int {
	session := GetSessionByOwners(owners, -1, -1, field, value, "", "")
	count, err := session.Count(&Cert{})
	if err != nil {
		panic(err)
	}

	return int(count)
}

func GetCerts(owner string) []*Cert {
	certs := []*Cert{}
	err := adapter.Engine.Where("owner = ? or owner = ? ", "admin", owner).Desc("created_time").Find(&certs, &Cert{})
	if err != nil {
		panic(err)
	}

	return certs
}

func GetCertsByOwners(owners []string) []*Cert {
	certs := []*Cert{}
	err := adapter.Engine.Desc("created_time").In("owner", owners).Find(&certs)
	if err != nil {
		panic(err)
	}

	return certs
}

func GetPaginationCerts(owner string, offset, limit int, field, value, sortField, sortOrder string) []*Cert {
	certs := []*Cert{}
	session := GetSession("", offset, limit, field, value, sortField, sortOrder)
	err := session.Where("owner = ? or owner = ? ", "admin", owner).Find(&certs)
	if err != nil {
		panic(err)
	}

	return certs
}

func GetPaginationCertsByOwners(owners []string, offset, limit int, field, value, sortField, sortOrder string) []*Cert {
	certs := []*Cert{}
	session := GetSessionByOwners(owners, offset, limit, field, value, sortField, sortOrder)
	err := session.Find(&certs)
	if err != nil {
		panic(err)
	}

	return certs
}

func GetGlobalCertsCount(field, value string) int {
	session := GetSession("", -1, -1, field, value, "", "")
	count, err := session.Count(&Cert{})
	if err != nil {
		panic(err)
	}

	return int(count)
}

func GetGlobleCerts() []*Cert {
	certs := []*Cert{}
	err := adapter.Engine.Desc("created_time").Find(&certs)
	if err != nil {
		panic(err)
	}

	return certs
}

func GetPaginationGlobalCerts(offset, limit int, field, value, sortField, sortOrder string) []*Cert {
	certs := []*Cert{}
	session := GetSession("", offset, limit, field, value, sortField, sortOrder)
	err := session.Find(&certs)
	if err != nil {
		panic(err)
	}

	return certs
}

func getCert(owner string, name string) *Cert {
	if owner == "" || name == "" {
		return nil
	}

	cert := Cert{Owner: owner, Name: name}
	existed, err := adapter.Engine.Get(&cert)
	if err != nil {
		panic(err)
	}

	if existed {
		return &cert
	} else {
		return nil
	}
}

func getCertByName(name string) *Cert {
	if name == "" {
		return nil
	}

	cert := Cert{Name: name}
	existed, err := adapter.Engine.Get(&cert)
	if err != nil {
		panic(err)
	}

	if existed {
		return &cert
	} else {
		return nil
	}
}

func GetCert(id string) *Cert {
	owner, name := util.GetOwnerAndNameFromId(id)
	return getCert(owner, name)
}

func UpdateCert(id string, cert *Cert) bool {
	owner, name := util.GetOwnerAndNameFromId(id)
	if getCert(owner, name) == nil {
		return false
	}

	if name != cert.Name {
		err := certChangeTrigger(name, cert.Name)
		if err != nil {
			return false
		}
	}
	affected, err := adapter.Engine.ID(core.PK{owner, name}).AllCols().Update(cert)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func AddCert(cert *Cert) bool {
	if cert.Certificate == "" || cert.PrivateKey == "" {
		certificate, privateKey := generateRsaKeys(cert.BitSize, cert.ExpireInYears, cert.Name, cert.Owner)
		cert.Certificate = certificate
		cert.PrivateKey = privateKey
	}

	affected, err := adapter.Engine.Insert(cert)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func DeleteCert(cert *Cert) bool {
	affected, err := adapter.Engine.ID(core.PK{cert.Owner, cert.Name}).Delete(&Cert{})
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func (p *Cert) GetId() string {
	return fmt.Sprintf("%s/%s", p.Owner, p.Name)
}

func getCertByApplication(application *Application) *Cert {
	if application.Cert != "" {
		return getCertByName(application.Cert)
	} else {
		return GetDefaultCert()
	}
}

func GetDefaultCert() *Cert {
	return getCert("admin", "cert-built-in")
}

func certChangeTrigger(oldName string, newName string) error {
	session := adapter.Engine.NewSession()
	defer session.Close()

	err := session.Begin()
	if err != nil {
		return err
	}

	application := new(Application)
	application.Cert = newName
	_, err = session.Where("cert=?", oldName).Update(application)
	if err != nil {
		return err
	}

	return session.Commit()
}
