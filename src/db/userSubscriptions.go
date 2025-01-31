package db

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// CREATE

func CreateUserSubscription(userID int, subscriptionID int, expiryTime time.Time) (err error) {
	var userSubFind UserSubscriptions
	db := DBSystem

	findUserSub := db.Where("user_id = ?", userID).First(&userSubFind)
	if !errors.Is(findUserSub.Error, gorm.ErrRecordNotFound) {
		return ErrUserSubscriptionExists
	}

	userSub := UserSubscriptions{UserID: userID, SubscriptionID: subscriptionID, UsedBandwidth: 0, Expiry: expiryTime}
	creatUserSub := db.Create(&userSub)
	if creatUserSub.Error != nil {
		combinedLogger.Error("Adding user to db " + creatUserSub.Error.Error())
		return ErrCreatingUserSubscription
	}
	return nil
}

// READ

func ReadUserSubscriptionFromID(userSubID int) (userSubscription UserSubscriptions, err error) {
	db := DBSystem

	findUserSub := db.Where("id = ?", userSubID).First(&userSubscription)
	if errors.Is(findUserSub.Error, gorm.ErrRecordNotFound) {
		err = ErrUserSubscriptionNotFound
	} else if findUserSub.Error != nil {
		err = ErrQuery
	}
	return
}

func ReadUserSubscriptionFromUserID(userID int) (userSubscription UserSubscriptions, err error) {
	db := DBSystem

	findUserSub := db.Where("user_id = ?", userID).First(&userSubscription)
	if errors.Is(findUserSub.Error, gorm.ErrRecordNotFound) {
		err = ErrUserSubscriptionNotFound
	} else if findUserSub.Error != nil {
		err = ErrQuery
	}
	return
}

// Gets all user subscriptions from the database
func ReadAllUserSubscriptions() (userSubs []UserSubscriptions, err error) {
	db := DBSystem

	findUserSubs := db.Find(&userSubs)
	if errors.Is(findUserSubs.Error, gorm.ErrRecordNotFound) {
		err = ErrUserSubscriptionsNotFound
	} else if findUserSubs.Error != nil {
		err = ErrQuery
	}
	return
}

func ReadUserSubscriptionWithSubscriptionID(subscriptionID int) (userSubs []UserSubscriptions, err error) {
	db := DBSystem

	findUserSubs := db.Where("subscription_id = ?", subscriptionID).Find(&userSubs)
	if errors.Is(findUserSubs.Error, gorm.ErrRecordNotFound) {
		err = ErrUserSubscriptionsNotFound
	} else if findUserSubs.Error != nil {
		err = ErrQuery
	}
	return
}

// UPDATE

func UpdateUserSubscription(userSubID int, usedBandwidth int, expiry time.Time) (err error) {
	db := DBSystem
	subs, err := ReadUserSubscriptionFromID(userSubID)
	if err != nil {
		return
	}

	if usedBandwidth > -1 {
		subs.UsedBandwidth = usedBandwidth
	}
	if !expiry.IsZero() {
		subs.Expiry = expiry
	}
	saveErrs := db.Save(&subs)
	if saveErrs.Error != nil {
		err = saveErrs.Error
		return
	}
	return nil
}

// DELETE

func DeleteUserSubscription(userSubID int) (err error) {
	db := DBSystem
	subs, err := ReadUserSubscriptionFromID(userSubID)
	if err != nil {
		return
	}

	_, readErr := readUserKeys(subs.UserID)
	if readErr != nil && readErr != ErrKeyNotFound { //check there is an error and not the error we want
		return readErr
	}
	delErr := db.Delete(&subs)
	if delErr.Error != nil {
		return delErr.Error
	}
	return nil
}

// MISC

// Checks to see if a new key can be added
func checkSubscriptionKeyAddition(userID int) (err error) {
	userSubscription, err := ReadUserSubscriptionFromUserID(userID)
	if err != nil {
		fmt.Println(err)

		return
	}
	if userSubscription.Expiry.After(time.Now()) {
		err = ErrSubscriptionExpired
		return
	}

	subscription, err := ReadSubscription(userSubscription.SubscriptionID)
	if err != nil {
		return
	}

	numKeys := subscription.NumberOfKeys
	userKeys, err := readUserKeys(userID)
	if err != nil {
		return
	}
	if len(userKeys) >= numKeys {
		err = ErrSubscriptionKeyMaxed
		return
	}
	return nil
}

//Validates the user has the subscription ID
func ValidateUsernameUserSubscription(userID int, userSubID int) (err error) {
	_, findUserErr := ReadUser(userID)
	if findUserErr != nil {
		return findUserErr
	}
	userSub, userSubErr := ReadUserSubscriptionFromUserID(userID)
	if userSubErr != nil {
		return userSubErr
	}

	if userSub.ID != userSubID {
		return ErrUserSubscriptionValidation
	}
	return nil
}
