package usecase

func (u Usecase) GetAllRedis() ([]string, error) {
	//var res *string
	var err error

	res, err := u.Cache.GetAll()
	if err != nil {
		return nil, err
	}
	return res, err
}

func (u Usecase) DeleteAllRedis() ([]string, error) {
	//var res *string
	var err error

	res, err := u.Cache.GetAll()
	if err != nil {
		return nil, err
	}

	for _, key := range res {
		err = u.Cache.Delete(key)
		if err != nil {
			return nil, err
		}
	}

	return res, err
}

func (u Usecase) GetRedis( key string) (*string, error) {
	var res *string
	var err error

	res, err = u.Cache.GetData(key)

	return res, err
}

func (u Usecase) UpsertRedis( key string, value string) (*string, error) {
	var res *string
	var err error

	err = u.Cache.SetStringData(key, value)
	if err != nil {
		return nil, err
	}

	res, err = u.Cache.GetData(key)

	return res, err
}

func (u Usecase) DeleteRedis( key string) error {
	var err error
	err = u.Cache.Delete(key)
	return err
}
