package restaurantbiz

import (
	"context"
	"food_delivery/modules/restaurant/restaurantmodel"
)

type CreateRestaurantStore interface {
	Create(ctx context.Context, data *restaurantmodel.RestaurantCreate) error
}

// private
type createRestaurantBiz struct {
	store CreateRestaurantStore
}

// public | export for outside to use
func NewCreateRestaurantBiz(store CreateRestaurantStore) *createRestaurantBiz {
	return &createRestaurantBiz{store: store}
}

// one of method `createRestaurantBiz`
func (biz *createRestaurantBiz) CreateRestaurant(
	ctx context.Context,
	data *restaurantmodel.RestaurantCreate) error {

	if err := data.Validate(); err != nil {
		return err
	}
	err := biz.store.Create(ctx, data)

	return err
}
