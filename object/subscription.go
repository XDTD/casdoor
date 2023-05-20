// Copyright 2023 The Casdoor Authors. All Rights Reserved.
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
	"time"

	"github.com/casdoor/casdoor/util"
	"github.com/xorm-io/core"
)

const defaultStatus = "Pending"

type Subscription struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`
	DisplayName string `xorm:"varchar(100)" json:"displayName"`
	Duration    int    `json:"duration"`

	Description string `xorm:"varchar(100)" json:"description"`
	Plan        string `xorm:"varchar(100)" json:"plan"`

	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`

	User string `xorm:"mediumtext" json:"user"`

	IsEnabled   bool   `json:"isEnabled"`
	Submitter   string `xorm:"varchar(100)" json:"submitter"`
	Approver    string `xorm:"varchar(100)" json:"approver"`
	ApproveTime string `xorm:"varchar(100)" json:"approveTime"`

	State string `xorm:"varchar(100)" json:"state"`
}

func NewSubscription(owner string, user string, plan string, duration int) *Subscription {
	id := util.GenerateId()[:6]
	return &Subscription{
		Name:        "Subscription_" + id,
		DisplayName: "New Subscription - " + id,
		Owner:       owner,
		User:        owner + "/" + user,
		Plan:        owner + "/" + plan,
		CreatedTime: util.GetCurrentTime(),
		State:       defaultStatus,
		Duration:    duration,
		StartDate:   time.Now(),
		EndDate:     time.Now().AddDate(0, 0, duration),
	}
}

func GetSubscriptionCount(owner, field, value string) int {
	session := GetSession(owner, -1, -1, field, value, "", "")
	count, err := session.Count(&Subscription{})
	if err != nil {
		panic(err)
	}

	return int(count)
}

func GetSubscriptionCountByOwners(owners []string, field, value string) int {
	session := GetSessionByOwners(owners, -1, -1, field, value, "", "")
	count, err := session.Count(&Subscription{})
	if err != nil {
		panic(err)
	}

	return int(count)
}

func GetSubscriptions(owner string) []*Subscription {
	subscriptions := []*Subscription{}
	err := adapter.Engine.Desc("created_time").Find(&subscriptions, &Subscription{Owner: owner})
	if err != nil {
		panic(err)
	}

	return subscriptions
}

func GetSubscriptionsByOwners(owners []string) []*Subscription {
	subscriptions := []*Subscription{}
	err := adapter.Engine.Desc("created_time").In("owner", owners).Find(&subscriptions)
	if err != nil {
		panic(err)
	}

	return subscriptions
}

func GetPaginationSubscriptions(owner string, offset, limit int, field, value, sortField, sortOrder string) []*Subscription {
	subscriptions := []*Subscription{}
	session := GetSession(owner, offset, limit, field, value, sortField, sortOrder)
	err := session.Find(&subscriptions)
	if err != nil {
		panic(err)
	}

	return subscriptions
}

func GetPaginationedSubscriptionsByOwners(owners []string, offset, limit int, field, value, sortField, sortOrder string) []*Subscription {
	subscriptions := []*Subscription{}
	session := GetSessionByOwners(owners, offset, limit, field, value, sortField, sortOrder)
	err := session.Find(&subscriptions)
	if err != nil {
		panic(err)
	}

	return subscriptions
}

func getSubscription(owner string, name string) *Subscription {
	if owner == "" || name == "" {
		return nil
	}

	subscription := Subscription{Owner: owner, Name: name}
	existed, err := adapter.Engine.Get(&subscription)
	if err != nil {
		panic(err)
	}

	if existed {
		return &subscription
	} else {
		return nil
	}
}

func GetSubscription(id string) *Subscription {
	owner, name := util.GetOwnerAndNameFromId(id)
	return getSubscription(owner, name)
}

func UpdateSubscription(id string, subscription *Subscription) bool {
	owner, name := util.GetOwnerAndNameFromId(id)
	if getSubscription(owner, name) == nil {
		return false
	}

	affected, err := adapter.Engine.ID(core.PK{owner, name}).AllCols().Update(subscription)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func AddSubscription(subscription *Subscription) bool {
	affected, err := adapter.Engine.Insert(subscription)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func DeleteSubscription(subscription *Subscription) bool {
	affected, err := adapter.Engine.ID(core.PK{subscription.Owner, subscription.Name}).Delete(&Subscription{})
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func (subscription *Subscription) GetId() string {
	return fmt.Sprintf("%s/%s", subscription.Owner, subscription.Name)
}
