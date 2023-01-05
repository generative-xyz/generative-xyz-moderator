package repository

import (
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r Repository) FindUserByWalletAddress(walletAddress string) (*entity.Users, error) {
	resp := &entity.Users{}

	cached, err := r.GetCache(utils.COLLECTION_USERS, walletAddress)
	if err ==  nil && cached != nil {
		err = helpers.Transform(cached, resp)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}

	usr, err := r.FilterOne(utils.COLLECTION_USERS, bson.D{{utils.KEY_WALLET_ADDRESS, walletAddress}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}

	r.CreateCache(utils.COLLECTION_USERS, walletAddress, resp)
	return resp, nil
}

func (r Repository) FindUserByID(userID string) (*entity.Users, error) {
	resp := &entity.Users{}

	usr, err := r.FindOne(utils.COLLECTION_USERS, userID)
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) FindUserByEmail(email string) (*entity.Users, error) {
	resp := &entity.Users{}

	f := bson.D{
		{"$and", 
			bson.A{
				bson.D{{"email", email}},
				bson.D{{utils.KEY_DELETED_AT, nil}},
			},
		},
	}

	usr, err := r.FilterOne(utils.COLLECTION_USERS, f)
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) CreateUser(insertedUser *entity.Users) error {

	err := r.InsertOne(insertedUser.TableName(), insertedUser)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) UpdateUserByWalletAddress(walletAdress string, updateddUser *entity.Users) (*mongo.UpdateResult, error) {
	filter := bson.D{{utils.KEY_WALLET_ADDRESS, walletAdress}}
	result, err := r.UpdateOne(updateddUser.TableName(), filter, updateddUser)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) UpdateUserByID(userID string, updateddUser *entity.Users) (*mongo.UpdateResult, error) {
	filter := bson.D{{utils.KEY_UUID, userID}}
	result, err := r.UpdateOne(updateddUser.TableName(), filter, updateddUser)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Repository) FindUserByAutoUserID(autoUserID int32) (*entity.Users, error) {
	resp := &entity.Users{}

	usr, err := r.FilterOne(utils.COLLECTION_USERS, bson.D{{utils.KEY_AUTO_USERID, autoUserID}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) ListUsers(filter entity.FilterUsers) (*entity.Pagination, error)  {
	users := []entity.Users{}
	resp := &entity.Pagination{}

	filter1 := bson.M{}
	filter1[utils.KEY_DELETED_AT] = nil

	if filter.Email != nil {
		if *filter.Email != "" {
			filter1["email"] = *filter.Email
		}
	}
	
	if filter.WalletAddress != nil {
		if *filter.WalletAddress != "" {
			filter1["wallet_address"] = *filter.WalletAddress
		}
	}
	
	if filter.UserType != nil {
		filter1["user_type"] = *filter.UserType
	}
	
	p, err := r.Paginate(utils.COLLECTION_USERS, filter.Page, filter.Limit, filter1, filter.SortBy, filter.Sort, &users)
	if err != nil {
		return nil, err
	}
	
	resp.Result = users
	//resp.Limit = p.Pagination.PerPage
	resp.Page = p.Pagination.Page
	// resp.Next = p.Pagination.Next
	// resp.Prev = p.Pagination.Prev
	// resp.TotalPage = p.Pagination.TotalPage
	resp.Total = p.Pagination.Total
	return resp, nil
}